package ratelimiter

import (
	"fmt"
	"strings"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

type limitKey struct {
	domain     string
	descriptor *rlsv3common.RateLimitDescriptor
}

func NewLimitKey(domain string, descriptor *rlsv3common.RateLimitDescriptor) limitKey {
	return limitKey{
		domain:     domain,
		descriptor: descriptor,
	}
}

func (lk *limitKey) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("ratelimit:v1:%s", lk.domain))
	for _, entry := range lk.descriptor.Entries {
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%s:%s", entry.Key, entry.Value))
	}
	return sb.String()
}
