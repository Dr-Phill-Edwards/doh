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
	RR = map[string]uint16{"A": 1, "NS": 2, "MX": 15, "SOA": 6, "TXT": 16}
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

func DecodeDomain(domain []byte, query string) (string, int) {
	name := ""
	index := 0
	for {
		if domain[index] == 0xC0 {
			name += query
			index += 2
			break
		}
		len := domain[index]
		index++
		if len == 0 {
			break
		}
		for ; len > 0; len-- {
			name += string(domain[index])
			index++
		}
		name += "."
	}
	return name, index
}

func FromBytes(response []byte) {
	var d DoH
	d.header.id = binary.BigEndian.Uint16(response)
	d.header.flags = binary.BigEndian.Uint16(response[2:])
	d.header.qdcount = binary.BigEndian.Uint16(response[4:])
	d.header.ancount = binary.BigEndian.Uint16(response[6:])
	d.header.nscount = binary.BigEndian.Uint16(response[8:])
	d.header.arcount = binary.BigEndian.Uint16(response[10:])

	question := response[12:]
	domain, index := DecodeDomain(question, "")
	d.question.name = append(d.question.name, question[:index]...)
	d.question.querytype = binary.BigEndian.Uint16(question[index:])
	d.question.queryclass = binary.BigEndian.Uint16(question[index+2:])

	answer := question[index+4:]
	for n := 0; n < int(d.header.ancount); n++ {
		index = DecodeAnswer(answer, domain)
		answer = answer[index:]
	}
}

func DecodeAnswer(answer []byte, domain string) int {
	reply, index := DecodeDomain(answer, domain)
	querytype := binary.BigEndian.Uint16(answer[index:])
	index += 2
	queryclass := int(binary.BigEndian.Uint16(answer[index:]))
	index += 4
	reply += " " + strconv.Itoa(int(binary.BigEndian.Uint16(answer[index:])))
	if queryclass == 1 {
		reply += " IN"
	}
	for key, code := range RR {
		if code == querytype {
			reply += " " + key + " "
		}
	}
	index += 2
	len := int(binary.BigEndian.Uint16(answer[index:]))
	index += 2
	if querytype == 1 {
		reply += strconv.Itoa(int(answer[index])) + "."
		reply += strconv.Itoa(int(answer[index+1])) + "."
		reply += strconv.Itoa(int(answer[index+2])) + "."
		reply += strconv.Itoa(int(answer[index+3]))
		index += 4
	} else if querytype == 2 {
		value, n := DecodeDomain(answer[index:], domain)
		reply += value
		index += n
	} else if querytype == 16 {
		reply += string(answer[index : index+len])
		index += len
	} else if querytype == 15 {
		reply += strconv.Itoa(int(binary.BigEndian.Uint16(answer[index:]))) + " "
		index += 2
		value, n := DecodeDomain(answer[index:], domain)
		reply += value
		index += n
	}
	fmt.Println(reply)
	return index
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
	return strings.TrimRight(string(b), "=")
}
