package limit

import (
	"fmt"
	"strings"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

var (
	keyVersion = "v1"
)

func BuildKey(domain string, descriptor *rlsv3common.RateLimitDescriptor) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ratelimit:%s:%s", keyVersion, domain))
	for _, entry := range descriptor.Entries {
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%s:%s", entry.Key, entry.Value))
	}
	return sb.String()
}
