package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// stream interceptor is used to print a log on each stream RPC.
func (b *Bootstrap) allStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	b.Logger.Info("stream rpc called", zap.String("method", info.FullMethod))

	return handler(srv, ss)
}

// allUnaryInterceptor interceptor checks the status of a service before running the gRPC function.
func (b *Bootstrap) allUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	b.Logger.Info("rpc called", zap.String("method", info.FullMethod))

	return handler(ctx, req)
}
