package command

import (
	"context"
)

// CommandService contains objects for command service.
type CommandService struct {
	repo Repository
}

// NewCommandService returns new command service.
func NewCommandService(ctx context.Context, repo Repository) *CommandService {
	return &CommandService{
		repo: repo,
	}
}

// Create creates new requested command and requests repository to put it into the storage.
func (s *CommandService) Create(ctx context.Context, command *Command) error {
	return nil
}

// List returns list of available commands stored in the database.
func (s *CommandService) List(ctx context.Context) ([]*Command, error) {
	return nil, nil
}

// Unload gets command by command's name and returns it.
func (s *CommandService) Unload(ctx context.Context, name string) (*Command, error) {
	return nil, nil
}
