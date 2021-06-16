package porkbun

import (
	"net/http"
	"time"

	"github.com/blmhemu/consul-ext-dns/config"
	dnsbackend "github.com/blmhemu/consul-ext-dns/dns"
	porkbun "github.com/blmhemu/porkbun-go"
	"k8s.io/apimachinery/pkg/util/sets"
)

const A = "A"

type PBClient struct {
	Client *porkbun.Client
	State  map[string]string // IP -> Porkbun ID
}

func NewBackend(cfg *config.Porkbun) (dnsbackend.Backend, error) {
	pbCfg := porkbun.Config{
		Auth: porkbun.Auth{
			APIKey:       cfg.APIKey,
			SecretAPIKey: cfg.SecretAPIKey,
		},
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	client, err := porkbun.NewClient(&pbCfg)
	if err != nil {
		return nil, err
	}
	pbclient := &PBClient{Client: client}
	pbclient.updateState(cfg.Domain)
	return pbclient, nil
}

func (p *PBClient) WriteRecords(newRecords sets.String) error {

	return nil
}

// Helper Land
func (p *PBClient) updateState(domain string) error {
	dnsResp, err := p.Client.RetrieveRecords(domain)
	if err != nil {
		return err
	}
	ipIDMap := getIPIDMap(dnsResp)
	p.State = ipIDMap
	return nil
}

func getIPIDMap(dnsResp *porkbun.DNSResponse) map[string]string {
	var ipIDMap map[string]string
	if dnsResp.Records == nil {
		return ipIDMap
	}
	for _, rec := range dnsResp.Records {
		if rec.Type == A {
			ipIDMap[rec.Content] = rec.ID
		}
	}
	return ipIDMap
}
