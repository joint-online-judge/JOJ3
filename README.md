# JOJ3

[![Go Report Card](https://goreportcard.com/badge/github.com/joint-online-judge/JOJ3)](https://goreportcard.com/report/github.com/joint-online-judge/JOJ3)
[![Go Reference](https://pkg.go.dev/badge/github.com/joint-online-judge/JOJ3.svg)](https://pkg.go.dev/github.com/joint-online-judge/JOJ3)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-joint--online--judge%2FJOJ3-blue.svg)](https://deepwiki.com/joint-online-judge/JOJ3)

## Table of Contents

- [Quick Start](#quick-start)
- [Workflow](#workflow)
- [Models](#models)
- [Project Structure](#project-structure)
- [Further Documentation](#further-documentation)


## Quick Start

1. Clone this repo in a Linux computer. For Windows, use [WSL 2](https://learn.microsoft.com/en-us/windows/wsl/install).

```bash
$ git clone ssh://git@focs.ji.sjtu.edu.cn:2222/JOJ/JOJ3.git
```

2. Install [Go](https://go.dev/doc/install). Also, make sure `make` and `git` are installed and all 3 programs are presented in `$PATH`.

    - If you have problem on connecting to the Go website and Go packages, download Go from [studygolang](https://studygolang.com/dl) and run `go env -w GOPROXY=https://goproxy.io,direct` to set the Go modules mirror proxy after installing Go.

3. Enable cgroup v2 for your OS. For WSL2, check [here](https://stackoverflow.com/a/73376219/13724598). Also, enable linger for the user you used to run `go-judge` if you are using `systemd`, e.g. if the user is `go-judge`, run `loginctl enable-linger go-judge`. So that you do not need root permission to run `go-judge` (it can create a nesting cgroup in its user slice).

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

## Workflow

These steps are executed within [runner-images](https://focs.ji.sjtu.edu.cn/git/JOJ/runner-images), as specified in the YAML files under `.gitea/workflows` in student repositories. Our customized [`act_runner`](https://github.com/focs-gitea/act_runner) ensures that only labeled images controlled by administrators can be used. Furthermore, our private Docker registry requires authentication for pushing images, so only administrator-created images can be used to run `joj3`.

Inside the container created by `act_runner`, we use `sudo -E -u tt` to elevate permissions and run `joj3` with environment variables provided by Gitea Actions. All sensitive files should be stored on the host machine under the `tt` user's directory and mounted into the runner (e.g., `/home/tt/.config`). Allowed mount directories are also limited in `act_runner` configuration. Since the default `student` user inside the container (created from runner-images) shares the same UID as the `student` user on the host, which does not have the permission to access `tt`'s files. This helps ensure data security.

Here are the steps `joj3` will run.

1. Parse the message.
    - It will use the git commit message from `HEAD`. The message should meet the [Conventional Commits specification](https://www.conventionalcommits.org/). We use `scope` and `description` here. Also, a suffix `[group]` will be used to decide which stages will be run later.
    - If `-tag` is specified, then it should equal to the scope of the message, or JOJ3 will not run.
2. Find the configuration file.
    - We have `conf-root` and `conf-name` specified in the CLI argument. Then the full path of configuration file is `<conf-root>/<scope>/<conf-name>`.
    - If that configuration file does not exist, and `fallback-conf-name` is passed, it will try to read `<conf-root>/<fallback-conf-name>`.
3. Generate stages.
    - We have an empty list of stages at the beginning.
    - We check all the stages from the configuration file. Stages with empty `group` field will always be added. Stages with non-empty `group` field requires that value (case insensitive) appears in the commit group. e.g. with commit msg `feat(h5/e3): joj msan [joj]`, stages with the following `group` field will run: `""`, `"joj"`. Currently, it does not support multiple groups within one commit. If the group specified in the commit message is `[all]`, then all groups will run.
    - Every stage needs to have an unique `name`, which means if two stages have the same name, only the first one will be added.
4. Run stages.
    - By default, all the stages will run sequentially.
    - Each stage contains a executor and multiple parsers. The executor executes the command and parsers parse the output generated by the executor. The parsers in one stage will run sequentially, and all the output will be aggregated (scores being summed up and comment being concatenated).
    - The parser can return a force quit, which means all the stages after it will be skipped, but the remaining parsers in the current stage will run.
5. Generate results.
    - Once the running of stages is done, it will generate a result file where the path is specified in the configuration file.
6. Run optional stages.
    - If pre-stages and post-stages is specified, it will run before step 4 and after step 5, responsively. The result of these optional stages will not affect regular stages. Now we run `joint-teapot` in post-stages as it needs the output file from regulara stages to post the issue to the corresponding repo for students.

## Models

The program parses the configuration file to run multiple stages.

Each stage contains an executor and multiple parsers. An executor takes a `Cmd` and returns an `ExecutorResult`, while a parser takes an `ExecutorResult` and its configuration and returns a `ParserResult` and `bool` to indicate whether we should skip the rest stages.

### `StageExecutor.Cmd` (executor config)

Check `Cmd` at <https://github.com/criyle/go-judge#rest-api-interface>.
Some difference:

-   `CopyInDir string`: set to non-empty string to add everything in that directory to `CopyIn`.
-   `CopyInCached map[string]string`: key: file name in the sandbox, value: file name used in `CopyOutCached`.
-   `LocalFile`: now supports the relative path

### `ExecutorResult` (executor result)

Check the `Result` at <https://github.com/criyle/go-judge#rest-api-interface>.

### `StageParser.Conf` (parser config)

Check <https://pkg.go.dev/github.com/joint-online-judge/JOJ3/internal/parser>. `type Conf` in each package defines the accepted config for each kind of parser.

### `ParserResult` (parser result)

-   `Score int`: score of the executor result.
-   `Comment string`: comment on the executor result.

## Project Structure

```
+---.gitea       # Gitea meta files
+---cmd          # Executable applications
+---examples     # Examples & testcases from submodules
+---internal     # Packages for internal use only
|   +---executor # Executors
|   +---parser   # Parsers
|   \---stage    # Structure and logic for stages
+---pkg          # Packages intended for use by other applications
\---scripts      # Various helper scripts
```

## Further Documentation

<https://pkg.go.dev/github.com/joint-online-judge/JOJ3>
<https://deepwiki.com/joint-online-judge/JOJ3>
