package handlers

import (
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
	cmd     *entities.Command
	service command.Service
}

// NewCommandWriter returns new CommandWriter object.
func NewCommandWriter(ctx context.Context, name string, service command.Service) *CommandWriter {
	return &CommandWriter{
		cmd: &entities.Command{
			Name: name,
		},
		service: service,
	}
}

// Write implements writing the data into the storage.
func (w *CommandWriter) Write(d []byte) (int, error) {
	w.cmd.Output = string(d)

	err := w.service.AppendOutput(context.Background(), w.cmd)
	if err != nil {
		return -1, fmt.Errorf("Write: append command output failed %w", err)
	}

	return len(d), nil
}

// RunCommand executes the command and stores the output.
func (h *CommandHandler) RunCommand(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c, ok := <-h.jobs:
			if !ok {
				logger.Log.Error("RunCommand: channel is closed")
				return
			}

			bashCmd := strings.Split(c.Script, " ")

			cmd := exec.CommandContext(ctx, bashCmd[0], bashCmd[1:]...)
			if cmd.Err != nil {
				logger.Log.With(zap.String("cmd_name", c.Name)).Error("RunCommand: set command failed",
					zap.Error(cmd.Err), zap.String("cmd", c.Script))

				return
			}

			h.procs.Store(c.Name, cmd)

			cmdWriter := NewCommandWriter(ctx, c.Name, h.Service)

			cmd.Stdout = cmdWriter
			cmd.Stderr = cmdWriter

			err := cmd.Start()
			if err != nil {
				logger.Log.With(zap.String("cmd_name", c.Name)).Error("RunCommand: start command failed",
					zap.Error(err), zap.String("cmd", c.Script))

				return
			}
		}
	}
}
