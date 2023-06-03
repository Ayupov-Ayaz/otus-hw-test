package api

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Info("grpc request", zap.String("method", info.FullMethod))

		return nil
	}
}
