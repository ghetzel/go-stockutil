// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/jackpal/gateway"
	"github.com/phayes/freeport"
)

var DefaultWaitForOpenConnectionTimeout = 5 * time.Second
var WaitForOpenCheckInterval = time.Second

// Periodically attempts to open a TCP socket to the given address, up to totaltime.  Returns nil if the socket
// was opened successfully (which will then immediately be closed), or an error if the function timed out.
// The check interval can be configured using the WaitForOpenCheckInterval package variable.
func WaitForOpen(network string, address string, totaltime time.Duration, timeouts ...time.Duration) error {
	var started = time.Now()
	var timeout time.Duration

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	} else {
		timeout = DefaultWaitForOpenConnectionTimeout
	}

	for time.Since(started) < totaltime {
		if conn, _ := net.DialTimeout("tcp", address, timeout); conn != nil {
			conn.Close()
			return nil
		}

		time.Sleep(WaitForOpenCheckInterval)
	}

	return fmt.Errorf("Timed out waiting for %s/%s to open", network, address)
}

// Retrieve an open ephemeral port.
func EphemeralPort() (int, error) {
	return freeport.GetFreePort()
}

// Takes an address in the form of "host:port", looks for port zero (e.g: ":0"),
// and gets an ephemeral local port and returns that address (e.g.: ":41327").
func ExpandPort(address string) string {
	if host, port, err := net.SplitHostPort(address); err == nil {
		if p := typeutil.Int(port); p == 0 {
			if p, err := EphemeralPort(); err == nil {
				return fmt.Sprintf("%v:%d", host, p)
			} else {
				panic(err.Error())
			}
		}
	}

	return address
}

// Retrieves the default gateway interface.
func DefaultGateway() (net.IP, error) {
	return gateway.DiscoverGateway()
}

type IPAddress struct {
	IP        net.IP
	Mask      net.IPMask
	Interface net.Interface
	Gateway   net.IP
}

// Return a list of routable IP addresses, along with their associated gateways and interfaces.
func RoutableAddresses() ([]*IPAddress, error) {
	var addresses = make([]*IPAddress, 0)

	var gateways []net.IP

	// get the default gateway (IPv4)
	if gw, err := DefaultGateway(); err == nil {
		gateways = append(gateways, gw)
	}

	// get the default gateway (IPv6)
	// TODO: need to implement DefaultGateway6()
	// if gw6, err := DefaultGateway6(); err == nil {
	// 	gateways = append(gateways, gw6)
	// }

	if len(gateways) == 0 {
		return nil, fmt.Errorf("no routable gateways found")
	}

	for _, gw := range gateways {
		if ifaces, err := net.Interfaces(); err == nil {
			// for each interface...
			for _, iface := range ifaces {
				if addrs, err := GetRoutableAddresses(gw, &iface); err == nil {
					addresses = append(addresses, addrs...)
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	return addresses, nil
}

// Returns all addresses on the given interface that can route to the given gateway.  If gateway is nil,
// the default gateway will be attempted.
func GetRoutableAddresses(gateway net.IP, iface *net.Interface) ([]*IPAddress, error) {
	var addresses = make([]*IPAddress, 0)

	if gateway == nil {
		if gw, err := DefaultGateway(); err == nil {
			gateway = gw
		} else {
			return nil, fmt.Errorf("failed to retrieve default gateway: %v", err)
		}
	}

	if addrs, err := iface.Addrs(); err == nil {
		// for each address on this interface...
		for _, addr := range addrs {
			// only consider IP addresses at the moment
			if network, ok := addr.(*net.IPNet); ok {
				// if this addresses network contains the gateway, we found a usable address
				if network.Contains(gateway) {
					addresses = append(addresses, &IPAddress{
						IP:        network.IP,
						Mask:      network.Mask,
						Interface: *iface,
						Gateway:   gateway,
					})
				}
			}
		}

		return addresses, nil
	} else {
		return nil, err
	}
}

// Retrieves the first routable IP address on any interface that falls inside of the
// system's default gateway network.  Will return nil if no IP could be found.
func DefaultAddress() *IPAddress {
	if addrs, err := RoutableAddresses(); err == nil && len(addrs) > 0 {
		return addrs[0]
	}

	return nil
}

// Like DefaultAddress, but specifically filters on IPv6 addresses.  If no routable
// IPv6 address is found, returns nil.
func DefaultAddress6() *IPAddress {
	if addrs, err := RoutableAddresses(); err == nil && len(addrs) > 0 {
		for _, addr := range addrs {
			if len(addr.IP) > 32 {
				return addr
			}
		}
	}

	return nil
}

// Return the current machine's Fully-qualified domain name,
func FQDN() string {
	if hostname, err := os.Hostname(); err == nil {
		if responses, err := net.LookupIP(hostname); err == nil {
			for _, addr := range responses {
				if ipv4 := addr.To4(); ipv4 != nil {
					if hosts, err := net.LookupAddr(ipv4.String()); err == nil && len(hosts) > 0 {
						for _, host := range hosts {
							if fqdn := strings.TrimSuffix(host, "."); fqdn != `` {
								return fqdn
							}
						}
					}
				}
			}
		}

		return hostname
	} else {
		return ``
	}
}
