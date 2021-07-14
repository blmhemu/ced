package dnsbackend

import "k8s.io/apimachinery/pkg/util/sets"

const Porkbun = "porkbun"

type Backend interface {
	WriteRecords(ips sets.String) error
}

var Default Backend
