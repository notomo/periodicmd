package periodicmd

import (
	"log/slog"
	"os"
	"strconv"
)

func SetupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if v, _ := strconv.ParseBool(os.Getenv("DEBUG")); v {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}
