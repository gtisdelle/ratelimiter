package server

import (
	"context"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	rlsv3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type limiter interface {
	Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (*rlsv3.RateLimitResponse, error)
}

type rateLimitServer struct {
	rlsv3.UnimplementedRateLimitServiceServer
	limiter limiter
}

func (s *rateLimitServer) ShouldRateLimit(ctx context.Context, req *rlsv3.RateLimitRequest) (*rlsv3.RateLimitResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is nil")
	}
	if req.Domain == "" {
		return nil, status.Errorf(codes.InvalidArgument, "domain is required")
	}

	res, err := s.limiter.Allow(ctx, req.Domain, uint64(req.HitsAddend), req.Descriptors)
	if err != nil {
		// TODO
		return nil, nil
	}

	return res, nil
}

func Register(s *grpc.Server, l limiter) {
	rlsv3.RegisterRateLimitServiceServer(s, &rateLimitServer{limiter: l})
}
