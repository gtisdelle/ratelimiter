package ratelimiter

import (
	"testing"

	ratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

func TestSingleEntry(t *testing.T) {
	entries := []*ratelimitv3.RateLimitDescriptor_Entry{{
		Key:   "foo",
		Value: "bar",
	}}
	descriptor := &ratelimitv3.RateLimitDescriptor{Entries: entries}
	limitKey := NewLimitKey("envoy", descriptor)

	result := limitKey.String()

	if result != "ratelimit:v1:envoy:foo:bar" {
		t.Fatalf("String() = %v; want \"ratelimit:v1:envoy:foo:bar\"", result)
	}
}

func TestMultiEntry(t *testing.T) {
	entries := []*ratelimitv3.RateLimitDescriptor_Entry{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key:   "baz",
			Value: "foobar",
		}}
	descriptor := &ratelimitv3.RateLimitDescriptor{Entries: entries}
	limitKey := NewLimitKey("envoy", descriptor)

	result := limitKey.String()

	if result != "ratelimit:v1:envoy:foo:bar:baz:foobar" {
		t.Fatalf("String() = %v; want \"ratelimit:v1:envoy:foo:bar:baz:foobar\"", result)
	}
}
