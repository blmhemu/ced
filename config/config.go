package config

type Config struct {
	Service string
	Consul  Consul
	DNS     DNS
}

type Consul struct {
	Addr     string
	Scheme   string
	CAFile   string
	CertFile string
	KeyFile  string
}

type DNS struct {
	Backend string
	Porkbun Porkbun
}

type Porkbun struct {
	APIKey       string
	SecretAPIKey string
	Domain       string
	Name         string
	TTL          string
}

var defaultConfig = Config{
	Service: "",
	Consul:  defaultConsul,
	DNS:     defaultDNS,
}

var defaultConsul = Consul{
	Addr:     "localhost:8500",
	Scheme:   "http",
	CAFile:   "",
	CertFile: "",
	KeyFile:  "",
}

var defaultDNS = DNS{
	Backend: "",
	Porkbun: defaultPorkbun,
}

var defaultPorkbun = Porkbun{
	APIKey:       "",
	SecretAPIKey: "",
	Domain:       "",
	Name:         "",
	TTL:          "300",
}
