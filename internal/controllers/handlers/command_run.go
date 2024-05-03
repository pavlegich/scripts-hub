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
