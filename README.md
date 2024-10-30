# JOJ3

[![Go Report Card](https://goreportcard.com/badge/github.com/joint-online-judge/JOJ3)](https://goreportcard.com/report/github.com/joint-online-judge/JOJ3)

## Quick Start

1. Clone this repo in a Linux computer. For Windows, use [WSL 2](https://learn.microsoft.com/en-us/windows/wsl/install).

```bash
$ git clone ssh://git@focs.ji.sjtu.edu.cn:2222/JOJ/JOJ3.git
```

2. Install [Go](https://go.dev/doc/install). Also, make sure `make` and `git` are installed and all 3 programs are presented in `$PATH`.

    - If you have problem on connecting to the Go website and Go packages, download Go from [studygolang](https://studygolang.com/dl) and run `go env -w GOPROXY=https://goproxy.io,direct` to set the Go modules mirror proxy after installing Go.

3. Enable cgroup v2 for your OS. For WSL2, check [here](https://stackoverflow.com/a/73376219/13724598). So that you do not need root permission to run `go-judge`.

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

6. Pull submodules. It might be slow, so only run it when the test branches are out of date.

```bash
$ # make sure you are in JOJ3 directory
$ make prepare-test
```

7. Build binaries in `/cmd`.

```bash
$ make
```

8. Check the functions of `joj3` with the `make test`, which should pass all the test cases. The cases used here are in `/examples`.

For now, the following checking tools are needed for test:

1. `clang`/`clang++`
2. `clang-tidy-18`
3. `cmake`
4. `make`
5. `cpplint`

```bash
$ make test
go test -coverprofile cover.out -v ./...
...
PASS
coverage: 74.0% of statements
ok      github.com/joint-online-judge/JOJ3/cmd/joj3  2.290s  coverage: 74.0% of statements
```

### For developers

1. Install [`pre-commit`](https://pre-commit.com/), [`golangci-lint`](https://golangci-lint.run).

2. Install the pre-commit hooks. It will run some checks before you commit.

```bash
$ pre-commit install
pre-commit installed at .git/hooks/pre-commit
```

3. You only need to run steps 5, 7, and 8 in the quick start during development. If the test cases need to be updated, step 6 is also needed.

## How does it work?

These steps are executed in runner-images. We use `sudo -u tt` to elevate the permission and run `joj3`. All the secret files should be stored in the host machine with user `tt` and mounted into the runner (e.g. `/home/tt/.config`). Since the runner uses user `student`, we can keep the data safe. A single call to `joj3` executable will run 2 parts:

1. Run JOJ3 stages
    1. Parse the message.
        - It will use the git commit message from `HEAD`. The message should meet the [Conventional Commits specification](https://www.conventionalcommits.org/). We use `scope` and `description` here.
        - If `-tag` is specified, then it should equal to the scope of the message, or JOJ3 will not run.
    2. Find the configuration file.
        - We have `conf-root` and `conf-name` specified in the CLI argument. Then the full path of configuration file is `<conf-root>/<scope>/<conf-name>`.
    3. Generate stages.
        - We have an empty list of stages at the beginning.
        - We check all the stages from the configuration file. Stages with empty `group` field will always be added. And stages with `group = joj` will be added when `description` contains "joj" (case insensitive). You can set arbitrary group keywords in config file, in `groupKeywords` field, but by default the only group keyword is joj.
        - Every stage needs to have an unique `name`, which means if two stages have the same name, only the first one will be added.
    4. Run stages.
        - By default, all the stages will run sequentially.
        - Each stage contains a executor and multiple parsers. The executor (currently only sandbox) executes the command and parsers parse the output generated by the executor. The parsers in one stage will run sequentially, and all the output will be aggregated (scores being summed up and comment being concatenated).
        - The parser can return a force quit, which means all the stages after it will be skipped, but the remaining parsers in the current stage will run.
    5. Generate results.
        - Once the running of stages is done, it will generate a result file where the path is specified in the configuration file.
2. Run Joint-Teapot.
    1. Generally speaking, it reads the JOJ3 results file and output results on Gitea.
    2. With `joint-teapot joj3-all`, it will do the following things:
        1. Create/Edit an issue in the submitter's repo to show the results.
        2. Update the scoreboard file in grading repo.
        3. Update the failed table file in grading repo.

## Components

### Binaries (under `/cmd` and `/pkg`)

#### JOJ3

JOJ3 itself. Parsers and executors are compiled into the JOJ3 binary.

#### Sample

Just a sample on how to write an executable that can be called by the executor.

#### HealthCheck

The repohealth check will return a json list to for check result. The structure follows the score-comment pattern.

HealthCheck currently includes, `reposize`, `forbidden file`, `Metafile existence`, `non-ascii character` in file and message, `release tag`, and `ci files invariance` check.

The workflow is `joj3` pass cli args to healthcheck binary. See `./cmd/healthcheck/main.go` to view all flags.

### Executors (under `/internal/executors`)

#### Dummy

Do not execute any command. Just return empty `ExecutorResult` slice.

#### Sandbox

Run the commands in `go-judge` and output the `ExecutorResult` slice. Note: we communicate with `go-judge` using gRPC, which means `go-judge` can run anywhere as the gRPC connection can be established. In deployment, `go-judge` runs in the host machine of the Gitea runner.

### Parsers (under `/internal/parsers`)

#### Clang Tidy

Parser for `clang-tidy`, check `/examples/clangtidy` on how to call `clang-tidy` with proper parameters.

#### Cppcheck

Parser for `cppcheck`, check `/examples/cppcheck` on how to call `cppcheck` with proper parameters.

#### Cpplint

Parser for `cpplint`, check `/examples/cpplint` on how to call `cpplint` with proper parameters.

#### Diff

Compare the specified output of `ExecutorResult` with the content of the answer file. If they are the same, then score will be given. Just like a normal online judge system.

#### Dummy

Does not parse the output of `ExecutorResult`. It just output what is set inside the configuration file as score and comment. Currently it is used to output metadata for `joint-teapot`.

In `joint-teapot`, it will take the content before `-` of the comment of the first stage with name `metadata` as the exercise name and record in the scoreboard. (e.g. If the comment is `p2-s2-0xdeadbeef`, then the exercise name is `p2`.)

The comment in `metadata` can also be used to skip teapot commands. With `skip-teapot` in the comment, teapot will not run. And with `skip-scoreboard`, `skip-failed-table`, and `skip-result-issue`, only the corresponding step will be skipped, while the others will be executed. (e.g. If the comment is `p2-s2-0xdeadbeef-skip-scoreboard-skip-result-issue`, then only failed table step in teapot will run.)

#### Healthcheck

Parser for the `healthcheck` binary mentioned before.

#### Keyword

Match the given keyword from the specified output of `ExecutorResult`. For each match, a deduction of score is given. Can be useful if we do not have a specific parser for a code quality tool. Check `/examples/keyword`.

#### Result Status

Only check if all the status of the executor is `StatusAccepted`. Can be used to check whether the command run in the executor exit normally.

#### Sample

Parser for the `sample` binary mentioned before. Only used as a sample.

## Models (for developers only)

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
