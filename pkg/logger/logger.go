package logger

import (
	"io"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/hoangquan/retail-store-api/internal/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	writers := []io.Writer{os.Stdout}

	if cfg.LogstashAddr != "" {
		conn, err := net.Dial("tcp", cfg.LogstashAddr)
		if err != nil {
			slog.Warn("Failed to connect to Logstash, logging to stdout only", "error", err)
		} else {
			writers = append(writers, conn)
		}
	}

	writer := io.MultiWriter(writers...)

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
