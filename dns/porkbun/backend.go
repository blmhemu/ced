package porkbun

import (
	"net/http"
	"time"

	"github.com/blmhemu/consul-ext-dns/config"
	"github.com/blmhemu/consul-ext-dns/dns"
)

type Porkbun struct {
	client    *http.Client
	recordMap map[string]map[string]string // DNS -> IP -> Porkbun ID
}

func NewBackend(cfg *config.Porkbun) (dns.Backend, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Porkbun{client: client}, nil
}

func (p *Porkbun) WriteRecords(dns string, newRecords []string) error {
	return nil
}
