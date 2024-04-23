package command

import (
	"context"
	"fmt"
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
func (s *CommandService) Create(ctx context.Context, command *Command) (int, error) {
	command, err := s.repo.CreateCommand(ctx, command)
	if err != nil {
		return -1, fmt.Errorf("Create: create command failed %w", err)
	}

	return command.ID, nil
}

// List returns list of available commands stored in the database.
func (s *CommandService) List(ctx context.Context) ([]*Command, error) {
	cmds, err := s.repo.GetAllCommands(ctx)
	if err != nil {
		return nil, fmt.Errorf("List: get commands list failed %w", err)
	}

	return cmds, nil
}

// Unload gets command by command's name and returns it.
func (s *CommandService) Unload(ctx context.Context, name string) (*Command, error) {
	cmd, err := s.repo.GetCommandByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("Unload: get command failed %w", err)
	}

	return cmd, nil
}
