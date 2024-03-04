# JOJ3

## Quick Start

In order to register sandbox executor, you need to run go-judge before running this program.

```bash
$ make clean && make && ./_example/simple/run.sh
rm -rf ./build/*
rm -rf *.out
go build -o ./build/joj3 ./cmd/joj3
++ dirname -- ./_example/simple/run.sh
+ DIRNAME=./_example/simple
+ cd ./_example/simple
+ ./../../build/joj3
+ cat ./joj3_result.json
[{"Name":"compile","Score":100,"Comment":"compile done, executor status: run time: 239591301 ns, memory: 57176064 bytes"},{"Name":"run","Score":100,"Comment":"executor status: run time: 1839200 ns, memory: 16826368 bytes"}]
+ rm -f ./joj3_result.json
+ cd -
```

## Models

The program parse the TOML file to run multiple stages.

Each stage contains a executor and parser.

Executor takes a `Cmd` and returns a `ExecutorResult`.

Parser takes a `ExecutorResult` and its config and returns a `ParserResult`.

### `Cmd`

Check `Cmd` in <https://github.com/criyle/go-judge#rest-api-interface>.

Some difference:

-   `CopyInCwd bool`: set to `true` to add everything in the current working directory to `CopyIn`.
-   `CopyInCached map[string]string`: key: file name in sandbox, value: file name used in `CopyOutCached`.
-   `LocalFile`: now support relative path

### `ExecutorResult`

Check `Result` in <https://github.com/criyle/go-judge#rest-api-interface>.

### `ParserResult`

-   `Score int`: score of the stage.
-   `Comment string`: comment of the stage.
