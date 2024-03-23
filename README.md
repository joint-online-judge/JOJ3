# JOJ3

## Quick Start

To register the sandbox executor, you need to run go-judge before running this program.

```bash
$ export CONF_GITEATOKEN="<YOUR_TOKEN>" && make clean && make && ./examples/success/run.sh && ./examples/compile_error/run.sh
rm -rf ./build/*
rm -rf *.out
go build -o ./build/joj3 ./cmd/joj3
++ dirname -- ./examples/success/run.sh
+ DIRNAME=./examples/success
+ cd ./examples/success
+ ./../../build/joj3
+ cat ./joj3_result.json
[{"Name":"compile","Results":[{"Score":0,"Comment":""}]},{"Name":"run","Results":[{"Score":100,"Comment":"executor status: run time: 2811900 ns, memory: 16658432 bytes"},{"Score":100,"Comment":"executor status: run time: 2578200 ns, memory: 13094912 bytes"}]}]
+ rm -f ./joj3_result.json
+ cd -
++ dirname -- ./examples/compile_error/run.sh
+ DIRNAME=./examples/compile_error
+ cd ./examples/compile_error
+ ./../../build/joj3
+ cat ./joj3_result.json
[{"Name":"compile","Results":[{"Score":0,"Comment":"Unexpected executor status: Nonzero Exit Status."}]}]
+ rm -f ./joj3_result.json
+ cd -
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
