// Package http contains object of command functions
// for activating the command handler in controller
// and commands handlers.
package http

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/pavlegich/scripts-hub/internal/domains/command"
	repo "github.com/pavlegich/scripts-hub/internal/domains/command/repository"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
)

// CommandHandler contains objects for work with command handlers.
type CommandHandler struct {
	Config  *config.Config
	Service command.Service
}

// Activate activates handler for command object.
func Activate(ctx context.Context, r *http.ServeMux, cfg *config.Config, db *sql.DB) {
	s := command.NewCommandService(ctx, repo.NewCommandRepository(ctx, db))
	newHandler(r, cfg, s)
}

// newHandler initializes handler for command object.
func newHandler(r *http.ServeMux, cfg *config.Config, s command.Service) {
	_ = &CommandHandler{
		Config:  cfg,
		Service: s,
	}

	// r.HandleFunc("/api/command", h.HandleCommand)
	// r.HandleFunc("/api/commands", h.HandleCommands)
}
