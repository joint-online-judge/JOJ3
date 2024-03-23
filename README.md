# JOJ3

## Quick Start

To register the sandbox executor, you need to run go-judge before running this program.

```bash
$ make test
go test -v ./...
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/cmd/dummy [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors        [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors/dummy  [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/dummy    [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers  [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/diff     [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors/sandbox        [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/resultstatus     [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage    [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/dummy [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/healthcheck   [no test files]
=== RUN   TestMain
=== RUN   TestMain/success
    main_test.go:96: stageResults: [{Name:compile Results:[{Score:0 Comment:}]} {Name:run Results:[{Score:100 Comment:executor status: run time: 1910200 ns, memory: 13529088 bytes} {Score:100 Comment:executor status: run time: 1703000 ns, memory: 15536128 bytes}]}]
=== RUN   TestMain/compile_error
    main_test.go:96: stageResults: [{Name:compile Results:[{Score:0 Comment:Unexpected executor status: Nonzero Exit Status.}]}]
=== RUN   TestMain/dummy
    main_test.go:96: stageResults: [{Name:dummy Results:[{Score:110 Comment:dummy comment + comment from toml conf}]}]
=== RUN   TestMain/dummy_error
    main_test.go:96: stageResults: [{Name:dummy Results:[{Score:0 Comment:Unexpected executor status: Nonzero Exit Status.
        Stderr: dummy negative score: -1}]}]
--- PASS: TestMain (0.29s)
    --- PASS: TestMain/success (0.27s)
    --- PASS: TestMain/compile_error (0.01s)
    --- PASS: TestMain/dummy (0.01s)
    --- PASS: TestMain/dummy_error (0.01s)
PASS
ok      focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/cmd/joj3  0.295s
```

### For developers

Install [`pre-commit`](https://pre-commit.com/), [`golangci-lint`](https://golangci-lint.run), [`goimports`](https://golang.org/x/tools/cmd/goimports), [`gofumpt`](https://github.com/mvdan/gofumpt).

Then install the pre-commit hooks. It will run some checks before you commit.

```bash
$ pre-commit install
pre-commit installed at .git/hooks/pre-commit
```

## Models
The program parses the TOML file to run multiple stages.

Each stage contains an executor and parser.

Executor takes a `Cmd` and returns an `ExecutorResult`.

Parser takes an `ExecutorResult` and its conf and returns a `ParserResult` and `bool` to indicate whether we should skip the rest stages.

### `Cmd`

Check `Cmd` at <https://github.com/criyle/go-judge#rest-api-interface>.

Some difference:

-   `CopyInCwd bool`: set to `true` to add everything in the current working directory to `CopyIn`.
-   `CopyInCached map[string]string`: key: file name in the sandbox, value: file name used in `CopyOutCached`.
-   `LocalFile`: now supports the relative path

### `ExecutorResult`

Check the `Result` at <https://github.com/criyle/go-judge#rest-api-interface>.

### `ParserResult`

-   `Score int`: score of the stage.
-   `Comment string`: comment on the stage.
