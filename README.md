# JOJ3

## Quick Start

1. Make sure you are in a Unix-like OS (Linux, MacOS). For Windows, use [WSL 2](https://learn.microsoft.com/en-us/windows/wsl/install).

2. Install [Go](https://go.dev/doc/install). Also make sure `make` and `git` are installed and all 3 programs are presented in `$PATH`.

3. Enable cgroups v2 for your OS. Check [here](https://stackoverflow.com/a/73376219/13724598). So that you do not need root permission to run `go-judge`.

4. Clone [go-judge](https://github.com/criyle/go-judge).
```bash
$ git clone https://github.com/criyle/go-judge && cd go-judge
$ go build -o ./tmp/go-judge ./cmd/go-judge
```

5. Run `go-judge`.
```bash
$ # make sure you are in go-judge directory
$ ./tmp/go-judge -enable-grpc -enable-debug -enable-metrics
```

6. Check the functions of `joj3` with the `make test`, which should pass all the test cases. The cases used here are in `/examples`.
```bash
$ # make sure you are in JOJ3 directory
$ make test
go test -coverprofile cover.out -v ./...
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/cmd/dummy         coverage: 0.0% of statements
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors        [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers  [no test files]
?       focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/healthcheck   [no test files]
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors/sandbox                coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/executors/dummy          coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/diff             coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/pkg/dummy         coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/stage            coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/dummy            coverage: 0.0% of statements
        focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/internal/parsers/resultstatus             coverage: 0.0% of statements
=== RUN   TestMain
=== RUN   TestMain/compile/success
    main_test.go:101: stageResults: [{Name:compile Results:[{Score:0 Comment:}]} {Name:run Results:[{Score:100 Comment:executor status: run time: 1867950 ns, memory: 10813440 bytes} {Score:100 Comment:executor status: run time: 1948947 ns, memory: 10813440 bytes}]}]
=== RUN   TestMain/compile/error
    main_test.go:101: stageResults: [{Name:compile Results:[{Score:0 Comment:Unexpected executor status: Nonzero Exit Status.}]}]
=== RUN   TestMain/dummy/success
    main_test.go:101: stageResults: [{Name:dummy Results:[{Score:110 Comment:dummy comment + comment from toml conf}]}]
=== RUN   TestMain/dummy/error
    main_test.go:101: stageResults: [{Name:dummy Results:[{Score:0 Comment:Unexpected executor status: Nonzero Exit Status.
        Stderr: dummy negative score: -1}]}]
--- PASS: TestMain (0.39s)
    --- PASS: TestMain/compile/success (0.36s)
    --- PASS: TestMain/compile/error (0.01s)
    --- PASS: TestMain/dummy/success (0.02s)
    --- PASS: TestMain/dummy/error (0.01s)
PASS
coverage: 68.5% of statements
ok      focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/cmd/joj3  0.403s  coverage: 68.5% of statements
```

### For developers

1. Install [`pre-commit`](https://pre-commit.com/), [`golangci-lint`](https://golangci-lint.run), [`goimports`](https://golang.org/x/tools/cmd/goimports), [`gofumpt`](https://github.com/mvdan/gofumpt).

2. Install the pre-commit hooks. It will run some checks before you commit.
```bash
$ pre-commit install
pre-commit installed at .git/hooks/pre-commit
```

3. You only need to run step 5 and 6 in quick start during development.

## Models

The program parses the configuration file to run multiple stages. It can create an issue on Gitea to report the result of each stage after all stages are done.

Each stage contains an executor and parser. An executor just executes a command and returns the original result (stdout, stderr, output files). We can limit the time and memory used by each command in the executor. We run all kinds of commands in executors of different stages, including code formatting, static check, compilation, and execution. A parser takes the result and the configuration of the stage to parse the result and return the score and comment. e.g. If in the current stage, the executor runs a `clang-tidy` command, then we can use the clang-tidy parser in the configuration file to parse the stdout of the executor result and check whether some of the rules are followed. We can deduct the score and add some comments based on the result, and return the score and comment as the output of this stage. This stage ends here and the next stage starts.

In codes, an executor takes a `Cmd` and returns an `ExecutorResult`, while a parser takes an `ExecutorResult` and its conf and returns a `ParserResult` and `bool` to indicate whether we should skip the rest stages.

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
