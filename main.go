package main

import (
	"io/ioutil"
	"log"

	"github.com/noelbundick/azssh/cmd"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	cmd.Execute()
}
