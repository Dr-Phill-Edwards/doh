package main

import (
	"doh/doh"
	"fmt"
)

func main() {
	d := doh.New()
	doh.SetQuestion(&d, "A", "www.example.com")
	doh.Print(&d)
	fmt.Println(doh.Encode(&d))
}
