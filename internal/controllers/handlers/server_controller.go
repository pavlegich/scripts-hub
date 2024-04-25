// Package handlers contains server controller object and
// methods for building the server route, command functions
// for activating the command handler in controller
// and commands handlers.
package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/pavlegich/scripts-hub/internal/controllers/middlewares"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	db  *sql.DB
	cfg *config.Config
}

// NewController creates and returns new server controller.
func NewController(ctx context.Context, db *sql.DB, cfg *config.Config) *Controller {
	return &Controller{
		db:  db,
		cfg: cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (*http.Handler, error) {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	Activate(ctx, router, c.cfg, c.db)

	handler := middlewares.Recovery(router)
	handler = middlewares.WithLogging(handler)

	return &handler, nil
}
