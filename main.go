package main

import (
	"congress/dns"
	"congress/logger"
	"congress/models"
	"congress/proxy"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func e(err error, format string) bool {
	if err != nil {
		logger.Default.Error(format, err)
		return false
	} else {
		return true
	}
}

func main() {
	var config models.Config
	if bytes, err := ioutil.ReadFile("congress.yamll"); e(err, "Failed to read file: %s") {
		if err := yaml.Unmarshal(bytes, &config); e(err, "Failed to parse yaml: %s") {
			dns := dns.New()
			http := proxy.New()

			for _, rule := range config.Rules {
				dns.AddRecord(rule.Host, "127.0.0.1")
				host := proxy.NewHost(rule.Host, rule.DefaultBackend)
				for _, path := range rule.Paths {
					host.AddPath(path.Path, proxy.PrefixPath, path.Backend)
				}
				http.AddHost(host)
			}

			if config.Congress.Dns.Enabled {
				go dns.Listen(int(config.Congress.Dns.Port))
			}
			if config.Congress.Proxy.Enabled {
				go http.Listen(int(config.Congress.Proxy.Port))
			}

			<-make(chan bool)
		}
	}
}
