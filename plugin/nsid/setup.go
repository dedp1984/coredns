package nsid

import (
	"os"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("nsid", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	nsid, err := nsidParse(c)
	if err != nil {
		return plugin.Error("nsid", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Nsid{Next: next, Data: nsid}
	})

	return nil
}

func nsidParse(c *caddy.Controller) (string, error) {
	// Use hostname as the default
	nsid, err := os.Hostname()
	if err != nil {
		nsid = "localhost"
	}
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) == 0 {
			return nsid, nil
		}
		nsid = strings.Join(args, " ")
		return nsid, nil
	}
	return nsid, nil
}
