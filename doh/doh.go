package doh

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

var RR map[string]uint16

type Header struct {
	id      uint16
	flags   uint16
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

type Question struct {
	name       []byte
	querytype  uint16
	queryclass uint16
}

type Answer struct {
	name       []byte
	querytype  uint16
	queryclass uint16
	ttl        uint16
	rdata      []byte
}

type DoH struct {
	header   Header
	question Question
	answer   Answer
}

func init() {
	RR = map[string]uint16{"A": 1, "NS": 2, "MX": 24, "SOA": 6, "TXT": 16}
}

func New() DoH {
	header := Header{0, 0x0100, 1, 0, 0, 0}
	question := Question{[]byte{}, 1, 1}
	d := DoH{header, question, Answer{}}
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
}

func ToBytes(d *DoH) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, d.header)
	if err != nil {
		fmt.Println(err)
	}
	query := buffer.Bytes()
	query = append(query, d.question.name...)
	query = append(query, 0, byte(d.question.querytype), 0, byte(d.question.queryclass))
	return query
}

func FromBytes(response []byte) string {
	var d DoH
	d.header.id = binary.BigEndian.Uint16(response)
	d.header.flags = binary.BigEndian.Uint16(response[2:])
	d.header.qdcount = binary.BigEndian.Uint16(response[4:])
	d.header.ancount = binary.BigEndian.Uint16(response[6:])
	d.header.nscount = binary.BigEndian.Uint16(response[8:])
	d.header.arcount = binary.BigEndian.Uint16(response[10:])

	question := response[12:]
	reply := ""
	index := 0
	for {
		len := question[index]
		index++
		if len == 0 {
			break
		}
		for ; len > 0; len-- {
			reply += string(question[index])
			index++
		}
		reply += "."
	}
	d.question.name = append(d.question.name, question[:index]...)
	d.question.querytype = binary.BigEndian.Uint16(question[index:])
	d.question.queryclass = binary.BigEndian.Uint16(question[index+2:])

	answer := question[index+4:]
	fmt.Println(answer)
	if answer[0] == 0xC0 {
		index = 2
	}
	reply += " " + strconv.Itoa(int(binary.BigEndian.Uint16(question[index+7:])))
	if answer[index+3] == 1 {
		reply += " IN"
	}
	for key, code := range RR {
		if byte(code) == answer[index+1] {
			reply += " " + key + " "
		}
	}
	if answer[index+1] == 1 {
		reply += strconv.Itoa(int(answer[index+10])) + "."
		reply += strconv.Itoa(int(answer[index+11])) + "."
		reply += strconv.Itoa(int(answer[index+12])) + "."
		reply += strconv.Itoa(int(answer[index+13]))
	}
	return reply
}

func Print(d *DoH) {
	query := ToBytes(d)
	for i := 0; i < len(query); i++ {
		fmt.Printf("%02x ", query[i])
	}
	fmt.Println()
}

func Encode(d *DoH) string {
	query := ToBytes(d)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(query)))
	base64.StdEncoding.Encode(b, query)
	return string(b)
}
