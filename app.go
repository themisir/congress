package main

import (
	"congress/dns"
	"congress/models"
	"congress/proxy"
)

type App struct {
	config *models.Config
}

func (a *App) startDns() {
	server := dns.New()
	for _, rule := range a.config.Rules {
		server.AddRecord(rule.Host, a.config.Congress.Ip)
	}
	if a.config.Congress.Dns.Fallback != "" {
		server.Fallback = dns.NewResolver(a.config.Congress.Dns.Fallback)
	}
	server.Listen(a.config.Congress.Dns.Port)
}

func (a *App) startProxy() {
	server := proxy.New()
	for _, rule := range a.config.Rules {
		host := proxy.NewHost(rule.Host, rule.DefaultBackend)
		for _, path := range rule.Paths {
			host.AddPath(path.Path, proxy.PrefixPath, path.Backend)
		}
		server.AddHost(host)
	}
	server.Listen(a.config.Congress.Proxy.Port)
}

func (a *App) Run() {
	if a.config.Congress.Dns.Enabled {
		go a.startDns()
	}
	if a.config.Congress.Proxy.Enabled {
		go a.startProxy()
	}

	if a.config.Congress.Dns.Enabled || a.config.Congress.Proxy.Enabled {
		<-make(chan bool) // dead loop
	}
}
