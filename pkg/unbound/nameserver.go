package unbound

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// UseDNSInternally is to change the Go program DNS only.
func (c *configurator) UseDNSInternally(ip net.IP) {
	net.DefaultResolver.PreferGo = true
	net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(ctx, "udp", net.JoinHostPort(ip.String(), "53"))
	}
}

// UseDNSSystemWide changes the nameserver to use for DNS system wide.
func (c *configurator) UseDNSSystemWide(ip net.IP, keepNameserver bool) error {
	const filepath = resolvConfFilepath
	file, err := c.openFile(filepath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return err
	}
	s := strings.TrimSuffix(string(data), "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	found := false
	if !keepNameserver { // default
		for i := range lines {
			if strings.HasPrefix(lines[i], "nameserver ") {
				lines[i] = "nameserver " + ip.String()
				found = true
			}
		}
	}
	if !found {
		lines = append(lines, "nameserver "+ip.String())
	}
	s = strings.Join(lines, "\n") + "\n"
	_, err = file.WriteString(s)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
