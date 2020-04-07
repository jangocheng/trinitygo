package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"trinitygo"

	helloworldpb "trinitygo/example/pb/helloworld"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func init() {
	const bufSize = 1024 * 1024
	lis = bufconn.Listen(bufSize)
	t := trinitygo.DefaultGRPC()
	t.InitGRPC()
	helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &Server{})
	go func() {
		if err := t.GetGRPCServer().Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestRuntimeKey(t *testing.T) {
	// md := metadata.Pairs("trace_id", uuid.New().String())
	// ctxWithAuth := metadata.NewOutgoingContext(context.Background(), md)
	ctx := context.Background()
	conn, _ := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	defer conn.Close()
	client := helloworldpb.NewGreeterClient(conn)
	_, err := client.SayHello(ctx, &helloworldpb.HelloRequest{Name: "124"})
	assert.NotEqual(t, nil, err, "err shouldnot be nil")
	if err != nil {
		status, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, status.Code(), "err should be nil")
	}

}

func TestRuntimeKeyWithAuth(t *testing.T) {
	md := metadata.Pairs("trace_id", uuid.New().String())
	ctxWithAuth := metadata.NewOutgoingContext(context.Background(), md)
	conn, _ := grpc.DialContext(ctxWithAuth, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	defer conn.Close()
	client := helloworldpb.NewGreeterClient(conn)
	_, err := client.SayHello(ctxWithAuth, &helloworldpb.HelloRequest{Name: "124"})
	assert.NotEqual(t, nil, err, "err shouldnot be nil")
	if err != nil {
		status, ok := status.FromError(err)
		fmt.Println(ok)
		assert.Equal(t, codes.InvalidArgument, status.Code(), "err should be nil")
	}

}