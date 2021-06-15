package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/blmhemu/consul-ext-dns/config"
	"github.com/hashicorp/consul/api"
)

var version = "0.1"

func main() {

	cfg, err := config.Load(os.Args, os.Environ())
	if err != nil {
		log.Printf("[FATAL] %s. %s", version, err)
		os.Exit(1)
	}
	if cfg == nil {
		fmt.Printf("%s %s\n", version, runtime.Version())
		return
	}

	// We recieve updates on this channel
	updates := make(chan []*api.ServiceEntry)

	// Watch for any changes in service
	go watchLB(cfg, updates)

	// Process any changes
	oldIPSet := make(map[string]struct{})
	for svcs := range updates {
		newIPSet := make(map[string]struct{})
		for _, svc := range svcs {
			newIPSet[svc.Node.Address] = struct{}{}
		}
		if changed(newIPSet, oldIPSet) {
			newIPs := make([]string, 0, len(newIPSet))
			for ip := range newIPSet {
				newIPs = append(newIPs, ip)
			}
			// No IPs means lb is down
			// either momentarily or for an extended period.
			// We refrain from setting empty A records in DNS
			// When lb is back online, we update the IPs anyways if there is a change
			if len(newIPs) != 0 {
				updateDNSRecords(newIPs)
				oldIPSet = newIPSet
			} else {
				fmt.Println("Lb is down")
			}
		}
	}
}

func changed(newIPSet, oldIPSet map[string]struct{}) bool {
	if len(newIPSet) != len(oldIPSet) {
		return true
	}
	for ip := range newIPSet {
		if _, ok := oldIPSet[ip]; !ok {
			return true
		}
	}
	return false
}

func updateDNSRecords(ips []string) {
	fmt.Println(ips)
}

type MyQueryOptions struct{ api.QueryOptions }

func (qo *MyQueryOptions) fetchUpdates(client *api.Client) ([]*api.ServiceEntry, error) {
	svccfg, qm, err := client.Health().Service("fabio", "", true, &qo.QueryOptions)

	if err != nil || qm.LastIndex <= qo.WaitIndex {
		qo.WaitIndex = 0
	} else {
		qo.WaitIndex = qm.LastIndex
	}
	return svccfg, err
}

func watchLB(cfg *config.Config, updates chan []*api.ServiceEntry) {
	client, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Addr,
		Scheme:  cfg.Consul.Scheme,
	})
	if err != nil {
		panic(err)
	}

	qo := MyQueryOptions{
		QueryOptions: api.QueryOptions{
			RequireConsistent: true,
			WaitIndex:         0,
		},
	}

	for {
		svccfg, err := qo.fetchUpdates(client)
		if err != nil {
			panic(err)
		}
		updates <- svccfg
	}
}
