// +build !windows

package azssh

import (
	"os"
	"os/signal"
	"syscall"
)

func pumpResize(resize chan<- TerminalSize) {
	var sigwinch = make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	for {
		<-sigwinch
		newSize := GetTerminalSize()
		resize <- newSize
	}
}
