package main

import (
	"doh/doh"
	"fmt"
)

func main() {
	d := doh.New()
	doh.Question(&d, "A", "www.example.com")
	doh.Print(&d)
	fmt.Println(doh.Encode(&d))
}
