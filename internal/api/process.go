package api

import (
	"encoding/json"
	"net/http"

	"github.com/danmrichards/firecracker-wrapper/internal/utils"
)

// StartProcessRequest requests a process to be started with the given args.
type StartProcessRequest struct {
	Exe  string
	Args []string
}

func (h *handler) startProcessHandler(w http.ResponseWriter, r *http.Request) {
	var startRequest StartProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&startRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if startRequest.Exe == "" {
		http.Error(w, "missing exe", http.StatusBadRequest)
		return
	}

	go h.startProcess(startRequest.Exe, startRequest.Args)

	w.Write([]byte("started"))
}

func (h *handler) startProcess(exe string, args []string) {
	out, err := utils.ExecCommand(exe, args...)
	if err != nil {
		h.logger.Error(err)
		return
	}

	h.logger.Info(out)
}
