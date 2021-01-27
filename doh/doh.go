package doh

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type doh struct {
	rr         map[string]byte
	header     []byte
	question   []byte
	querytype  []byte
	queryclass []byte
	query      []byte
}

func New() doh {
	var h []byte
	h = append(h, 0, 0)
	h = append(h, 1, 0)
	h = append(h, 0, 1)
	h = append(h, 0, 0)
	h = append(h, 0, 0)
	h = append(h, 0, 0)
	d := doh{map[string]byte{"A": 1, "NS": 2, "MX": 24, "SOA": 6, "TXT": 16}, h, []byte{}, []byte{0, 1}, []byte{0, 1}, []byte{}}
	return d
}

func Question(d *doh, rr string, domain string) {
	d.querytype[1] = d.rr[rr]
	d.question = d.question[:0]
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		d.question = append(d.question, byte(len(part)))
		for i := 0; i < len(part); i++ {
			d.question = append(d.question, part[i])
		}
	}
	d.question = append(d.question, 0)
	d.query = append(d.header, d.question...)
	d.query = append(d.query, d.querytype...)
	d.query = append(d.query, d.queryclass...)
}

func Print(d *doh) {
	for i := 0; i < len(d.query); i++ {
		fmt.Printf("%02x ", d.query[i])
	}
	fmt.Println()
}

func Encode(d *doh) string {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(d.query)))
	base64.StdEncoding.Encode(b, d.query)
	return string(b)
}
