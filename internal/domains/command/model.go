// Package command contains object and methods
// for interacting with commands.
package command

import "context"

// Command contains data for commands.
type Command struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Script string `json:"script"`
}

// Service describes methods for communication between
// handlers and repositories.
type Service interface {
	Create(ctx context.Context, command *Command) error
	List(ctx context.Context) ([]*Command, error)
	Unload(ctx context.Context, name string) (*Command, error)
}

// Repository describes methods related with commands
// for interaction with database.
type Repository interface {
	CreateCommand(ctx context.Context, command *Command) error
	GetAllCommands(ctx context.Context) ([]*Command, error)
	GetCommandByName(ctx context.Context, name string) (*Command, error)
}
