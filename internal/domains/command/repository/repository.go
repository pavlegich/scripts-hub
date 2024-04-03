// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"

	"github.com/pavlegich/scripts-hub/internal/domains/command"
)

// Repository contains storage objects for storing the commands.
type Repository struct {
	db *sql.DB
}

// NewCommandRepository returns new commands repository object.
func NewCommandRepository(ctx context.Context, db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CreateCommand stores new command into the storage.
func (r *Repository) CreateCommand(ctx context.Context, command *command.Command) error {
	return nil
}

// GetAllCommands gets and returns all the commands from the storage.
func (r *Repository) GetAllCommands(ctx context.Context) ([]*command.Command, error) {
	return nil, nil
}

// GetCommandByName gets and returns the requested by name command from the storage.
func (r *Repository) GetCommandByName(ctx context.Context, name string) (*command.Command, error) {
	return nil, nil
}
