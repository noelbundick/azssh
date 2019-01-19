package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
)

func dial(url string) *websocket.Conn {
	log.Println("connect:", url)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

// read from ws
func pumpOutput(c *websocket.Conn, w io.Writer) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		w.Write(message)
	}
}

// send to ws
func pumpInput(c *websocket.Conn, r io.Reader, done chan interface{}) {
	data := make([]byte, 1)
	for {
		r.Read(data)

		err := c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("send:", err)
		}
	}
}

// ConnectToWebsocket wires up STDIN and STDOUT to a websocket, allowing you to use it as a terminal
func ConnectToWebsocket(url string) {
	// disable input buffering
	// do not display entered characters on the screen
	if runtime.GOOS == "linux" {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	} else if runtime.GOOS == "darwin" {
		exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	}

	done := make(chan interface{})

	c := dial(url)
	defer c.Close()

	go pumpOutput(c, os.Stdout)
	go pumpInput(c, os.Stdin, done)

	<-done
}
