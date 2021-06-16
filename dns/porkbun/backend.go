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
	Domain string
	Name   string
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
	pbclient := &PBClient{Client: client, Domain: cfg.Domain, Name: cfg.Name}
	pbclient.updateState()
	return pbclient, nil
}

func (p *PBClient) WriteRecords(newRecords sets.String) error {
	currRecords := p.fetchIPSet()
	deletions := currRecords.Difference(newRecords)
	additions := newRecords.Difference(*currRecords)
	// No additions / deletions ? Good return
	if !deletions.HasAny() && !additions.HasAny() {
		return nil
	}
	// This means the LB is down
	// No need to update until it is back up again
	if !newRecords.HasAny() {
		return nil
	}
	// Porkbun does not have delete all function.
	// To minimize API calls, we instead use edit when available
	for {
		newip, newany := additions.PopAny()
		oldip, oldany := deletions.PopAny()
		if newany {
			dnsRecord := porkbun.DNSRecord{
				Type:    A,
				Content: newip,
			}
			if p.Name != "" {
				dnsRecord.Name = p.Name
			}
			if oldany {
				p.Client.EditRecord(p.Domain, p.State[oldip], &dnsRecord)
			} else {
				p.Client.CreateRecord(p.Domain, &dnsRecord)
			}
		} else if oldany {
			p.Client.DeleteRecord(p.Domain, p.State[oldip])
		} else {
			break
		}
	}
	return nil
}

// Helper Land
func (p *PBClient) fetchIPSet() *sets.String {
	ipSet := sets.NewString()
	for ip := range p.State {
		ipSet.Insert(ip)
	}
	return &ipSet
}

func (p *PBClient) updateState() error {
	dnsResp, err := p.Client.RetrieveRecords(p.Domain)
	if err != nil {
		return err
	}
	ipIDMap := getIPIDMap(dnsResp)
	p.State = ipIDMap
	return nil
}

func getIPIDMap(dnsResp *porkbun.DNSResponse) map[string]string {
	ipIDMap := make(map[string]string)
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
