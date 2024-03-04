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
2024/03/04 14:33:27 INFO stage start name=compile
2024/03/04 14:33:27 INFO executor run start cmd="{Args:[/usr/bin/g++ a.cc -o a] Env:[PATH=/usr/bin:/bin] Files:[0xc000114340 0xc000114380 0xc0001143c0] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyInCached:map[] CopyInCwd:true CopyOut:[stdout stderr] CopyOutCached:[a] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 14:33:27 INFO executor run done result="{Status:Accepted ExitStatus:0 Error: Time:355.171ms RunTime:356.069198ms Memory:53.9 MiB Files:map[stderr:len:0 stdout:len:0] FileIDs:map[a:VNQ3A3QC] FileError:[]}"
2024/03/04 14:33:27 INFO parser run start config="map[comment:compile done score:100]"
2024/03/04 14:33:27 INFO parser run done result="&{Score:100 Comment:compile done, executor status: run time: 356069198 ns, memory: 56561664 bytes}"
2024/03/04 14:33:27 INFO stage start name=run
2024/03/04 14:33:27 INFO executor run start cmd="{Args:[./a] Env:[PATH=/usr/bin:/bin] Files:[0xc000114400 0xc000114440 0xc000114480] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyInCached:map[a:a] CopyInCwd:false CopyOut:[stdout stderr] CopyOutCached:[] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 14:33:27 INFO executor run done result="{Status:Accepted ExitStatus:0 Error: Time:1.393ms RunTime:2.2294ms Memory:12.8 MiB Files:map[stderr:len:0 stdout:len:2] FileIDs:map[] FileError:[]}"
2024/03/04 14:33:27 INFO parser run start config="map[score:100 stdoutPath:1.stdout]"
2024/03/04 14:33:27 INFO parser run done result="&{Score:100 Comment:}"
2024/03/04 14:33:27 INFO stage result name=compile score=100 comment="compile done, executor status: run time: 356069198 ns, memory: 56561664 bytes"
2024/03/04 14:33:27 INFO stage result name=run score=100 comment=""
2024/03/04 14:33:27 INFO executor cleanup start name=dummy
2024/03/04 14:33:27 INFO executor cleanup done name=dummy
2024/03/04 14:33:27 INFO executor cleanup start name=sandbox
2024/03/04 14:33:27 INFO executor cleanup done name=sandbox
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
