package dns

import "github.com/miekg/dns"

type DnsResolver struct {
	server string
}

func NewResolver(s string) *DnsResolver {
	return &DnsResolver{s}
}

func (d *DnsResolver) Resolve(q dns.Question) (*dns.Msg, error) {
	msg := &dns.Msg{Question: []dns.Question{q}}
	return dns.Exchange(msg, d.server)
}
