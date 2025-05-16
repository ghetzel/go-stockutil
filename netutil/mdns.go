package netutil

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ghetzel/go-defaults"
	"github.com/ghetzel/go-stockutil/rxutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/grandcat/zeroconf"
)

var registered sync.Map

type ZeroconfOptions struct {
	Context       context.Context
	Limit         int
	Timeout       time.Duration `default:"30s"`
	Service       string        `default:"_http._tcp"`
	Domain        string        `default:".local"`
	CheckInterval time.Duration `default:"100ms"`
	MatchInstance string
	MatchPort     string
	MatchHostname string
	MatchAddress  string
}

type Service struct {
	Hostname   string          `json:"hostname"`
	Instance   string          `json:"instance"`
	Service    string          `json:"service"`
	Domain     string          `json:"domain"`
	Port       int             `json:"port"`
	Text       []string        `json:"txt"`
	Address    string          `json:"address"`
	Addresses  []net.IP        `json:"addresses,omitempty"`
	Interfaces []net.Interface `json:"interfaces,omitempty"`
}

func (self *Service) String() string {
	return strings.TrimSuffix(
		fmt.Sprintf("%s.%s%s:%d %v", self.Instance, self.Service, self.Domain, self.Port, self.Text),
		` []`,
	)
}

type ServiceFunc func(*Service) bool

func isEntryMatch(options *ZeroconfOptions, entry *zeroconf.ServiceEntry) bool {
	if rx := options.MatchInstance; rx == `` || rxutil.IsMatchString(rx, entry.Instance) {
		return true
	} else if rx := options.MatchPort; rx == `` || rxutil.IsMatchString(rx, typeutil.String(entry.Port)) {
		return true
	} else if rx := options.MatchHostname; rx == `` || rxutil.IsMatchString(rx, typeutil.String(entry.HostName)) {
		return true
	} else if rx := options.MatchAddress; rx == `` {
		return true
	} else {
		for _, ip := range entry.AddrIPv4 {
			if rxutil.IsMatchString(rx, ip.String()) {
				return true
			}
		}

		for _, ip := range entry.AddrIPv6 {
			if rxutil.IsMatchString(rx, ip.String()) {
				return true
			}
		}
	}

	return false
}

// Perform Multicast DNS discovery on the local network, calling the fn callback for each
// discovered service.
func ZeroconfDiscover(options *ZeroconfOptions, fn ServiceFunc) error {
	if fn == nil {
		return fmt.Errorf("Must provide a callback function to receive discover services")
	}

	if options == nil {
		options = new(ZeroconfOptions)
	}

	defaults.SetDefaults(options)

	if options.Context == nil {
		options.Context = context.Background()
	}

	var found = 0

	// setup mDNS resolver
	if resolver, err := zeroconf.NewResolver(
		zeroconf.SelectIPTraffic(zeroconf.IPv4AndIPv6),
	); err == nil {
		var entries = make(chan *zeroconf.ServiceEntry)
		ctx, cancel := context.WithTimeout(options.Context, options.Timeout)
		defer cancel()

		// receive discovered services
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				if isEntryMatch(options, entry) {
					found += 1
					var addrs = make([]net.IP, 0)
					addrs = append(addrs, entry.AddrIPv4...)
					addrs = append(addrs, entry.AddrIPv6...)
					var addr = ``

					if len(addrs) > 0 {
						addr = fmt.Sprintf("%v:%d", addrs[0], entry.Port)
					}

					// fire off callback for this service
					if !fn(&Service{
						Hostname:  entry.HostName,
						Instance:  entry.Instance,
						Service:   entry.Service,
						Port:      entry.Port,
						Domain:    entry.Domain,
						Text:      entry.Text,
						Addresses: addrs,
						Address:   addr,
					}) {
						cancel()
					}
				}

				if options.Limit > 0 && found >= options.Limit {
					cancel()
				}
			}
		}(entries)

		// actually start mDNS discovery
		if err := resolver.Browse(ctx, options.Service, options.Domain, entries); err == nil {
			select {
			case <-ctx.Done():
			}

			return nil
		} else {
			return fmt.Errorf("browse error: %v", err)
		}
	} else {
		return err
	}
}

// Register the given service in Multicast DNS.  Returns an ID that can be used to unregister
// the service later.
func ZeroconfRegister(svc *Service) (string, error) {
	if svc == nil {
		return ``, fmt.Errorf("Must provide a service configuration to register mDNS")
	} else if svc.Instance == `` {
		svc.Instance = stringutil.UUID().String()
	} else if svc.Service == `` {
		return ``, fmt.Errorf("Must provide a service type")
	} else if svc.Domain == `` {
		return ``, fmt.Errorf("Must provide a service domain")
	} else if svc.Port == 0 {
		return ``, fmt.Errorf("Must specify a service port")
	}

	var slug = fmt.Sprintf("%x", sha256.Sum256(
		[]byte(fmt.Sprintf("%s.%s%s:%d", svc.Instance, svc.Service, svc.Domain, svc.Port)),
	))

	if _, ok := registered.Load(slug); ok {
		return slug, fmt.Errorf("A service matching these parameters is already registered")
	}

	if server, err := zeroconf.Register(
		svc.Instance,
		svc.Service,
		svc.Domain,
		svc.Port,
		svc.Text,
		svc.Interfaces,
	); err == nil {
		registered.Store(slug, server)
		return slug, nil
	} else {
		return ``, err
	}
}

// Unregister a previously-registered service.
func ZeroconfUnregister(id string) {
	defer registered.Delete(id)

	if s, ok := registered.Load(id); ok {
		if server, ok := s.(*zeroconf.Server); ok {
			server.Shutdown()
		}
	}
}

// Unregister all Multicast DNS services.
func ZeroconfUnregisterAll() {
	registered.Range(func(key, value any) bool {
		ZeroconfUnregister(typeutil.String(key))
		return true
	})
}
