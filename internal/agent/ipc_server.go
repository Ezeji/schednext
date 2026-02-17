package agent

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"

	"schednext/internal/util"
)

const socketPath = "/run/schednext/schednext.sock"

func StartIPCServer(homeRoot string) {
	os.MkdirAll(filepath.Dir(socketPath), 0755)
	os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	log.Println("IPC listening on", socketPath)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			go handleConn(conn, homeRoot)
		}
	}()
}

func handleConn(conn net.Conn, homeRoot string) {
	defer conn.Close()

	var req IPCRequest
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		return
	}

	resp := handleRequest(req, homeRoot)
	json.NewEncoder(conn).Encode(resp)
}

func handleRequest(req IPCRequest, homeRoot string) IPCResponse {
	cfgPath := filepath.Join(homeRoot, req.User, configName)

	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0644)
	if err != nil {
		return IPCResponse{false, "config not found"}
	}
	defer f.Close()

	if err := lockFile(f); err != nil {
		return IPCResponse{false, "config busy"}
	}
	defer unlockFile(f)

	var cfg Config
	if err := util.ReadConfig(cfgPath, &cfg); err != nil {
		return IPCResponse{false, "invalid config"}
	}

	var result string
	for i := range cfg.Jobs {
		if cfg.Jobs[i].ID == req.JobID {
			switch req.Action {
			case "start":
				cfg.Jobs[i].Enabled = true
				result = "Job started"
			case "stop":
				cfg.Jobs[i].Enabled = false
				result = "Job stopped"
			default:
				return IPCResponse{false, "unknown action"}
			}

			util.WriteConfigAtomic(cfgPath, &cfg)
			return IPCResponse{true, result}
		}
	}

	return IPCResponse{false, "job not found"}
}
