package server

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	system      = ""
	pingTimeout = 200 * time.Millisecond
	interval    = time.Second
)

func StartReadinessReporter(parent context.Context, hs *health.Server, rdb *redis.Client) context.CancelFunc {
	ctx, cancel := context.WithCancel(parent)
	hs.SetServingStatus(system, healthpb.HealthCheckResponse_NOT_SERVING)

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				timeout, c := context.WithTimeout(ctx, pingTimeout)
				err := rdb.Ping(timeout).Err()
				c()
				if err != nil {
					hs.SetServingStatus(system, healthpb.HealthCheckResponse_NOT_SERVING)
				} else {
					hs.SetServingStatus(system, healthpb.HealthCheckResponse_SERVING)
				}
			case <-ctx.Done():
				hs.SetServingStatus(system, healthpb.HealthCheckResponse_NOT_SERVING)
				return
			}
		}
	}()

	return cancel
}
