package porkbun

import (
	"fmt"
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

// How to handle the bunch of errors ??? Currently just loggin them.
func (p *PBClient) WriteRecords(newRecords sets.String) error {
	// This means the LB is down
	// No need to update until it is back up again
	if !newRecords.HasAny() {
		return nil
	}
	currRecords := p.fetchIPSet()
	deletions := currRecords.Difference(newRecords)
	additions := newRecords.Difference(*currRecords)
	// No additions / deletions ? Good return
	if !deletions.HasAny() && !additions.HasAny() {
		return nil
	}
	// Porkbun does not have delete all function.
	// To minimize API calls, we instead use (this seemingly complex) edits when available
	anyErrors := false
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
				err := p.Client.EditRecord(p.Domain, p.State[oldip], &dnsRecord)
				if err != nil {
					fmt.Printf("[ERROR] Failed editing record %s", err)
					anyErrors = true
				} else {
					p.State[newip] = p.State[oldip]
					delete(p.State, oldip)
				}
			} else {
				id, err := p.Client.CreateRecord(p.Domain, &dnsRecord)
				if err != nil {
					fmt.Printf("[ERROR] Failed creating record %s", err)
					anyErrors = true
				} else {
					p.State[newip] = id
					delete(p.State, oldip)
				}
			}
		} else if oldany {
			if err := p.Client.DeleteRecord(p.Domain, p.State[oldip]); err != nil {
				fmt.Printf("[ERROR] Failed deleting record %s", err)
				anyErrors = true
			} else {
				delete(p.State, oldip)
			}
		} else {
			break
		}
	}
	if anyErrors {
		if err := p.updateState(); err != nil {
			fmt.Printf("[ERROR] Failed updating state %s", err)
		}
		return fmt.Errorf("There were some errors in some of the calls. Please check the logs.")
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

// Removes existing state and fetches new state from remote
func (p *PBClient) updateState() error {
	dnsResp, err := p.Client.RetrieveRecords(p.Domain)
	if err != nil {
		return err
	}
	p.State = getIPIDMap(dnsResp)
	return nil
}

func getIPIDMap(dnsResp *porkbun.DNSResponse) map[string]string {
	ipIDMap := make(map[string]string)
	if len(dnsResp.Records) == 0 {
		return ipIDMap
	}
	for _, rec := range dnsResp.Records {
		if rec.Type == A {
			ipIDMap[rec.Content] = rec.ID
		}
	}
	return ipIDMap
}
