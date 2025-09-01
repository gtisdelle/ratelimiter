package keyfmt

import (
	"fmt"
	"testing"

	rlsv3c "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

func TestBuildKeyWithOneEntry(t *testing.T) {
	expected := fmt.Sprintf("ratelimit:%s:foo:bar:baz", keyVersion)

	result := BuildKey("foo", &rlsv3c.RateLimitDescriptor{Entries: []*rlsv3c.RateLimitDescriptor_Entry{{
		Key:   "bar",
		Value: "baz",
	}}})

	if result != expected {
		t.Fatalf("BuildKey() = \"%s\"; want \"%s\"", result, expected)
	}
}

func TestBuildKeyWithTwoEntries(t *testing.T) {
	expected := fmt.Sprintf("ratelimit:%s:foo:bar:baz:second:value", keyVersion)

	result := BuildKey("foo", &rlsv3c.RateLimitDescriptor{Entries: []*rlsv3c.RateLimitDescriptor_Entry{{
		Key:   "bar",
		Value: "baz",
	}, {
		Key:   "second",
		Value: "value",
	}}})

	if result != expected {
		t.Fatalf("BuildKey() = \"%s\"; want \"%s\"", result, expected)
	}
}
