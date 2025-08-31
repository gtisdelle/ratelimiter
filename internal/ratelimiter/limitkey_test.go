package ratelimiter

import (
	"testing"

	ratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSingleEntry(t *testing.T) {
	entries := []*ratelimitv3.RateLimitDescriptor_Entry{{
		Key:   "foo",
		Value: "bar",
	}}
	descriptor := &ratelimitv3.RateLimitDescriptor{Entries: entries}
	limitKey := NewLimitKey("envoy", 1, descriptor)

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
	limitKey := NewLimitKey("envoy", 1, descriptor)

	result := limitKey.String()

	if result != "ratelimit:v1:envoy:foo:bar:baz:foobar" {
		t.Fatalf("String() = %v; want \"ratelimit:v1:envoy:foo:bar:baz:foobar\"", result)
	}
}

func TestHitsAddendSet(t *testing.T) {
	entries := []*ratelimitv3.RateLimitDescriptor_Entry{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key:   "baz",
			Value: "foobar",
		}}
	descriptor := &ratelimitv3.RateLimitDescriptor{Entries: entries, HitsAddend: &wrapperspb.UInt64Value{Value: 2}}
	limitKey := NewLimitKey("envoy", 1, descriptor)

	result := limitKey.Hits()

	if result != 2 {
		t.Fatalf("Hits() = %v; want 2", result)
	}
}

func TestHitsAddendNotSet(t *testing.T) {
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
	limitKey := NewLimitKey("envoy", 1, descriptor)

	result := limitKey.Hits()

	if result != 1 {
		t.Fatalf("Hits() = %v; want 2", result)
	}
}
