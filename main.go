package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/gtisdelle/ratelimiter/internal/ratelimiter"
	ratelimitv1 "github.com/gtisdelle/ratelimiter/proto/ratelimit/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimitServer struct {
	ratelimitv1.UnimplementedRateLimitServiceServer
	limiter ratelimiter.RateLimiter
}

func (s *rateLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimitv1.RateLimitRequest) (*ratelimitv1.RateLimitResponse, error) {
	allowed, err := s.limiter.Allow(req.Domain)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "rate limit check failed: %v", err)
	}

	return &ratelimitv1.RateLimitResponse{
		Allowed: allowed,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	grpcServer := grpc.NewServer()
	clock := ratelimiter.NewClock()
	store := ratelimiter.NewRedisStore(rdb)
	limit := 10
	windowSize := time.Duration(20) * time.Second
	limiter := ratelimiter.NewRateLimiter(store, clock, limit, windowSize)
	ratelimitv1.RegisterRateLimitServiceServer(grpcServer, &rateLimitServer{limiter: limiter})

	log.Println("rate limiter is listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
