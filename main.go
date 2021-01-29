package main

import (
	"doh/doh"
	"fmt"
	"io/ioutil"
	"net/http"
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
	reply := doh.FromBytes(body)
	fmt.Println(reply)
}

func main() {
	d := doh.New()
	doh.SetQuestion(&d, "A", "www.example.com")
	Request(&d)
}
