package command

import (
	"context"
	"fmt"

	"github.com/pavlegich/scripts-hub/internal/entities"
	repo "github.com/pavlegich/scripts-hub/internal/repository"
)

// Service describes methods for communication between
// handlers and repositories.
//
//go:generate mockgen -destination=../../mocks/mock_Service.go -package=mocks github.com/pavlegich/scripts-hub/internal/service/command Service
type Service interface {
	Create(ctx context.Context, command *entities.Command) (int, error)
	List(ctx context.Context) ([]*entities.Command, error)
	Unload(ctx context.Context, name string) (*entities.Command, error)
	Update(ctx context.Context, command *entities.Command) error
	Delete(ctx context.Context, name string) error
}

// CommandService contains objects for command service.
type CommandService struct {
	repo repo.Repository
}

// NewCommandService returns new command service.
func NewCommandService(ctx context.Context, repo repo.Repository) *CommandService {
	return &CommandService{
		repo: repo,
	}
}

// Create creates new requested command and requests repository to put it into the storage.
func (s *CommandService) Create(ctx context.Context, c *entities.Command) (int, error) {
	command, err := s.repo.CreateCommand(ctx, c)
	if err != nil {
		return -1, fmt.Errorf("Create: create command failed %w", err)
	}

	return command.ID, nil
}

// List returns list of available commands stored in the database.
func (s *CommandService) List(ctx context.Context) ([]*entities.Command, error) {
	cmds, err := s.repo.GetAllCommands(ctx)
	if err != nil {
		return nil, fmt.Errorf("List: get commands list failed %w", err)
	}

	return cmds, nil
}

// Unload gets command by command's name and returns it.
func (s *CommandService) Unload(ctx context.Context, name string) (*entities.Command, error) {
	cmd, err := s.repo.GetCommandByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("Unload: get command failed %w", err)
	}

	return cmd, nil
}

// Update updates command by command's name.
func (s *CommandService) Update(ctx context.Context, c *entities.Command) error {
	err := s.repo.UpdateCommandByName(ctx, c)
	if err != nil {
		return fmt.Errorf("Update: update command failed %w", err)
	}

	return nil
}

// Delete delete command from the DB and returns it.
func (s *CommandService) Delete(ctx context.Context, name string) error {
	err := s.repo.DeleteCommandByName(ctx, name)
	if err != nil {
		return fmt.Errorf("Delete: delete command failed %w", err)
	}

	return nil
}
