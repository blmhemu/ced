package dns

type Backend interface {
	WriteRecords(dns string, ips []string) error
}

var Default Backend
