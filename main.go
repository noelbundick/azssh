package main

import (
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	url := Provision()
	Connect(url)
}
