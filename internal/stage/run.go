package stage

import (
	"fmt"
	"log/slog"

	"code.gitea.io/sdk/gitea"
)

func Run(stages []Stage) []StageResult {
	var stageResults []StageResult
	for _, stage := range stages {
		slog.Debug("stage start", "name", stage.Name)
		slog.Debug("executor run start", "cmds", stage.ExecutorCmds)
		executor, ok := executorMap[stage.ExecutorName]
		if !ok {
			slog.Error("executor not found", "name", stage.ExecutorName)
			break
		}
		executorResults, err := executor.Run(stage.ExecutorCmds)
		if err != nil {
			slog.Error("executor run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Debug("executor run done", "results", executorResults)
		slog.Debug("parser run start", "conf", stage.ParserConf)
		parser, ok := parserMap[stage.ParserName]
		if !ok {
			slog.Error("parser not found", "name", stage.ParserName)
			break
		}
		parserResults, end, err := parser.Run(executorResults, stage.ParserConf)
		if err != nil {
			slog.Error("parser run error", "name", stage.ExecutorName, "error", err)
			break
		}
		slog.Debug("parser run done", "results", parserResults)
		stageResults = append(stageResults, StageResult{
			Name:    stage.Name,
			Results: parserResults,
		})
		if end {
			break
		}
	}
	return stageResults
}

func Cleanup() {
	for name, executor := range executorMap {
		slog.Debug("executor cleanup start", "name", name)
		err := executor.Cleanup()
		if err != nil {
			slog.Error("executor cleanup error", "name", name, "error", err)
		}
		slog.Debug("executor cleanup done", "name", name)
	}
}

func Submit(url, token, owner, repo string, results []StageResult) error {
	c, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		return err
	}
	body := "# Stages\n"
	totalScore := 0
	for _, result := range results {
		content := fmt.Sprintf("## %s\n", result.Name)
		for _, r := range result.Results {
			content += fmt.Sprintf(
				"Score: %d\nComment: %s\n\n", r.Score, r.Comment,
			)
			totalScore += r.Score
		}
		body += content
	}
	body += fmt.Sprintf("# Total Score\n%d\n", totalScore)
	issue, resp, err := c.CreateIssue(owner, repo, gitea.CreateIssueOption{
		Title: "JOJ3 result test", Body: body,
	})
	slog.Debug("create issue", "issue", issue, "resp", resp)
	if err != nil {
		return err
	}
	return nil
}
