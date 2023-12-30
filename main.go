package main

import (
	"log/slog"
	"os"

	"github.com/notomo/periodicmd/periodicmd"
	"github.com/notomo/periodicmd/pkg/datelib"
	"github.com/urfave/cli/v2"
)

const (
	paramConfig     = "config"
	paramStartDate  = "start-date"
	paramOffsetDays = "offset-days"
	paramDryRun     = "dry-run"
)

func main() {
	app := &cli.App{
		Name:  "periodicmd",
		Usage: "generate and execute periodic commands",

		Action: func(c *cli.Context) error {
			periodicmd.SetupLogger()

			config, err := periodicmd.ReadConfig(c.String(paramConfig))
			if err != nil {
				return err
			}

			return periodicmd.Run(
				c.Context,
				config.Tasks,
				c.String(paramStartDate),
				c.Int(paramOffsetDays),
				c.Bool(paramDryRun),
				os.Stdout,
				os.Stderr,
			)
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     paramConfig,
				Usage:    "config file path",
				Required: true,
			},
			&cli.StringFlag{
				Name:  paramStartDate,
				Value: datelib.Today(),
				Usage: "yyyy-mm-dd format date (default: today)",
			},
			&cli.IntFlag{
				Name:  paramOffsetDays,
				Value: 30,
				Usage: "offset days to calculate endDate (= startDate + offsetDays)",
			},
			&cli.BoolFlag{
				Name:  paramDryRun,
				Value: false,
				Usage: "does not execute command",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
