package cmap

import "testing"

func TestAll(t *testing.T) {
	cm := New()
	if len(cm) != DefaultShardCount {
		t.Fatalf("wanted len(cm) == %v, got %v", DefaultShardCount, len(cm))
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
