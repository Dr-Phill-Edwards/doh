package main

import (
	"doh/doh"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Request(d *doh.DoH) {
	uri := "https://dns.google/dns-query?dns=" + doh.Encode(d)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Add("accept", "application/dns-message")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error " + err.Error())
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	doh.FromBytes(body)
}

func main() {
	if len(os.Args) != 3 || doh.RR[os.Args[1]] == 0 {
		fmt.Println("Usage: " + os.Args[0] + " A|NS|MX|SOA|TXT url")
		os.Exit(1)
	}
	d := doh.New()
	doh.SetQuestion(&d, os.Args[1], os.Args[2])
	Request(&d)
}
