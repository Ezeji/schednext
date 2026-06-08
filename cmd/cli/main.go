package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"schednext/internal/model"
)

const socketPath = "/run/schednext/schednext.sock"

func main() {
	if len(os.Args) < 3 {
		log.Println("usage: schednext <start|stop> <job-id>")
		os.Exit(1)
	}

	req := model.IPCRequest{
		Action: os.Args[1],
		User:   "schednext-runtime",
		JobID:  os.Args[2],
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Println("agent not running")
		os.Exit(1)
	}
	defer conn.Close()

	json.NewEncoder(conn).Encode(req)

	var resp model.IPCResponse
	json.NewDecoder(conn).Decode(&resp)

	if resp.OK {
		log.Println("✔ ", resp.Message)
	} else {
		log.Println("✖ ", resp.Message)
		os.Exit(1)
	}
}
