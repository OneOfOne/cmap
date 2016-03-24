package cmap

import (
	"bytes"
	"encoding/json"
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
	cnt := 0
	for range ch {
		cnt++
		breakFn()
	}
	if cnt != 1 {
		t.Fatalf("expected only 1 value, got %v", cnt)
	}
}

func TestJSON(t *testing.T) {
	cm := New()
	for i := 0; i < 100; i++ {
		cm.Set(strconv.Itoa(i), i)
	}
	j, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}
	ncm, err := NewFromJSON(bytes.NewReader(j), func(v interface{}) interface{} {
		if v.(json.Number) == "99" {
			return IgnoreValue
		}
		return jsonNumberToInt(v)
	})
	if err != nil {
		t.Fatal(err)
	}
	if ln := ncm.Len(); ln != cm.Len()-1 {
		t.Fatalf("wanted %d len, got %d", cm.Len()-1, ln)
	}
	for i := 0; i < 99; i++ { // 99 because we dropped 99
		if v, ok := ncm.Get(strconv.Itoa(i)).(int); !ok || v != i {
			t.Fatalf("wanted %v, got %v", i, v)
		}
	}
}

func jsonNumberToInt(v interface{}) int {
	n, _ := v.(json.Number).Int64()
	return int(n)
}
