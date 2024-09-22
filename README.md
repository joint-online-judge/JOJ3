# JOJ3

## Quick Start

1. Make sure you are in a Unix-like OS (Linux, MacOS). For Windows, use [WSL 2](https://learn.microsoft.com/en-us/windows/wsl/install).

2. Install [Go](https://go.dev/doc/install). Also, make sure `make` and `git` are installed and all 3 programs are presented in `$PATH`.

3. Enable cgroup v2 for your OS. Check [here](https://stackoverflow.com/a/73376219/13724598). So that you do not need root permission to run `go-judge`.

4. Clone [go-judge](https://github.com/criyle/go-judge).

```bash
$ git clone https://github.com/criyle/go-judge && cd go-judge
$ go build -o ./tmp/go-judge ./cmd/go-judge
```

5. Run `go-judge`.

```bash
$ # make sure you are in go-judge directory
$ ./tmp/go-judge -http-addr 0.0.0.0:5050 -grpc-addr 0.0.0.0:5051 -monitor-addr 0.0.0.0:5052 -enable-grpc -enable-debug -enable-metrics
```

6. Pull submodules. It might be slow, so only run it when necessary.

```bash
$ # make sure you are in JOJ3 directory
$ make prepare-test
```

7. Build binaries in `/cmd`.

```bash
$ make
```

8. Check the functions of `joj3` with the `make test`, which should pass all the test cases. The cases used here are in `/examples`.

Note: you may fail the test if the checking tools are not installed. e.g. For the test case `cpplint/sillycode`, you need to install `cpplint` in `/usr/bin` or `/usr/local/bin`.

```bash
$ make test
go test -coverprofile cover.out -v ./...
...
PASS
coverage: 74.0% of statements
ok      focs.ji.sjtu.edu.cn/git/FOCS-dev/JOJ3/cmd/joj3  2.290s  coverage: 74.0% of statements
```

### For developers

1. Install [`pre-commit`](https://pre-commit.com/), [`golangci-lint`](https://golangci-lint.run), [`goimports`](https://golang.org/x/tools/cmd/goimports), [`gofumpt`](https://github.com/mvdan/gofumpt).

2. Install the pre-commit hooks. It will run some checks before you commit.

```bash
$ pre-commit install
pre-commit installed at .git/hooks/pre-commit
```

3. You only need to run steps 5, 7, and 8 in the quick start during development. If the test cases need to be updated, step 6 is also needed.

## Models

The program parses the configuration file to run multiple stages.

Each stage contains an executor and parser. An executor just executes a command and returns the original result (stdout, stderr, output files). We can limit the time and memory used by each command in the executor. We run all kinds of commands in executors of different stages, including code formatting, static check, compilation, and execution. A parser takes the result and the configuration of the stage to parse the result and return the score and comment. e.g. If in the current stage, the executor runs a `clang-tidy` command, then we can use the clang-tidy parser in the configuration file to parse the stdout of the executor result and check whether some of the rules are followed. We can deduct the score and add some comments based on the result, and return the score and comment as the output of this stage. This stage ends here and the next stage starts.

In codes, an executor takes a `Cmd` and returns an `ExecutorResult`, while a parser takes an `ExecutorResult` and its conf and returns a `ParserResult` and `bool` to indicate whether we should skip the rest stages.

### `Cmd`

Check `Cmd` at <https://github.com/criyle/go-judge#rest-api-interface>.

Some difference:

-   `CopyInDir string`: set to non-empty string to add everything in that directory to `CopyIn`.
-   `CopyInCached map[string]string`: key: file name in the sandbox, value: file name used in `CopyOutCached`.
-   `LocalFile`: now supports the relative path

### `ExecutorResult`

Check the `Result` at <https://github.com/criyle/go-judge#rest-api-interface>.

### `ParserResult`

-   `Score int`: score of the stage.
-   `Comment string`: comment on the stage.

## Binaries (under `/cmd` and `/pkg`)

### Sample

Just a sample on how to write an executable that can be called by the executor.

### HealthCheck

The repohealth check will return a json list to for check result. The structure follows the score-comment pattern.

HealthCheck currently includes, `reposize`, `forbidden file`, `Metafile existence`, `non-ascii character` in file and message, `release tag`, and `ci files invariance` check.

The workflow is `joj3` pass cli args to healthcheck binary. See `./cmd/healthcheck/main.go` to view all flags.

## Executors (under `/internal/executors`)

### Dummy

Do not execute any command. Just return empty `ExecutorResult` slice.

### Sandbox

Run the commands in `go-judge` and output the `ExecutorResult` slice.

## Parsers (under `/internal/parsers`)

### Clang Tidy

Parser for `clang-tidy`, check `/examples/clangtidy` on how to call `clang-tidy` with proper parameters.

### Cppcheck

Parser for `cppcheck`, check `/examples/cppcheck` on how to call `cppcheck` with proper parameters.

### Cpplint

Parser for `cpplint`, check `/examples/cpplint` on how to call `cpplint` with proper parameters.

### Diff

Compare the specified output of `ExecutorResult` with the content of the answer file. If they are the same, then score will be given. Just like a normal online judge system.

### Dummy

Does not parse the output of `ExecutorResult`. It just output what is set inside the configuration file as score and comment. Currently it is used to output metadata for `joint-teapot`.

In `joint-teapot`, it will take the content before `-` of the comment of the first stage with name `metadata` as the exercise name and record in the scoreboard. (e.g. If the comment is `p2-s2-0xdeadbeef`, then the exercise name is `p2`.)

### Healthcheck

Parser for the `healthcheck` binary mentioned before.

### Keyword

Match the given keyword from the specified output of `ExecutorResult`. For each match, a deduction of score is given. Can be useful if we do not have a specific parser for a code quality tool. Check `/examples/keyword`.

### Result Status

Only check if all the status of the executor is `StatusAccepted`. Can be used to check whether the command run in the executor exit normally.

### Sample

Parser for the `sample` binary mentioned before. Only used as a sample.
