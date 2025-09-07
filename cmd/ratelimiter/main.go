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

	"github.com/gtisdelle/ratelimiter/internal/limit"
	"github.com/gtisdelle/ratelimiter/internal/server"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	bucketSize = flag.Int("bucket", 10, "size of the token bucket")
	rate       = flag.Int("rate", 1, "rate that the token bucket refills")
	listenAddr = flag.String("listen", ":50051", "port")
	reflect    = flag.Bool("reflect", false, "enable server reflection (use for dev only)")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})

	grpcServer := grpc.NewServer()
	clock := limit.NewClock()
	store := limit.NewRedisStore(rdb, clock, limit.Config{BucketSize: *bucketSize, Rate: *rate})
	limiter := limit.NewLimiter(store, *bucketSize)
	server.Register(grpcServer, limiter)

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
