package main

import (
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	token := GetToken()
	url := ProvisionCloudShell(token)
	ConnectToWebsocket(url)
}
