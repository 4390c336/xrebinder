package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var round int

func main() {
	round = 0
	dns.HandleFunc(".", Hanlder)
	server := &dns.Server{Addr: "", Net: "udp"}
	server.ListenAndServe()
}

func Hanlder(w dns.ResponseWriter, req *dns.Msg) {
	//not A rec
	if req.Question[0].Qtype != dns.TypeA {
		w.WriteMsg(req)
		return
	}

	//we can also do this based on the req.Id
	round = (round + 1) % 2

	sub := strings.Split(req.Question[0].Name, ".")

	if len(sub) < 2 {
		return
	}

	a := LongStr2Ip(sub[round])

	//logging
	fmt.Println(req.Id, req.Question[0].Name, w.RemoteAddr().String(), a.String())

	m := new(dns.Msg)
	m.SetReply(req)

	rec := &dns.A{
		Hdr: dns.RR_Header{
			Name:   req.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: a,
	}
	m.Answer = append(m.Answer, rec)
	w.WriteMsg(m)
}

func LongStr2Ip(a string) net.IP {
	ip, _ := strconv.Atoi(a)
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
