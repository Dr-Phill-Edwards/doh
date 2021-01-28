package doh

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
)

var RR map[string]byte
var Query []byte

type Header struct {
	id      int16
	flags   uint16
	qdcount int16
	ancount int16
	nscount int16
	arcount int16
}

type Question struct {
}

type DoH struct {
	header     []byte
	question   []byte
	querytype  []byte
	queryclass []byte
}

func init() {
	RR = map[string]byte{"A": 1, "NS": 2, "MX": 24, "SOA": 6, "TXT": 16}
}

func New() DoH {
	header := Header{0, 0x0100, 1, 0, 0, 0}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, header)
	fmt.Println(buffer.Bytes())
	var h []byte
	h = append(h, 0, 0)
	h = append(h, 1, 0)
	h = append(h, 0, 1)
	h = append(h, 0, 0)
	h = append(h, 0, 0)
	h = append(h, 0, 0)
	d := DoH{h, []byte{}, []byte{0, 1}, []byte{0, 1}}
	return d
}

func SetQuestion(d *DoH, rr string, domain string) {
	d.querytype[1] = RR[rr]
	d.question = d.question[:0]
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		d.question = append(d.question, byte(len(part)))
		for i := 0; i < len(part); i++ {
			d.question = append(d.question, part[i])
		}
	}
	d.question = append(d.question, 0)
	Query = append(d.header, d.question...)
	Query = append(Query, d.querytype...)
	Query = append(Query, d.queryclass...)
}

func Print(d *DoH) {
	for i := 0; i < len(Query); i++ {
		fmt.Printf("%02x ", Query[i])
	}
	fmt.Println()
}

func Encode(d *DoH) string {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(Query)))
	base64.StdEncoding.Encode(b, Query)
	return string(b)
}
