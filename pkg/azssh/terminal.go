package azssh

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh/terminal"
)

var sigwinch = make(chan os.Signal, 1)

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

func GetTerminalSize() TerminalSize {
	cols, rows, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		cols = 80
		rows = 30
	}
	return TerminalSize{
		Rows: rows,
		Cols: cols,
	}
}

func pumpSigwinch(resize chan<- TerminalSize) {
	for {
		<-sigwinch
		newSize := GetTerminalSize()
		resize <- newSize
	}
}

// ConnectToWebsocket wires up STDIN and STDOUT to a websocket, allowing you to use it as a terminal
func ConnectToWebsocket(url string, resize chan<- TerminalSize) {
	// disable input buffering
	// do not display entered characters on the screen
	if runtime.GOOS == "linux" {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	} else if runtime.GOOS == "darwin" {
		exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	}

	// hook into terminal resizes
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		signal.Notify(sigwinch, syscall.SIGWINCH)
		go pumpSigwinch(resize)
	}

	done := make(chan interface{})

	c := dial(url)
	defer c.Close()

	go pumpOutput(c, os.Stdout)
	go pumpInput(c, os.Stdin, done)

	<-done
}
