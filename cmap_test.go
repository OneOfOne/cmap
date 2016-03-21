package cmap

import (
	"strconv"
	"testing"
)

func TestSetHasGet(t *testing.T) {
	cm := New()
	if len(cm.shards) != DefaultShardCount {
		t.Fatalf("wanted len(cm) == %v, got %v", DefaultShardCount, len(cm.shards))
	}

	cm.Set("key", "value")

	if cm.Len() != 1 {
		t.Fatalf("wanted cm.Len() == 1, got %v", cm.Len())
	}

	if cm.Has("hi") {
		t.Fatal("found a key that shouldn't have been found")
	}
	if v, ok := cm.Get("key").(string); !ok || v != "value" {
		t.Fatalf("wanted `value`, got %v", v)
	}
}

func TestIter(t *testing.T) {
	cm := New()
	for i := uint64(0); i < 100; i++ {
		k := strconv.FormatUint(^i, 10)
		cm.Set(k, k)
	}
	ch, breakFn := cm.IterBuffered(20)
	for kv := range ch {
		t.Logf("%+v", kv)
		breakFn()
	}
	_ = breakFn
}
