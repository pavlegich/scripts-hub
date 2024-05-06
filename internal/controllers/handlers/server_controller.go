// Package handlers contains server controller object and
// methods for building the server route, command functions
// for activating the command handler in controller
// and command handlers.
package handlers

import (
	"context"
	"net/http"

	"github.com/pavlegich/scripts-hub/internal/controllers/middlewares"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
	"github.com/pavlegich/scripts-hub/internal/repository"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	repo repository.Repository
	cfg  *config.Config
}

// NewController creates and returns new server controller.
func NewController(ctx context.Context, repo repository.Repository, cfg *config.Config) *Controller {
	return &Controller{
		repo: repo,
		cfg:  cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (*http.Handler, error) {
	router := http.NewServeMux()

	commandsActivate(ctx, router, c.repo, c.cfg)

	handler := middlewares.Recovery(router)
	handler = middlewares.WithLogging(handler)

	return &handler, nil
}
