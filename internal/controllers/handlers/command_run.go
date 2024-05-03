package handlers

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pavlegich/scripts-hub/internal/entities"
	"github.com/pavlegich/scripts-hub/internal/infra/logger"
	"github.com/pavlegich/scripts-hub/internal/service/command"
	"go.uber.org/zap"
)

// CommandWriter contains data for writing the command.
type CommandWriter struct {
	name    string
	buf     *bytes.Buffer
	service command.Service
}

// NewCommandWriter returns new CommandWriter object.
func NewCommandWriter(ctx context.Context, name string, service command.Service) *CommandWriter {
	return &CommandWriter{
		name:    name,
		buf:     &bytes.Buffer{},
		service: service,
	}
}

// Write implements writing the data into the storage.
func (w *CommandWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)

	cmd := &entities.Command{
		Name:   w.name,
		Output: w.buf.String(),
	}

	err := w.service.Update(context.Background(), cmd)
	if err != nil {
		return -1, fmt.Errorf("Write: update command failed %w", err)
	}

	return len(d), nil
}

// RunCommand executes the command and stores the output.
func (h *CommandHandler) RunCommand(name string, script string) {
	ctx := context.Background()

	bashCmd := strings.Split(script, " ")

	cmd := exec.CommandContext(ctx, bashCmd[0], bashCmd[1:]...)
	h.procs.Store(name, cmd)

	cmdWriter := NewCommandWriter(ctx, name, h.Service)

	cmd.Stdout = cmdWriter
	cmd.Stderr = cmdWriter

	err := cmd.Start()
	if err != nil {
		logger.Log.With(zap.String("cmd_name", name)).Error("RunCommand: execute command failed",
			zap.Error(err))

		cmd.Process.Kill()
		return
	}

	err = cmd.Wait()
	if err != nil {
		logger.Log.With(zap.String("cmd_name", name)).Error("RunCommand: wait command failed",
			zap.Error(err))

		cmd.Process.Kill()
		return
	}
}
