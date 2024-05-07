// Package errors contains variables with error explainations.
package errors

import "errors"

var (
	ErrCmdNotFound      = errors.New("command not found")
	ErrCmdAlreadyExists = errors.New("command already exists")
)
