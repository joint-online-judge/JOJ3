// Package stage provides functionality to run stages. Each stage contains an
// executor and multiple parsers. The executor executes the command and parsers
// parse the output generated by the executor.
package stage

import (
	"fmt"
	"log/slog"
)

func Run(stages []Stage) (
	stageResults []StageResult, forceQuitStageName string, err error,
) {
	var executorResults []ExecutorResult
	var parserResults []ParserResult
	var tmpParserResults []ParserResult
	slog.Info("stage run start")
	for _, stage := range stages {
		func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error(
						"panic recovered",
						"stageName", stage.Name,
						"panic", r,
					)
					stageResults = append(stageResults, StageResult{
						Name: stage.Name,
						Results: []ParserResult{
							{
								Score: 0,
								Comment: "JOJ3 internal error. " +
									"Please contact the administrator.\n\n",
							},
						},
						ForceQuit: true,
					})
					forceQuitStageName = stage.Name
					err = fmt.Errorf("panic in stage %s: %v", stage.Name, r)
				}
			}()
			slog.Info("stage start", "name", stage.Name)
			slog.Info(
				"executor run start",
				"stageName", stage.Name,
				"name", stage.Executor.Name,
			)
			slog.Debug(
				"executor run start",
				"stageName", stage.Name,
				"name", stage.Executor.Name,
				"cmds", stage.Executor.Cmds,
			)
			executor, ok := executorMap[stage.Executor.Name]
			if !ok {
				slog.Error(
					"executor not found",
					"stageName", stage.Name,
					"name", stage.Executor.Name,
				)
				err = fmt.Errorf("executor not found: %s", stage.Executor.Name)
				return
			}
			executorResults, err = executor.Run(stage.Executor.Cmds)
			if err != nil {
				slog.Error(
					"executor run error",
					"stageName", stage.Name,
					"name", stage.Executor.Name,
					"error", err,
				)
				return
			}
			for i, executorResult := range executorResults {
				slog.Debug(
					"executor run done",
					"stageName", stage.Name,
					"case", i,
					"name", stage.Executor.Name,
					"result", executorResult,
				)
			}
			slog.Debug(
				"executor run done",
				"stageName", stage.Name,
				"name", stage.Executor.Name,
				"summary", SummarizeExecutorResults(executorResults),
			)
			parserResults = []ParserResult{}
			stageDetail := StageDetail{
				Name:        stage.Name,
				CaseDetails: make([]CaseDetail, len(executorResults)),
			}
			parserScoresMap := map[string][]int{}
			for _, parser := range stage.Parsers {
				parserScoresMap[parser.Name] = make([]int, len(executorResults))
			}
			for _, stageParser := range stage.Parsers {
				slog.Info(
					"parser run start",
					"stageName", stage.Name,
					"name", stageParser.Name,
				)
				slog.Debug(
					"parser run start",
					"stageName", stage.Name,
					"name", stageParser.Name,
					"conf", stageParser.Conf,
				)
				parser, ok := parserMap[stageParser.Name]
				if !ok {
					slog.Error(
						"parser not found",
						"stageName", stage.Name,
						"name", stageParser.Name,
					)
					err = fmt.Errorf("parser not found: %s", stageParser.Name)
					return
				}
				var parserForceQuit bool
				tmpParserResults, parserForceQuit, err = parser.Run(
					executorResults, stageParser.Conf)
				if err != nil {
					slog.Error(
						"parser run error",
						"stageName", stage.Name,
						"name", stageParser.Name,
						"error", err,
					)
					forceQuitStageName = stage.Name
					break
				}
				for i, parserResult := range tmpParserResults {
					parserScoresMap[stageParser.Name][i] += parserResult.Score
				}
				if parserForceQuit {
					slog.Error(
						"parser force quit",
						"stageName", stage.Name,
						"name", stageParser.Name,
					)
					forceQuitStageName = stage.Name
				}
				slog.Debug(
					"parser run done",
					"stageName", stage.Name,
					"name", stageParser.Name,
					"results", tmpParserResults,
				)
				if len(parserResults) == 0 {
					parserResults = tmpParserResults
				} else {
					for i := range len(parserResults) {
						parserResults[i].Score += tmpParserResults[i].Score
						parserResults[i].Comment += tmpParserResults[i].Comment + "\n\n"
					}
				}
			}
			for i := range executorResults {
				caseDetail := CaseDetail{
					Index:          i,
					ExecutorResult: executorResults[i],
					ParserScores:   make(map[string]int),
				}
				for name, scores := range parserScoresMap {
					caseDetail.ParserScores[name] = scores[i]
				}
				stageDetail.CaseDetails[i] = caseDetail
			}
			stageResults = append(stageResults, StageResult{
				Name:      stage.Name,
				Results:   parserResults,
				ForceQuit: forceQuitStageName != "",
			})
			slog.Debug("stage done", "name", stage.Name, "stageDetail", stageDetail)
		}()
		if forceQuitStageName != "" {
			break
		}
	}
	return stageResults, forceQuitStageName, err
}

func Cleanup() {
	slog.Info("stage cleanup start")
	for name, executor := range executorMap {
		err := executor.Cleanup()
		if err != nil {
			slog.Error("executor cleanup error", "name", name, "error", err)
		}
	}
}
