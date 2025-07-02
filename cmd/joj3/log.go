package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
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

func newSlogAttrs(csvPath string) []slog.Attr {
	actor := env.GetActor()
	actorName := fmt.Sprintf("Name(%s)", actor)
	actorID := fmt.Sprintf("ID(%s)", actor)

	if csvPath != "" {
		file, err := os.Open(csvPath)
		if err != nil {
			slog.Error("open csv", "error", err)
		} else {
			defer file.Close()
			reader := csv.NewReader(file)
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					slog.Error("read csv", "error", err)
					break
				}
				if len(row) >= 3 && row[2] == actor {
					actorName = row[0]
					actorID = row[1]
					break
				}
			}
		}
	}

	return []slog.Attr{
		slog.String("runID", env.GetRunID()),
		slog.String("confName", env.GetConfName()),
		slog.String("actor", actor),
		slog.String("actorName", actorName),
		slog.String("actorID", actorID),
		slog.String("repository", env.GetRepository()),
		slog.String("sha", env.GetSha()),
		slog.String("ref", env.GetRef()),
	}
}

func setupSlog(conf *conf.Conf) error {
	logPath := conf.LogPath
	attrs := newSlogAttrs(conf.ActorCsvPath)
	handlers := []slog.Handler{}
	if logPath != "" {
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0o750); err != nil {
			return err
		}
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
		debugJSONFile, err := os.OpenFile(logPath+".ndjson",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		if err != nil {
			return err
		}
		debugJSONHandler := slog.NewJSONHandler(debugJSONFile,
			&slog.HandlerOptions{Level: slog.LevelDebug})
		handlers = append(handlers, debugJSONHandler.WithAttrs(attrs))
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
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"setup slog attrs",
		attrs...,
	)
	return nil
}
