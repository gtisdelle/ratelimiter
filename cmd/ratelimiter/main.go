package main

import (
	"context"
	"flag"
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

var (
	limit      = flag.Int("limit", 10, "Requests Per Window")
	windowSize = flag.Duration("window", 20*time.Second, "Window Size")
	listenAddr = flag.String("listen", ":50051", "Port")
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
	flag.Parse()

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})

	grpcServer := grpc.NewServer()
	clock := ratelimiter.NewClock()
	store := ratelimiter.NewRedisStore(rdb)
	limiter := ratelimiter.NewRateLimiter(store, clock, *limit, *windowSize)
	ratelimitv1.RegisterRateLimitServiceServer(grpcServer, &rateLimitServer{limiter: limiter})

	log.Printf("rate limiter is listening on addr %s...", *listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
