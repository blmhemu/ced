package cloudflare

import (
	"net/http"

	"github.com/blmhemu/consul-ext-dns/config"
	"github.com/blmhemu/consul-ext-dns/dns"
)

type cf struct {
	client *http.Client
}

func NewBackend(cfg *config.Cloudflare) (dns.Backend, error) {
	return &cf{client: http.DefaultClient}, nil
}

func (c *cf) WriteRecords(dns string, ips []string) error {
	return nil
}
