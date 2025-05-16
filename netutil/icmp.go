// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var WaitForPingReply = 5 * time.Second
var WaitForPingCheckInterval = time.Second

// Periodically attempts to send an ICMP ECHO request (a "ping") to the given IP address, up to totaltime.
// Returns nil if an ECHO reply was received, or an error if the function timed out.
// The check interval can be configured using the WaitForPingCheckInterval package variable.
func WaitForPing(addr any, totaltime time.Duration) error {
	var started = time.Now()

	var address net.IP

	if a, ok := addr.(net.IP); ok {
		address = a
	} else if a, ok := addr.(string); ok {
		address = net.ParseIP(a)
	} else if a, ok := addr.(*IPAddress); ok {
		address = a.IP
	} else {
		return fmt.Errorf("Address must be a string, net.IP, or *IPAddress")
	}

	for time.Since(started) < totaltime {
		if err := OnePingOnly(address, nil, WaitForPingReply); err == nil {
			return nil
		} else if strings.HasPrefix(err.Error(), `fatal:`) {
			return err
		}

		time.Sleep(WaitForOpenCheckInterval)
	}

	return fmt.Errorf("Timed out attempting to ping %s after %v", address, time.Since(started))
}

// Like WaitForPing, but will identify the default gateway and ping that address.
func WaitForGatewayPing(totaltime time.Duration) error {
	if addr := DefaultAddress(); addr != nil {
		return WaitForPing(addr.Gateway, totaltime)
	} else {
		return fmt.Errorf("no default gateway found")
	}
}

// Like WaitForGatewayPing, but specifically pings an IPv6 gateway
func WaitForGatewayPing6(totaltime time.Duration) error {
	if addr := DefaultAddress6(); addr != nil {
		return WaitForPing(addr.Gateway, totaltime)
	} else {
		return fmt.Errorf("no default gateway found")
	}
}

// Send a single ICMP ECHO request packet to the given address on the given interface and wait for
// up to timeout for a reply.
func OnePingOnly(dest net.IP, source *IPAddress, timeout time.Duration) error {
	switch runtime.GOOS {
	case `darwin`:
	case `linux`:
	default:
		return fmt.Errorf("fatal: not supported on %v", runtime.GOOS)
	}

	var proto string
	var icmptyp icmp.Type

	if source == nil {
		source = DefaultAddress()

		if source == nil {
			return fmt.Errorf("fatal: no default interface to ping from")
		}
	}

	if len(source.IP) > 32 {
		proto = `udp6`
		icmptyp = ipv6.ICMPTypeEchoRequest
	} else {
		proto = `udp4`
		icmptyp = ipv4.ICMPTypeEcho
	}

	if conn, err := icmp.ListenPacket(proto, source.IP.String()+`%`+source.Interface.Name); err == nil {
		defer conn.Close()

		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			return fmt.Errorf("bad timeout: %v", err)
		}

		var req = icmp.Message{
			Type: icmptyp,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}

		if icmpbin, err := req.Marshal(nil); err == nil {
			if _, err := conn.WriteTo(icmpbin, &net.UDPAddr{
				IP:   source.IP,
				Zone: source.Interface.Name,
			}); err != nil {
				return fmt.Errorf("failed to send ping: %v", err)
			}

			var replybin = make([]byte, 1500)

			if n, _, err := conn.ReadFrom(replybin); err == nil {
				if _, err := icmp.ParseMessage(58, replybin[:n]); err == nil {
					return nil
				} else {
					return fmt.Errorf("bad ICMP reply: %v", err)
				}
			} else {
				return fmt.Errorf("failed to read reply: %v", err)
			}
		} else {
			return fmt.Errorf("fatal: bad outgoing ICMP packet: %v", err)
		}
	} else {
		return fmt.Errorf("fatal: failed to setup ping reply listener: %v", err)
	}
}
