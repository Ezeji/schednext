package agent

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"

	"schednext/internal/util"
	"schednext/internal/model"
)

const socketPath = "/run/schednext/schednext.sock"

var ipcListener net.Listener

func StartIPCServer(optPath string) {
	os.MkdirAll(filepath.Dir(socketPath), 0755)
	os.Remove(socketPath)

	ipcListener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	log.Println("IPC listening on", socketPath)

	go func() {
		for {
			conn, err := ipcListener.Accept()
			if err != nil {
				continue
			}
			go handleConn(conn, optPath)
		}
	}()
}

func ShutdownIPC() {

	if ipcListener != nil {
		ipcListener.Close()
	}

	os.Remove(socketPath)
}

func handleConn(conn net.Conn, optPath string) {
	defer conn.Close()

	var req model.IPCRequest
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		return
	}

	resp := handleRequest(req, optPath)
	json.NewEncoder(conn).Encode(resp)
}

func handleRequest(req model.IPCRequest, optPath string) model.IPCResponse {
	cfgPath := filepath.Join(optPath, req.User, configName)

	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0644)
	if err != nil {
		return model.IPCResponse{false, "config not found"}
	}
	defer f.Close()

	if err := lockFile(f); err != nil {
		return model.IPCResponse{false, "config busy"}
	}
	defer unlockFile(f)

	var cfg model.Config
	if err := util.ReadConfig(cfgPath, &cfg); err != nil {
		return model.IPCResponse{false, "invalid config"}
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
				return model.IPCResponse{false, "unknown action"}
			}

			util.WriteConfigAtomic(cfgPath, &cfg)
			return model.IPCResponse{true, result}
		}
	}

	return model.IPCResponse{false, "job not found"}
}
