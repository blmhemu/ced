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
}

type Cloudflare struct {
	APIToken string
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
}

var defaultCloudflare = Cloudflare{
	APIToken: "",
}
