package sandbox

import (
	"context"
	"testing"
)

func TestTokenAuth(t *testing.T) {
	auth := newTokenAuth("secret")
	metadata, err := auth.GetRequestMetadata(context.Background())
	if err != nil || metadata["authorization"] != "Bearer secret" {
		t.Fatalf("GetRequestMetadata() = %v, %v", metadata, err)
	}
	if auth.RequireTransportSecurity() {
		t.Fatal("RequireTransportSecurity() = true")
	}
}

func TestCreateGRPCConnection(t *testing.T) {
	conn, err := createGRPCConnection("passthrough:///unused", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}
