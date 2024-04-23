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
	"github.com/pavlegich/scripts-hub/internal/domains/command"
	errs "github.com/pavlegich/scripts-hub/internal/errors"
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
func (r *Repository) CreateCommand(ctx context.Context, c *command.Command) (*command.Command, error) {
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
func (r *Repository) GetAllCommands(ctx context.Context) ([]*command.Command, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, script FROM commands`)
	if err != nil {
		return nil, fmt.Errorf("GetAllCommands: read rows from table failed %w", err)
	}
	defer rows.Close()

	cmdsList := make([]*command.Command, 0)
	for rows.Next() {
		var c command.Command
		err = rows.Scan(&c.ID, &c.Name, &c.Script)
		if err != nil {
			return nil, fmt.Errorf("GetAllCommands: scan row failed %w", err)
		}
		cmdsList = append(cmdsList, &c)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetAllCommands: rows.Err %w", err)
	}

	return cmdsList, nil
}

// GetCommandByName gets and returns the requested by name command from the storage.
func (r *Repository) GetCommandByName(ctx context.Context, name string) (*command.Command, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, script FROM commands WHERE name = $1`, name)

	var c command.Command
	err := row.Scan(&c.ID, &c.Name, &c.Script)
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
