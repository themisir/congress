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

package dns

import (
	"congress/logger"
	"fmt"

	"github.com/miekg/dns"
)

type DnsServer struct {
	records  map[string]string
	Fallback *DnsResolver
}

var log = logger.Default.Child("dns")

func New() *DnsServer {
	return &DnsServer{records: make(map[string]string)}
}

func (s *DnsServer) AddRecord(domain string, target string) {
	s.records[domain+"."] = target
}

func (s *DnsServer) Print() {
	for host, ip := range s.records {
		log.Debug("%s A %s", host, ip)
	}
}

func (s *DnsServer) Listen(port uint) error {
	log.Info("Listening at port %d", port)
	return dns.ListenAndServe(fmt.Sprintf(":%d", port), "udp", s)
}

func (s *DnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if r.Opcode == dns.OpcodeQuery && len(r.Question) > 0 {
		question := r.Question[0]

		log.Debug("---> %s", question.Name)

		if question.Qtype == dns.TypeA {
			if ip, ok := s.records[question.Name]; ok {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", question.Name, ip))
				if err != nil {
					log.Error("Failed to create dns response: %s", err)
				} else {
					msg := new(dns.Msg)
					msg.Answer = []dns.RR{rr}
					msg.SetReply(r)
					w.WriteMsg(msg)
					return
				}
			}
		}

		if s.Fallback != nil {
			log.Debug("Resolving %s using fallback server", question.Name)

			res, err := s.Fallback.Resolve(question)

			if err != nil {
				log.Error("Failed to resolve DNS request using fallback server: %s", err)
			} else {
				msg := new(dns.Msg)
				msg.Answer = res.Answer
				msg.SetReply(r)
				w.WriteMsg(msg)
			}
		}
	}
}
