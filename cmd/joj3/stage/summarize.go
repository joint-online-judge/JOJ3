package stage

import (
	"log/slog"
	"os"

	"github.com/joint-online-judge/JOJ3/cmd/joj3/conf"
	"github.com/joint-online-judge/JOJ3/internal/stage"
)

func Summarize(
	conf *conf.Conf, stageResults []stage.StageResult, stageForceQuit bool,
) {
	actor := os.Getenv("GITHUB_ACTOR")
	repository := os.Getenv("GITHUB_REPOSITORY")
	ref := os.Getenv("GITHUB_REF")
	workflow := os.Getenv("GITHUB_WORKFLOW")
	totalScore := 0
	for _, stageResult := range stageResults {
		for _, result := range stageResult.Results {
			totalScore += result.Score
		}
	}
	slog.Info(
		"stage summary",
		"name", conf.Name,
		"totalScore", totalScore,
		"forceQuit", stageForceQuit,
		"actor", actor,
		"repository", repository,
		"ref", ref,
		"workflow", workflow,
	)
}
