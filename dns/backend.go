package dns

const Cloudflare = "Cloudflare"
const Porkbun = "Porkbun"

type Backend interface {
	WriteRecords(dns string, ips []string) error
}

var Default Backend
