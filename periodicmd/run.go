package periodicmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/notomo/periodicmd/pkg/datelib"
	"github.com/notomo/periodicmd/pkg/trilib"
)

func Run(
	ctx context.Context,
	tasks []Task,
	targetStartDate string,
	offsetDays int,
	dryRun bool,
	stdoutWriter io.Writer,
	stderrWriter io.Writer,
) error {
	targetStart, err := datelib.Parse(targetStartDate)
	if err != nil {
		return fmt.Errorf("parse start date: %w", err)
	}

	targetEnd := targetStart.AddDate(0, 0, offsetDays)

	logger := slog.Default()
	errs := []error{}
	for _, task := range tasks {
		if err := runTask(
			ctx,
			task,
			targetStart,
			targetEnd,
			dryRun,
			stdoutWriter,
			stderrWriter,
		); err != nil {
			logger.Error(err.Error())
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func runTask(
	ctx context.Context,
	task Task,
	targetStart time.Time,
	targetEnd time.Time,
	dryRun bool,
	stdoutWriter io.Writer,
	stderrWriter io.Writer,
) error {
	periodicStart, err := datelib.Parse(task.StartDate)
	if err != nil {
		return fmt.Errorf("parse task state date: %w", err)
	}

	dates := datelib.PeriodicDates(
		periodicStart,
		targetStart,
		targetEnd,
		task.Frequency.Years,
		task.Frequency.Months,
		task.Frequency.Weeks,
		task.Frequency.Days,
	)
	createCmds := []Command{}
	for _, date := range dates {
		targetDate := date.Format(time.DateOnly)
		cmd, err := resolveCommand(
			ctx,
			task.CreateCommand,
			targetDate,
			map[string]any{
				"date":   targetDate,
				"output": "",
			},
			stdoutWriter,
			stderrWriter,
		)
		if err != nil {
			return fmt.Errorf("resolve create command: %w", err)
		}
		createCmds = append(createCmds, *cmd)
	}

	for _, cmd := range createCmds {
		if err := cmd.Run(dryRun); err != nil {
			return fmt.Errorf("run create command: %w", err)
		}
	}

	triCmds := trilib.Make(createCmds)
	linkCmds := []Command{}
	for _, tri := range triCmds {
		if task.LinkCommand == "" {
			continue
		}

		cmd, err := resolveCommand(
			ctx,
			task.LinkCommand,
			tri.Current.Date,
			map[string]any{
				"previous": map[string]string{
					"date":   tri.Previous.Date,
					"output": tri.Previous.Output(),
				},
				"current": map[string]string{
					"date":   tri.Current.Date,
					"output": tri.Current.Output(),
				},
				"next": map[string]string{
					"date":   tri.Next.Date,
					"output": tri.Next.Output(),
				},
			},
			stdoutWriter,
			stderrWriter,
		)
		if err != nil {
			return fmt.Errorf("resolve link command: %w", err)
		}
		linkCmds = append(linkCmds, *cmd)
	}

	for _, cmd := range linkCmds {
		if err := cmd.Run(dryRun); err != nil {
			return fmt.Errorf("run link command: %w", err)
		}
	}

	return nil
}
