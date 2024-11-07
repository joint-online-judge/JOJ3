package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

var runningTest bool

type multiHandler struct {
	handlers []slog.Handler
}

var runID string

func init() {
	runID = fmt.Sprintf("%06s%04s",
		strconv.FormatInt(time.Now().Unix(), 36),
		strconv.FormatInt(int64(os.Getpid()), 36))
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func setupSlog(logPath string) error {
	handlers := []slog.Handler{}
	if logPath != "" {
		// File handler for debug logs
		debugFile, err := os.OpenFile(logPath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		debugFileHandler := slog.NewTextHandler(debugFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		handlers = append(handlers, debugFileHandler)
	}
	stderrLogLevel := slog.LevelInfo
	if runningTest {
		stderrLogLevel = slog.LevelDebug
	}
	// Stderr handler for info logs and above
	stderrHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: stderrLogLevel,
	})
	handlers = append(handlers, stderrHandler)
	// Create a multi-handler
	multiHandler := &multiHandler{handlers: handlers}
	multiHandlerWithAttrs := multiHandler.WithAttrs(
		[]slog.Attr{slog.String("runID", runID)})
	// Set the default logger
	logger := slog.New(multiHandlerWithAttrs)
	slog.SetDefault(logger)
	if logPath != "" {
		slog.Info("debug log", "path", logPath)
	}
	return nil
}
