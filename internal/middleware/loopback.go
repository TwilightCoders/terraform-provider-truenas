package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// loopbackHTTPClient returns a plaintext client that can only connect to loopback addresses.
//
// TrueNAS revokes an API key sent over plaintext unless the connection reaches middlewared from a
// loopback address, as it does through an SSH tunnel to 127.0.0.1. The check runs at dial time on
// the resolved address, so neither a hostname nor DNS can route plaintext off the machine.
func loopbackHTTPClient() *http.Client {
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	return &http.Client{Transport: &http.Transport{
		Proxy: nil, // an HTTP proxy would carry the key off-host
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				if !ip.IP.IsLoopback() {
					return nil, &NotLoopbackError{Host: host, Address: ip.IP.String()}
				}
			}
			if len(ips) == 0 {
				return nil, fmt.Errorf("%s did not resolve", host)
			}
			var lastErr error
			for _, ip := range ips {
				conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip.IP.String(), port))
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			return nil, lastErr
		},
	}}
}

// NotLoopbackError reports an attempt to use plaintext transport to a non-loopback address.
type NotLoopbackError struct {
	Host    string
	Address string
}

func (e *NotLoopbackError) Error() string {
	return fmt.Sprintf("insecure_loopback only connects to loopback addresses, but %s resolves to %s; "+
		"TrueNAS would revoke the API key. Use an SSH tunnel to 127.0.0.1 or connect over HTTPS.", e.Host, e.Address)
}
