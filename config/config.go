package config

type Config struct {
	Consul Consul
	DNS    DNS
}

type Consul struct {
	Addr   string
	Scheme string
}

type DNS struct {
	Backend    string
	Cloudflare Cloudflare
	Porkbun    Porkbun
}

type Cloudflare struct {
	APIToken string
	Email    string
}

type Porkbun struct {
	APIKey       string
	SecretAPIKey string
	Domain       string
	Name         string
}

var defaultConfig = Config{
	Consul: defaultConsul,
	DNS:    defaultDNS,
}

var defaultConsul = Consul{
	Addr:   "localhost:8500",
	Scheme: "http",
}

var defaultDNS = DNS{
	Backend:    "",
	Cloudflare: defaultCloudflare,
	Porkbun:    defaultPorkbun,
}

var defaultCloudflare = Cloudflare{
	APIToken: "",
	Email:    "",
}

var defaultPorkbun = Porkbun{
	APIKey:       "",
	SecretAPIKey: "",
	Domain:       "",
	Name:         "",
}
