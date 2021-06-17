package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/magiconair/properties"
)

func Load(args, environ []string) (cfg *Config, err error) {
	var props *properties.Properties

	cmdline, path, version, err := parse(args)
	switch {
	case err != nil:
		return nil, err
	case version:
		return nil, nil
	case path != "":
		switch {
		case strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"):
			props, err = properties.LoadURL(path)
		case path != "":
			props, err = properties.LoadFile(path, properties.UTF8)
		}
		if err != nil {
			return nil, err
		}
	}
	envprefix := []string{"FABIO_", ""}
	return load(cmdline, environ, envprefix, props)
}

var errInvalidConfig = errors.New("invalid or missing path to config file")

// parse extracts the version and config file flags from the command
// line arguments and returns the individual parts. Test flags are
// ignored.
func parse(args []string) (cmdline []string, path string, version bool, err error) {
	if len(args) < 1 {
		panic("missing exec name")
	}

	// always copy the name of the executable
	cmdline = args[:1]

	// parse rest of the arguments
	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch {
		// version flag
		case arg == "-v" || arg == "-version" || arg == "--version":
			return nil, "", true, nil

		// config file without '='
		case arg == "-cfg" || arg == "--cfg":
			if i >= len(args)-1 {
				return nil, "", false, errInvalidConfig
			}
			path = args[i+1]
			i++

		// config file with '='. needs unquoting
		case strings.HasPrefix(arg, "-cfg=") || strings.HasPrefix(arg, "--cfg="):
			if strings.HasPrefix(arg, "-cfg=") {
				path = arg[len("-cfg="):]
			} else {
				path = arg[len("--cfg="):]
			}
			switch {
			case path == "":
				return nil, "", false, errInvalidConfig
			case path[0] == '\'':
				path = strings.Trim(path, "'")
			case path[0] == '"':
				path = strings.Trim(path, "\"")
			}
			if path == "" {
				return nil, "", false, errInvalidConfig
			}

		// ignore test flags
		case strings.HasPrefix(arg, "-test."):
			continue

		default:
			cmdline = append(cmdline, arg)
		}
	}
	return cmdline, path, false, nil
}

func load(cmdline, environ, envprefix []string, props *properties.Properties) (cfg *Config, err error) {
	cfg = &Config{}
	f := NewFlagSet(cmdline[0], flag.ExitOnError)

	// dummy values which were parsed earlier
	f.String("cfg", "", "Path or URL to config file")
	f.Bool("v", false, "Show version")
	f.Bool("version", false, "Show version")

	f.StringVar(&cfg.Service, "service", defaultConfig.Service, "Consul service to track (Generally the load balancer)")
	f.StringVar(&cfg.Consul.Addr, "consul.addr", defaultConfig.Consul.Addr, "Address of Consul agent")
	f.StringVar(&cfg.Consul.Scheme, "consul.scheme", defaultConfig.Consul.Scheme, "Scheme of Consul agent (http/https)")
	f.StringVar(&cfg.DNS.Backend, "dns.backend", defaultConfig.DNS.Backend, "Name of DNS backend to use")
	f.StringVar(&cfg.DNS.Porkbun.APIKey, "dns.porkbun.apikey", defaultConfig.DNS.Porkbun.APIKey, "APIKey to connect to porkbun")
	f.StringVar(&cfg.DNS.Porkbun.SecretAPIKey, "dns.porkbun.secretapikey", defaultConfig.DNS.Porkbun.SecretAPIKey, "SecretAPIKey to connect to porkbun")
	f.StringVar(&cfg.DNS.Porkbun.Domain, "dns.porkbun.domain", defaultPorkbun.Domain, "Porkbun domain to be updated")
	f.StringVar(&cfg.DNS.Porkbun.Name, "dns.porkbun.name", defaultPorkbun.Name, "Porkbun subdmain domain to be updated \nUse * for wildcard record \nUnset or \"\" for root domain")

	// Parse configuration
	if err := f.ParseFlags(cmdline[1:], environ, envprefix, props); err != nil {
		return nil, err
	}

	// A service must be provided
	if cfg.Service == "" {
		return nil, fmt.Errorf("No service was provided to track")
	}
	cfg.DNS.Backend = strings.ToLower(cfg.DNS.Backend)

	return cfg, nil
}
