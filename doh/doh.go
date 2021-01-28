package doh

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
)

var RR map[string]int16
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
	name       []byte
	querytype  int16
	queryclass int16
}

type DoH struct {
	header   Header
	question Question
}

func init() {
	RR = map[string]int16{"A": 1, "NS": 2, "MX": 24, "SOA": 6, "TXT": 16}
}

func New() DoH {
	header := Header{0, 0x0100, 1, 0, 0, 0}
	question := Question{[]byte{}, 1, 1}
	d := DoH{header, question}
	return d
}

func SetQuestion(d *DoH, rr string, domain string) {
	d.question.querytype = RR[rr]
	d.question.name = d.question.name[:0]
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		d.question.name = append(d.question.name, byte(len(part)))
		for i := 0; i < len(part); i++ {
			d.question.name = append(d.question.name, part[i])
		}
	}
	d.question.name = append(d.question.name, 0)

	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, d.header)
	if err != nil {
		fmt.Println(err)
	}
	Query = buffer.Bytes()
	Query = append(Query, d.question.name...)
	Query = append(Query, 0, byte(d.question.querytype), 0, byte(d.question.queryclass))
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
