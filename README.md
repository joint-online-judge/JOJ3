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
2024/03/04 02:57:50 INFO stage start name=compile
2024/03/04 02:57:50 INFO sandbox run cmd="{Args:[/usr/bin/g++ a.cc -o a] Env:[PATH=/usr/bin:/bin] Files:[0xc00007f540 0xc00007f580 0xc00007f5c0] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyInCached:map[] CopyInCwd:true CopyOut:[stdout stderr] CopyOutCached:[a] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 02:57:50 INFO sandbox run ret="results:{status:Accepted  time:298002000  runTime:298694146  memory:57880576  files:{key:\"stderr\"  value:\"\"}  files:{key:\"stdout\"  value:\"\"}  fileIDs:{key:\"a\"  value:\"DMTRJR3V\"}}"
2024/03/04 02:57:50 INFO executor done result="{Status:Accepted ExitStatus:0 Error: Time:298.002ms RunTime:298.694146ms Memory:55.2 MiB Files:map[stderr:len:0 stdout:len:0] FileIDs:map[a:DMTRJR3V] FileError:[]}"
2024/03/04 02:57:50 INFO parser done result="&{Score:100 Comment:compile done, executor status: run time: 298694146 ns, memory: 57880576 bytes}"
2024/03/04 02:57:50 INFO stage start name=run
2024/03/04 02:57:50 INFO sandbox run cmd="{Args:[./a] Env:[PATH=/usr/bin:/bin] Files:[0xc00007f600 0xc00007f640 0xc00007f680] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyInCached:map[a:a] CopyInCwd:false CopyOut:[stdout stderr] CopyOutCached:[] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 02:57:50 INFO sandbox run ret="results:{status:Accepted  time:1122000  runTime:1723994  memory:15384576  files:{key:\"stderr\"  value:\"\"}  files:{key:\"stdout\"  value:\"2\\n\"}}"
2024/03/04 02:57:50 INFO executor done result="{Status:Accepted ExitStatus:0 Error: Time:1.122ms RunTime:1.723994ms Memory:14.7 MiB Files:map[stderr:len:0 stdout:len:2] FileIDs:map[] FileError:[]}"
2024/03/04 02:57:50 INFO parser done result="&{Score:100 Comment:}"
2024/03/04 02:57:50 INFO stage result name=compile score=100 comment="compile done, executor status: run time: 298694146 ns, memory: 57880576 bytes"
2024/03/04 02:57:50 INFO stage result name=run score=100 comment=""
2024/03/04 02:57:50 INFO sandbox cleanup
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
