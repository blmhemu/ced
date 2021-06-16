package dnsbackend

import "k8s.io/apimachinery/pkg/util/sets"

const Cloudflare = "Cloudflare"
const Porkbun = "Porkbun"

type Backend interface {
	WriteRecords(ips sets.String) error
}

var Default Backend
