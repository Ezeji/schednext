package agent

type IPCRequest struct {
	Action string `json:"action"` // start | stop | status
	User   string `json:"user"`
	JobID  string `json:"jobId"`
}

type IPCResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}
