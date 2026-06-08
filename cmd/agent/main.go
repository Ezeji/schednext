package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"schednext/internal/agent"
	"schednext/internal/statelens"

	"bazil.org/fuse"
)

func main() {

	conn, err := statelens.Mount("/statelens")
	if err != nil {
		log.Fatal(err)
	}

	go agent.RunAgent("/opt")

	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-sigChan

	log.Printf("received %v signal", sig)

	agent.Shutdown()
	agent.ShutdownIPC()

	fuse.Unmount("/statelens")

	if conn != nil {
		conn.Close()
	}

	log.Println("shutdown complete")
}