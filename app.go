/*
	Copyright 2021 Misir Jafarov

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

			http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"congress/dns"
	"congress/logger"
	"congress/models"
	"congress/proxy"
)

type App struct {
	config *models.Config
}

func (a *App) startDns() {
	server := dns.New()
	config := a.config.Congress.Dns

	for _, rule := range a.config.Rules {
		server.AddRecord(rule.Host, a.config.Congress.Ip)
	}

	if config.Fallback != "" {
		server.Fallback = dns.NewResolver(config.Fallback)
	}

	server.Listen(config.Port)
}

func (a *App) startProxy() {
	server := proxy.New()
	config := a.config.Congress.Proxy

	for _, rule := range a.config.Rules {
		host := proxy.NewHost(rule.Host, rule.DefaultBackend)
		for _, path := range rule.Paths {
			host.AddPath(path.Path, proxy.PrefixPath, path.Backend)
		}
		server.AddHost(host)
	}

	if config.TLS != nil {
		go func() {
			if err := server.ListenTLS(config.TLS.Port, config.TLS.CertFile, config.TLS.KeyFile); err != nil {
				logger.Default.Error("Failed to start proxy server with TLS: %s", err)
			}
		}()
	}

	server.Listen(config.Port)
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
