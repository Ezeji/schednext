package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

const socketPath = "/run/schednext/schednext.sock"

type Request struct {
	Action string `json:"action"`
	User   string `json:"user"`
	JobID  string `json:"jobId"`
}

type Response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) < 4 {
		log.Println("usage: schednext <start|stop> <user> <job-id>")
		os.Exit(1)
	}

	req := Request{
		Action: os.Args[1],
		User:   os.Args[2],
		JobID:  os.Args[3],
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Println("agent not running")
		os.Exit(1)
	}
	defer conn.Close()

	json.NewEncoder(conn).Encode(req)

	var resp Response
	json.NewDecoder(conn).Decode(&resp)

	if resp.OK {
		log.Println("✔ ", resp.Message)
	} else {
		log.Println("✖ ", resp.Message)
		os.Exit(1)
	}
}
