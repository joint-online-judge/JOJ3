package stage

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func Write(outputPath string, results []stage.StageResult) error {
	slog.Info("output result start", "path", outputPath)
	slog.Debug("output result start", "path", outputPath, "results", results)
	content, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath,
		append(content, []byte("\n")...), 0o600)
}
