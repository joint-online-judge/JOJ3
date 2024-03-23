package sandbox

import (
	"context"
	"log/slog"

	"github.com/criyle/go-judge/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// copied from https://github.com/criyle/go-judger-demo/blob/master/apigateway/main.go
func createExecClient(execServer, token string) (pb.ExecutorClient, error) {
	conn, err := createGRPCConnection(execServer, token)
	if err != nil {
		slog.Error("gRPC connection", "error", err)
		return nil, err
	}
	return pb.NewExecutorClient(conn), nil
}

func createGRPCConnection(addr, token string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(newTokenAuth(token)))
	}
	return grpc.Dial(addr, opts...)
}

type tokenAuth struct {
	token string
}

func newTokenAuth(token string) credentials.PerRPCCredentials {
	return &tokenAuth{token: token}
}

// Return value is mapped to request headers.
func (t *tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (
	map[string]string, error,
) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (*tokenAuth) RequireTransportSecurity() bool {
	return false
}
