package sandbox

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/criyle/go-judge/pb"
	"github.com/joint-online-judge/JOJ3/internal/stage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeExecutorClient struct {
	pb.ExecutorClient
	exec       func(context.Context, *pb.Request) (*pb.Response, error)
	fileDelete func(context.Context, *pb.FileID) (*emptypb.Empty, error)
}

func (f fakeExecutorClient) Exec(
	ctx context.Context, req *pb.Request, _ ...grpc.CallOption,
) (*pb.Response, error) {
	return f.exec(ctx, req)
}

func (f fakeExecutorClient) FileDelete(
	ctx context.Context, id *pb.FileID, _ ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return f.fileDelete(ctx, id)
}

func TestExecRPCTimeout(t *testing.T) {
	tests := []struct {
		name string
		cmds []stage.Cmd
		want time.Duration
	}{
		{name: "margin only", want: rpcTimeoutMargin},
		{
			name: "largest clock limit",
			cmds: []stage.Cmd{{ClockLimit: uint64(time.Minute)}, {ClockLimit: uint64(2 * time.Minute)}},
			want: 2*time.Minute + rpcTimeoutMargin,
		},
		{
			name: "cpu limit fallback",
			cmds: []stage.Cmd{{CPULimit: uint64(time.Minute)}},
			want: 2*time.Minute + rpcTimeoutMargin,
		},
		{
			name: "clock limit takes precedence over cpu limit",
			cmds: []stage.Cmd{{ClockLimit: uint64(time.Minute), CPULimit: uint64(10 * time.Minute)}},
			want: time.Minute + rpcTimeoutMargin,
		},
		{
			name: "cpu multiplication overflow",
			cmds: []stage.Cmd{{CPULimit: math.MaxUint64}},
			want: time.Duration(math.MaxInt64),
		},
		{
			name: "duration overflow",
			cmds: []stage.Cmd{{ClockLimit: math.MaxUint64}},
			want: time.Duration(math.MaxInt64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := execRPCTimeout(tt.cmds); got != tt.want {
				t.Fatalf("execRPCTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunAppliesComputedRPCDeadline(t *testing.T) {
	var remaining time.Duration
	client := fakeExecutorClient{
		exec: func(ctx context.Context, _ *pb.Request) (*pb.Response, error) {
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("Exec context has no deadline")
			}
			remaining = time.Until(deadline)
			return &pb.Response{}, nil
		},
	}
	executor := &Sandbox{
		execClient: client,
		cachedMap:  make(map[string]string),
	}
	want := 2*time.Minute + rpcTimeoutMargin
	_, err := executor.Run(context.Background(), []stage.Cmd{{
		Args:     []string{"true"},
		CPULimit: uint64(time.Minute),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if remaining > want || remaining < want-time.Second {
		t.Fatalf("RPC deadline remaining = %v, want approximately %v", remaining, want)
	}
}

func TestCleanupJoinsDeleteErrorsAndClearsCache(t *testing.T) {
	deleteErr := errors.New("delete failed")
	deleted := 0
	client := fakeExecutorClient{
		fileDelete: func(context.Context, *pb.FileID) (*emptypb.Empty, error) {
			deleted++
			return nil, deleteErr
		},
	}
	executor := &Sandbox{
		execClient: client,
		cachedMap:  map[string]string{"one": "1", "two": "2"},
	}
	err := executor.Cleanup(context.Background())
	if !errors.Is(err, deleteErr) {
		t.Fatalf("Cleanup() error = %v, want delete error", err)
	}
	if deleted != 2 {
		t.Fatalf("FileDelete called %d times, want 2", deleted)
	}
	if len(executor.cachedMap) != 0 {
		t.Fatalf("cached files not cleared: %v", executor.cachedMap)
	}
}
