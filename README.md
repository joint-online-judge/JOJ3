# JOJ3

In order to register sandbox executor, you need to run go-judge before running this program.

```bash
$ make clean && make && ./build/joj3
rm -rf ./build/*
rm -rf *.out
go build -o ./build/joj3 ./cmd/joj3
2024/03/03 18:01:11 INFO stage start name="stage 0"
2024/03/03 18:01:11 INFO sandbox run cmd="{Args:[ls] Env:[PATH=/usr/bin:/bin] Files:[0xc0000aa340 0xc0000aa380 0xc0000aa3c0] CPULimit:10000000000 RealCPULimit:0 ClockLimit:0 MemoryLimit:104857600 StackLimit:0 ProcLimit:50 CPURateLimit:0 CPUSetLimit: CopyIn:map[] CopyOut:[stdout stderr] CopyOutCached:[] CopyOutMax:0 CopyOutDir: TTY:false StrictMemoryLimit:false DataSegmentLimit:false AddressSpaceLimit:false}"
2024/03/03 18:01:11 INFO sandbox run ret="results:{status:Accepted time:1162000 runTime:3847400 memory:14929920 files:{key:\"stderr\" value:\"\"} files:{key:\"stdout\" value:\"stderr\\nstdout\\n\"}}"
2024/03/03 18:01:11 INFO executor done result="{Status:Accepted ExitStatus:0 Error: Time:1.162ms RunTime:3.8474ms Memory:14.2 MiB Files:map[stderr:len:0 stdout:len:14] FileIDs:map[] FileError:[]}"
2024/03/03 18:01:11 INFO parser done result="&{Score:100 Comment:dummy comment for stage 0, executor status: run time: 3847400 ns, memory: 14929920 bytes}"
stage 0: score: 100, comment: dummy comment for stage 0, executor status: run time: 3847400 ns, memory: 14929920 bytes
```
