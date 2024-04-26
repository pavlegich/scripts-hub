package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pavlegich/scripts-hub/internal/entities"
	errs "github.com/pavlegich/scripts-hub/internal/errors"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
	"github.com/pavlegich/scripts-hub/internal/infra/logger"
	"github.com/pavlegich/scripts-hub/internal/repository"
	"github.com/pavlegich/scripts-hub/internal/service/command"
	"go.uber.org/zap"
)

// CommandHandler contains objects for work with command handlers.
type CommandHandler struct {
	Config  *config.Config
	Service command.Service
}

// commandsActivate activates handler for command object.
func commandsActivate(ctx context.Context, r *http.ServeMux, repo repository.Repository, cfg *config.Config) {
	s := command.NewCommandService(ctx, repo)
	newHandler(r, cfg, s)
}

// newHandler initializes handler for command object.
func newHandler(r *http.ServeMux, cfg *config.Config, s command.Service) {
	h := &CommandHandler{
		Config:  cfg,
		Service: s,
	}

	r.HandleFunc("/command", h.HandleCommand)
	r.HandleFunc("/commands", h.HandleCommands)
}

// HandleCommand handles request to create or get the command.
func (h *CommandHandler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateCommand(w, r)
	case http.MethodGet:
		h.HandleGetCommand(w, r)
	default:
		logger.Log.Error("HandleCommand: incorrect method",
			zap.String("method", r.Method))

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleCreateCommand handles request to create and execute new command.
func (h *CommandHandler) HandleCreateCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Error("HandleCreateCommand: incorrect method",
			zap.String("method", r.Method))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var req entities.Command
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleCreateCommand: read request body failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleCreateCommand: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bashCmd := strings.Split(req.Script, " ")
	cmd := exec.CommandContext(ctx, bashCmd[0], bashCmd[1:]...)
	err = cmd.Start()
	if err != nil {
		logger.Log.Error("HandleCreateCommand: start command failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commandID, err := h.Service.Create(ctx, &req)
	if err != nil {
		logger.Log.Error("HandleCreateCommand: create command failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrCmdAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"command_id": strconv.Itoa(commandID),
	}
	out, err := json.Marshal(resp)
	if err != nil {
		logger.Log.Error("HandleCreateCommand: marshal command to JSON failed",
			zap.Int("command_id", commandID),
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

// HandleGetCommand handles request to get the requested command.
func (h *CommandHandler) HandleGetCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req entities.Command
	want := map[string]struct{}{
		"name": {},
	}

	queries := r.URL.Query()
	for val := range queries {
		_, ok := want[val]
		if !ok {
			logger.Log.Error("HandleGetCommand: incorrect query",
				zap.String("query", val))

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(queries[val]) != 1 {
			logger.Log.Error("HandleGetCommand: incorrect number of queries",
				zap.Int("queries_count", len(queries)))

			w.WriteHeader(http.StatusBadRequest)
		}

		req.Name = queries[val][0]
	}

	command, err := h.Service.Unload(ctx, req.Name)
	if err != nil {
		logger.Log.With(zap.String("cmd_name", req.Name)).
			Error("HandleGetCommand: get command failed", zap.Error(err))

		if errors.Is(err, errs.ErrCmdNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmdJSON, err := json.Marshal(command)
	if err != nil {
		logger.Log.With(zap.String("cmd_name", req.Name)).
			Error("HandleGetCommand: marshal command failed", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cmdJSON)
}

// HandleCommands handles request to get list of commands.
func (h *CommandHandler) HandleCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Log.Error("HandleCommands: incorrect method",
			zap.String("method", r.Method))

		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	commands, err := h.Service.List(ctx)
	if err != nil {
		logger.Log.Error("HandleGetCommands: get commands list failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrCmdNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmdsJSON, err := json.Marshal(commands)
	if err != nil {
		logger.Log.Error("HandleCommands: marshal command failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cmdsJSON)
}
