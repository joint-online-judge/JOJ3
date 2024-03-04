# JOJ3

In order to register sandbox executor, you need to run go-judge before running this program.

```bash
$ make clean && make && ./build/joj3
rm -rf ./build/*
rm -rf *.out
go build -o ./build/joj3 ./cmd/joj3
2024/03/04 01:00:33 INFO stage start name=compile
2024/03/04 01:00:33 INFO sandbox run cmd="{Args:[/usr/bin/g++ a.cc -o a] Env:[PATH=/usr/bin:/bin] Files:[0xc00007e380 0xc00007e3c0 0xc00007e400] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[a.cc:{Src:<nil> Content:0xc0000245f0 FileID:<nil> Name:<nil> Max:<nil> Symlink:<nil> StreamIn:false StreamOut:false Pipe:false}] CopyInCached:map[] CopyOut:[stdout stderr] CopyOutCached:[a] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 01:00:33 INFO sandbox run ret="results:{status:Accepted  time:327939000  runTime:328796901  memory:57540608  files:{key:\"stderr\"  value:\"\"}  files:{key:\"stdout\"  value:\"\"}  fileIDs:{key:\"a\"  value:\"YCYTGTCQ\"}}"
2024/03/04 01:00:33 INFO executor done result="{Status:Accepted ExitStatus:0 Error: Time:327.939ms RunTime:328.796901ms Memory:54.9 MiB Files:map[stderr:len:0 stdout:len:0] FileIDs:map[a:YCYTGTCQ] FileError:[]}"
2024/03/04 01:00:33 INFO parser done result="&{Score:100 Comment:compile done, executor status: run time: 328796901 ns, memory: 57540608 bytes}"
2024/03/04 01:00:33 INFO stage start name=run
2024/03/04 01:00:33 INFO sandbox run cmd="{Args:[a] Env:[PATH=/usr/bin:/bin] Files:[0xc00007e440 0xc00007e480 0xc00007e4c0] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyInCached:map[a:a] CopyOut:[stdout stderr] CopyOutCached:[] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/04 01:00:33 INFO sandbox run ret="results:{status:Accepted  time:1334000  runTime:2083023  memory:15384576  files:{key:\"stderr\"  value:\"\"}  files:{key:\"stdout\"  value:\"2\\n\"}}"
2024/03/04 01:00:33 INFO executor done result="{Status:Accepted ExitStatus:0 Error: Time:1.334ms RunTime:2.083023ms Memory:14.7 MiB Files:map[stderr:len:0 stdout:len:2] FileIDs:map[] FileError:[]}"
2024/03/04 01:00:33 INFO parser done result="&{Score:100 Comment:run done, executor status: run time: 2083023 ns, memory: 15384576 bytes}"
compile: score: 100, comment: compile done, executor status: run time: 328796901 ns, memory: 57540608 bytes
run: score: 100, comment: run done, executor status: run time: 2083023 ns, memory: 15384576 bytes
2024/03/04 01:00:33 INFO sandbox cleanup
```
