package main

import (
	"context"
	"log"
	"net"

	ratelimitv1 "github.com/gtisdelle/ratelimiter/proto/ratelimit/v1"
	"google.golang.org/grpc"
)

type rateLimitServer struct {
	ratelimitv1.UnimplementedRateLimitServiceServer
}

func (s *rateLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimitv1.RateLimitRequest) (*ratelimitv1.RateLimitResponse, error) {
	return &ratelimitv1.RateLimitResponse{
		Allowed: true,
		Message: "Hello, world!",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ratelimitv1.RegisterRateLimitServiceServer(grpcServer, &rateLimitServer{})

	log.Println("rate limiter is listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
