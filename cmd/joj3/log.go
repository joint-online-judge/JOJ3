package main

import (
	"context"
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/env"
	"github.com/joint-online-judge/JOJ3/internal/conf"
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

func newSlogAttrs(csvPath string) (attrs []slog.Attr) {
	attrs = []slog.Attr{
		slog.String("runID", env.Attr.RunID),
		slog.String("confName", env.Attr.ConfName),
		slog.String("actor", env.Attr.Actor),
		slog.String("actorName", env.Attr.ActorName),
		slog.String("actorID", env.Attr.ActorID),
		slog.String("repository", env.Attr.Repository),
		slog.String("sha", env.Attr.Sha),
		slog.String("ref", env.Attr.Ref),
	}
	// if csvPath is empty, just return
	if csvPath == "" {
		return
	}
	file, err := os.Open(csvPath)
	if err != nil {
		slog.Error("open csv", "error", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("read csv", "error", err)
			return
		}
		if len(row) < 3 {
			continue
		}
		actor := row[2]
		if actor == env.Attr.Actor {
			env.Attr.ActorName = row[0]
			env.Attr.ActorID = row[1]
			return []slog.Attr{
				slog.String("runID", env.Attr.RunID),
				slog.String("confName", env.Attr.ConfName),
				slog.String("actor", env.Attr.Actor),
				slog.String("actorName", env.Attr.ActorName),
				slog.String("actorID", env.Attr.ActorID),
				slog.String("repository", env.Attr.Repository),
				slog.String("sha", env.Attr.Sha),
				slog.String("ref", env.Attr.Ref),
			}
		}
	}
	return
}

func setupSlog(conf *conf.Conf) error {
	logPath := conf.LogPath
	attrs := newSlogAttrs(conf.ActorCsvPath)
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
		handlers = append(handlers, debugTextHandler.WithAttrs(attrs))
		// Json file handler for debug logs
		debugJsonFile, err := os.OpenFile(logPath+".ndjson",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		debugJsonHandler := slog.NewJSONHandler(debugJsonFile,
			&slog.HandlerOptions{Level: slog.LevelDebug})
		handlers = append(handlers, debugJsonHandler.WithAttrs(attrs))
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
	// Set the default logger
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)
	if logPath != "" {
		slog.Info("debug log", "path", logPath)
	}
	return nil
}
