package periodicmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

type Command struct {
	cmd       *exec.Cmd
	logWriter io.Writer
	stdout    *bytes.Buffer
	Date      string
}

func (c *Command) Run(dryRun bool) error {
	log := fmt.Sprintf("%s\n", strings.Join(c.cmd.Args, " "))
	if _, err := c.logWriter.Write([]byte(log)); err != nil {
		return fmt.Errorf("write command log: %w", err)
	}

	if dryRun {
		return nil
	}

	if err := c.cmd.Run(); err != nil {
		return fmt.Errorf("run command: %w", err)
	}
	return nil
}

func (c *Command) Output() string {
	if c == nil || c.stdout == nil {
		return ""
	}
	return c.stdout.String()
}

func resolveCreateCommand(
	ctx context.Context,
	commandTemplate string,
	targetDate string,
	stdoutWriter io.Writer,
	stderrWriter io.Writer,
) (*Command, error) {
	return resolveCommand(
		ctx,
		commandTemplate,
		targetDate,
		map[string]any{
			"date":   targetDate,
			"output": "",
		},
		stdoutWriter,
		stderrWriter,
	)
}

func resolveLinkCommand(
	ctx context.Context,
	commandTemplate string,
	previousCommand Command,
	currentCommand Command,
	nextCommand Command,
	stdoutWriter io.Writer,
	stderrWriter io.Writer,
) (*Command, error) {
	return resolveCommand(
		ctx,
		commandTemplate,
		currentCommand.Date,
		map[string]any{
			"previous": map[string]string{
				"date":   previousCommand.Date,
				"output": previousCommand.Output(),
			},
			"current": map[string]string{
				"date":   currentCommand.Date,
				"output": currentCommand.Output(),
			},
			"next": map[string]string{
				"date":   nextCommand.Date,
				"output": nextCommand.Output(),
			},
		},
		stdoutWriter,
		stderrWriter,
	)
}

func resolveCommand(
	ctx context.Context,
	commandTemplate string,
	targetDate string,
	templateParams map[string]any,
	stdoutWriter io.Writer,
	stderrWriter io.Writer,
) (*Command, error) {
	t, err := template.New("").Parse(commandTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse command template: %w", err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, templateParams); err != nil {
		return nil, fmt.Errorf("execute command template: %w", err)
	}

	words, err := shellwords.Parse(b.String())
	if err != nil {
		return nil, fmt.Errorf("parse shell words: %w", err)
	}

	if len(words) == 0 {
		return nil, errors.New("command must not be empty")
	}

	executable := words[0]
	args := words[1:]
	cmd := exec.CommandContext(ctx, executable, args...)

	var stdoutBuffer bytes.Buffer
	cmd.Stdout = io.MultiWriter(stdoutWriter, &stdoutBuffer)
	cmd.Stderr = stderrWriter

	return &Command{
		cmd:       cmd,
		logWriter: stdoutWriter,
		stdout:    &stdoutBuffer,
		Date:      targetDate,
	}, nil
}
