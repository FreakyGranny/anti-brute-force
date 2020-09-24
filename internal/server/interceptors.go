package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// LoggingInterceptor logs incoming requests.
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Warn().Msg("unable to log grpc request")

		return h, err
	}

	log.Info().
		Str("method", info.FullMethod).
		Str("latency", time.Since(start).String()).
		Str("ip", p.Addr.String()).
		Msg("")

	return h, err
}
