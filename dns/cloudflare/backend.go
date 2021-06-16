package cloudflare

import (
	"fmt"
	"os"

	"github.com/blmhemu/consul-ext-dns/config"
	"github.com/blmhemu/consul-ext-dns/dns"
	cfapi "github.com/cloudflare/cloudflare-go"
)

type cf struct {
	api *cfapi.API
}

func NewBackend(cfg *config.Cloudflare) (dns.Backend, error) {
	api, err := cfapi.New(cfg.APIToken, cfg.APIToken)
	if err != nil {
		fmt.Printf("[ERROR] Error creating cloudflare client %s", err)
		os.Exit(1)
	}
	return &cf{api: api}, nil
}

func (c *cf) WriteRecords(dns string, ips []string) error {
	return nil
}
