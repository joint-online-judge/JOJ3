# JOJ3

## Try `tiger`

```bash
$ make clean && make && sudo ./build/tiger python3 -c 'bytearray(1024 * 1024); 10 ** 10 ** 3; print("out"); import sys; print("err", file=sys.stderr)'
rm -rf ./build/*
rm -rf *.out
go build -o ./build/tiger ./cmd/tiger
2024/03/01 01:25:34 INFO process created pid=3148763
2024/03/01 01:25:34 INFO done success time=16.80708ms
ReturnCode: 0
Stdout: out

Stderr: err

TimedOut: false
TimeNs: 0
MemoryByte: 0
```
