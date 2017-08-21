package stringcmap

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestJSON(t *testing.T) {
	cm := New()
	for i := 0; i < 100; i++ {
		si := strconv.Itoa(i)
		cm.Set(si, si)
	}

	j, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var mwj MapWithJSON

	if err = json.Unmarshal(j, &mwj); err != nil {
		t.Fatal(err)
	}

	okeys, nkeys := cm.Keys(), mwj.Keys()

	sort.Strings(okeys)
	sort.Strings(nkeys)

	if !reflect.DeepEqual(okeys, nkeys) {
		t.Fatal("!reflect.DeepEqual(okeys, nkeys)")
	}
}

func TestJSONType(t *testing.T) {
	cm := New()
	cm.Set("1", uint64(1))

	j, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var mwj MapWithJSON
	mwj.UnmarshalValueFn = func(j json.RawMessage) (interface{}, error) {
		var u uint64
		err = json.Unmarshal(j, &u)
		return u, err
	}

	if err = json.Unmarshal(j, &mwj); err != nil {
		t.Fatal(err)
	}

	for kv := range cm.Iter(context.Background(), 0) {
		if v, ok := kv.Value.(uint64); !ok || kv.Key != "1" || v != 1 {
			t.Fatalf("bad kv: %#+v", kv)
		}
	}
}
