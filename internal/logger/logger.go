package logger

import (
	"log/slog"
	"os"

	"github.com/staidinput/slogcolor"
)

func Init(env string) {
	var handler slog.Handler
	switch env {
	case "development":
		handler = slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05",
		})
	case "production":
		handler = slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
			Level:      slog.LevelInfo,
			TimeFormat: "15:04:05",
		})
	}
	slog.SetDefault(slog.New(handler))
}
