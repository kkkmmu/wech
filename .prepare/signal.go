package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sich := make(chan os.Signal)
	signal.Notify(sich, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGINT)

	for s := range sich {
		fmt.Println("Received signale")
		fmt.Printf("%#v\n", s)
		fmt.Printf("%s\n", s)
		os.Exit(0)
	}
}
