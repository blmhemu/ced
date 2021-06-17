package main

import (
	// See if we can use zerolog
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/blmhemu/ced/config"
	dnsbackend "github.com/blmhemu/ced/dns"
	"github.com/blmhemu/ced/dns/porkbun"
	"github.com/blmhemu/ced/exit"
	"github.com/hashicorp/consul/api"
	"k8s.io/apimachinery/pkg/util/sets"
)

var version = "dev"

func main() {

	cfg, err := config.Load(os.Args, os.Environ())
	if err != nil {
		exit.Fatalf("[FATAL] %s. %s", version, err)
	}
	if cfg == nil {
		fmt.Printf("%s %s\n", version, runtime.Version())
		return
	}
	exit.Listen(func(s os.Signal) {})

	// We recieve updates on this channel
	updates := make(chan []*api.ServiceEntry)

	// Initialize a DNS backend
	initBackend(cfg)

	// Watch for any changes in service
	go watchLB(cfg, updates)

	// Process changes
	go pushUpdatesToBackend(updates)

	// Wait till end
	exit.Wait()
	log.Print("[INFO] Down")
}

func initBackend(cfg *config.Config) {
	var err error
	switch cfg.DNS.Backend {
	case dnsbackend.Porkbun:
		dnsbackend.Default, err = porkbun.NewBackend(&cfg.DNS.Porkbun)
	default:
		log.Printf("[FATAL] Please provide a valid DNS backend")
		exit.Exit(1)
	}
	if err != nil {
		log.Printf("[FATAL] Cannot initialize DNS backend %s", err)
		exit.Exit(1)
	}

}

func watchLB(cfg *config.Config, updates chan []*api.ServiceEntry) {
	client, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Addr,
		Scheme:  cfg.Consul.Scheme,
	})
	if err != nil {
		panic(err)
	}

	qo := ConsulQueryOpts{
		QueryOptions: api.QueryOptions{
			RequireConsistent: true,
			WaitIndex:         0,
		},
	}

	for {
		svccfg, err := qo.fetchUpdates(cfg.Service, client)
		if err != nil {
			panic(err)
		}
		updates <- svccfg
	}
}

func pushUpdatesToBackend(updates chan []*api.ServiceEntry) {
	// Continously fetch and push updates
	for svcs := range updates {
		newIPSet := sets.NewString()
		for _, svc := range svcs {
			newIPSet.Insert(svc.Node.Address)
		}
		dnsbackend.Default.WriteRecords(newIPSet)
	}
}

// Wrapper around existing QueryOptions to impl some methods
// Helper land
type ConsulQueryOpts struct{ api.QueryOptions }

func (qo *ConsulQueryOpts) fetchUpdates(service string, client *api.Client) ([]*api.ServiceEntry, error) {
	svccfg, qm, err := client.Health().Service(service, "", true, &qo.QueryOptions)
	if err != nil || qm.LastIndex <= qo.WaitIndex {
		qo.WaitIndex = 0
	} else {
		qo.WaitIndex = qm.LastIndex
	}
	return svccfg, err
}
