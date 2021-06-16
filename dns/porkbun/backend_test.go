package porkbun

import (
	"encoding/json"
	"testing"
)

func TestFetchAllRecords(t *testing.T) {
	lee := dnsRecordWithAuth{
		Auth: Auth{
			APIKey:       "lol",
			SecretAPIKey: "lol2",
		},
		DNSRecord: DNSRecord{
			ID: "id1",
		},
	}
	s, _ := json.Marshal(lee)

	t.Logf(string(s))
	// resp, err := c.config.Client.Post(fmt.Sprintf(PORKBUN_DNS_RETRIEVE, dns), "application/json")
	// if err != nil {
	// 	return "", err
	// }
	// var body string
	// resp.Body.Read(body)

}
