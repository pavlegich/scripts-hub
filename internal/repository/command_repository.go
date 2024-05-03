// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pavlegich/scripts-hub/internal/entities"
	errs "github.com/pavlegich/scripts-hub/internal/errors"
)

// Repository describes methods related with commands
// for interaction with database.
//
//go:generate mockgen -destination=../mocks/mock_Repository.go -package=mocks github.com/pavlegich/scripts-hub/internal/repository Repository
type Repository interface {
	CreateCommand(ctx context.Context, command *entities.Command) (*entities.Command, error)
	GetAllCommands(ctx context.Context) ([]*entities.Command, error)
	GetCommandByName(ctx context.Context, name string) (*entities.Command, error)
	AppendCommandOutputByName(ctx context.Context, command *entities.Command) error
	DeleteCommandByName(ctx context.Context, name string) error
}

// CommandRepository contains storage objects for storing the commands.
type CommandRepository struct {
	db *sql.DB
}

// NewCommandRepository returns new commands repository object.
func NewCommandRepository(ctx context.Context, db *sql.DB) *CommandRepository {
	return &CommandRepository{
		db: db,
	}
}

// CreateCommand stores new command into the storage.
func (r *CommandRepository) CreateCommand(ctx context.Context, c *entities.Command) (*entities.Command, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO commands (name, script) 
	VALUES ($1, $2) RETURNING id`, c.Name, c.Script)

	var id int
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, fmt.Errorf("CreateCommand: %w", errs.ErrCmdAlreadyExists)
		}

		return nil, fmt.Errorf("CreateCommand: scan row failed %w", err)
	}

	c.ID = id

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("CreateCommand: row.Err %w", err)
	}

	return c, nil
}

// GetAllCommands gets and returns all the commands from the storage.
func (r *CommandRepository) GetAllCommands(ctx context.Context) ([]*entities.Command, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, script, output FROM commands`)
	if err != nil {
		return nil, fmt.Errorf("GetAllCommands: read rows from table failed %w", err)
	}
	defer rows.Close()

	cmdsList := make([]*entities.Command, 0)
	for rows.Next() {
		var c entities.Command
		err = rows.Scan(&c.ID, &c.Name, &c.Script, &c.Output)
		if err != nil {
			return nil, fmt.Errorf("GetAllCommands: scan row failed %w", err)
		}
		cmdsList = append(cmdsList, &c)
	}

	if len(cmdsList) == 0 {
		return nil, fmt.Errorf("GetAllCommands: nothing to return %w", errs.ErrCmdNotFound)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetAllCommands: rows.Err %w", err)
	}

	return cmdsList, nil
}

// GetCommandByName gets and returns the requested by name command from the storage.
func (r *CommandRepository) GetCommandByName(ctx context.Context, name string) (*entities.Command, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, script, output FROM commands WHERE name = $1`, name)

	var c entities.Command
	err := row.Scan(&c.ID, &c.Name, &c.Script, &c.Output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("GetCommandByName: nothing to get, %w", errs.ErrCmdNotFound)
		}
		return nil, fmt.Errorf("GetCommandByName: scan row failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("GetCommandByName: row.Err %w", err)
	}

	return &c, nil
}

// AppendCommandOutputByName updates output for requested command in the storage.
func (r *CommandRepository) AppendCommandOutputByName(ctx context.Context, c *entities.Command) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("AppendCommandOutputByName: begin transaction failed %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `SELECT output FROM commands WHERE name = $1`, c.Name)
	var out string
	err = row.Scan(&out)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("AppendCommandOutputByName: nothing to get, %w", errs.ErrCmdNotFound)
		}
		return fmt.Errorf("AppendCommandOutputByName: scan row failed %w", err)
	}

	err = row.Err()
	if err != nil {
		return fmt.Errorf("AppendCommandOutputByName: row.Err %w", err)
	}

	c.Output = out + c.Output
	_, err = tx.ExecContext(ctx, `UPDATE commands SET output = $1 WHERE name = $2`, c.Output, c.Name)
	if err != nil {
		return fmt.Errorf("AppendCommandOutputByName: update command failed %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("AppendCommandOutputByName: commit transaction failed %w", err)
	}

	return nil
}

// DeleteCommandByName deletes command from the storage and returns it.
func (r *CommandRepository) DeleteCommandByName(ctx context.Context, name string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM commands WHERE name = $1`, name)

	if err != nil {
		return fmt.Errorf("DeleteCommandByName: delete command failed %w", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteCommandByName: couldn't get rows affected %w", err)
	}
	if rowsCount == 0 {
		return fmt.Errorf("DeleteCommandByName: nothing to delete, %w", errs.ErrCmdNotFound)
	}

	return nil
}
