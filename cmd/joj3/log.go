package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
)

var runningTest bool

type multiHandler struct {
	handlers []slog.Handler
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

func getSlogAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("runID", env.Attr.RunID),
		slog.String("confName", env.Attr.ConfName),
		slog.String("actor", env.Attr.Actor),
		slog.String("repository", env.Attr.Repository),
	}
}

func setupSlog(logPath string) error {
	handlers := []slog.Handler{}
	if logPath != "" {
		// Text file handler for debug logs
		debugTextFile, err := os.OpenFile(logPath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		debugTextHandler := slog.NewTextHandler(debugTextFile,
			&slog.HandlerOptions{Level: slog.LevelDebug})
		handlers = append(handlers, debugTextHandler)
		// Json file handler for debug logs
		debugJsonFile, err := os.OpenFile(logPath+".ndjson",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		debugJsonHandler := slog.NewJSONHandler(debugJsonFile,
			&slog.HandlerOptions{Level: slog.LevelDebug})
		handlers = append(handlers, debugJsonHandler)
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
	if runningTest {
		stderrJSONHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: stderrLogLevel,
		})
		handlers = append(handlers, stderrJSONHandler)
	}
	// Create a multi-handler
	multiHandler := &multiHandler{handlers: handlers}
	multiHandlerWithAttrs := multiHandler.WithAttrs(getSlogAttrs())
	// Set the default logger
	logger := slog.New(multiHandlerWithAttrs)
	slog.SetDefault(logger)
	if logPath != "" {
		slog.Info("debug log", "path", logPath)
	}
	return nil
}
