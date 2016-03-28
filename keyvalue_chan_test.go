package cmap

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestChan(t *testing.T) {
	if reflect.TypeOf((**KeyValue)(nil)).Elem().Kind() == reflect.Interface {
		t.Skip("interfaces aren't supported by this test.")
	}
	var (
		N = 10000
		iv reflect.Value
		typ  = reflect.TypeOf(nilValue)
		zero = reflect.Zero(typ).Interface().(*KeyValue)
		ch = newSizeKeyValueChan(100)
	)
	if testing.Short() {
		N = 1000
	}
	go ch.Send(zero, true)
	if v, ok := ch.Recv(true); !ok || !reflect.DeepEqual(v, zero) {
		t.Fatal("!ok || v != zero")
	}
	for {
		var ok bool
		if iv, ok = quick.Value(typ, rand.New(rand.NewSource(43))); !ok {
			t.Logf("!ok creating a random value")
			return
		}
		if iv.Kind() == reflect.Ptr && iv.IsNil() {
			continue
		}
		break
	}
	rv, ok := iv.Interface().(*KeyValue)
	if !ok {
		t.Fatal("wrong value type")
	}
	go func() {
		for i := 0; i < N; i++ {
			ch.Send(rv, true)
		}
		ch.Close()
	}()
	for i := 0; i < N; i++ {
		v, ok := ch.Recv(true)
		if !ok {
			t.Fatal("!ok")
		}
		if !reflect.DeepEqual(v, rv) {
			t.Fatalf("wanted %%v, got %%v", rv, v)
		}
	}
}
