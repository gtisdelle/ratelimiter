package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	ratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	"github.com/gtisdelle/ratelimiter/internal/ratelimiter"
	"github.com/gtisdelle/ratelimiter/internal/server"
	ratelimitv1 "github.com/gtisdelle/ratelimiter/proto/ratelimit/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpchealth "google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var quietCodes = map[codes.Code]struct{}{
	codes.Canceled:           {},
	codes.InvalidArgument:    {},
	codes.NotFound:           {},
	codes.AlreadyExists:      {},
	codes.PermissionDenied:   {},
	codes.ResourceExhausted:  {},
	codes.FailedPrecondition: {},
	codes.Aborted:            {},
	codes.OutOfRange:         {},
	codes.Unauthenticated:    {},
}

var (
	bucketSize = flag.Int("bucket", 10, "size of the token bucket")
	rate       = flag.Int("rate", 1, "rate that the token bucket refills")
	listenAddr = flag.String("listen", ":50051", "port")
	reflect    = flag.Bool("reflect", false, "enable server reflection (use for dev only)")
)

type rateLimitServer struct {
	ratelimitv1.UnimplementedRateLimitServiceServer
	limiter ratelimiter.RateLimiter
}

func (s *rateLimitServer) ShouldRateLimit(ctx context.Context, req *ratelimitv1.RateLimitRequest) (*ratelimitv1.RateLimitResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is nil")
	}
	if req.Domain == "" {
		return nil, status.Errorf(codes.InvalidArgument, "domain is required")
	}

	descriptors := []*ratelimitv3.RateLimitDescriptor{
		{Entries: []*ratelimitv3.RateLimitDescriptor_Entry{{Key: "type", Value: "legacy"}}}}
	allowed, err := s.limiter.Allow(ctx, req.Domain, descriptors)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "rate limit check failed: %v", err)
	}

	return &ratelimitv1.RateLimitResponse{
		Allowed: allowed,
	}, nil
}

func unaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m, err := handler(ctx, req)
	if err == nil {
		return m, err
	}
	code := status.Code(err)
	if _, ok := quietCodes[code]; !ok {
		log.Printf("rpc=%s code=%s err=%v", info.FullMethod, code, err)
	}
	return m, err
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryLoggingInterceptor))
	clock := ratelimiter.NewClock()
	store := ratelimiter.NewRedisStore(rdb, clock, ratelimiter.Config{BucketSize: *bucketSize, Rate: *rate})
	limiter := ratelimiter.NewRateLimiter(store)
	ratelimitv1.RegisterRateLimitServiceServer(grpcServer, &rateLimitServer{limiter: limiter})

	if *reflect {
		log.Println("reflection enabled")
		reflection.Register(grpcServer)
	}

	hs := grpchealth.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, hs)
	healthCancel := server.StartReadinessReporter(context.Background(), hs, rdb)
	defer healthCancel()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("rate limiter is listening on addr %s...", *listenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received, gracefully shutting down…")

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("server gracefully closed")

	// need to handle this case so that hung requests cannot prevent a shutdown
	case <-time.After(10 * time.Second):
		log.Println("timeout reached, forcing grpc server to stop...")
		grpcServer.Stop()
	}
}
