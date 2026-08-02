package sandbox

import (
	"math"
	"testing"
	"time"

	"github.com/joint-online-judge/JOJ3/internal/stage"
)

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
