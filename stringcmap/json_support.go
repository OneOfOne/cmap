package stringcmap

import (
	"bytes"
	"encoding/json"
)

// WithJSON returns a MapWithJSON with the specific unmarshal func.
func (cm *CMap) WithJSON(fn func(raw json.RawMessage) (interface{}, error)) *MapWithJSON {
	return &MapWithJSON{
		CMap:             cm,
		UnmarshalValueFn: fn,
	}
}

// MapWithJSON is a CMap with UnmarshalJSON support.
// Usage:
// 	var mwj MapWithJSON
// 	json.Unmarshal(`{"key":"value"}`, &mwj)
type MapWithJSON struct {
	*CMap
	UnmarshalValueFn func(raw json.RawMessage) (interface{}, error)
}

// UnmarshalJSON implements json.Unmarshaler
func (mwj *MapWithJSON) UnmarshalJSON(p []byte) error {
	if mwj.CMap == nil {
		mwj.CMap = New()
	}

	if mwj.UnmarshalValueFn == nil {
		var in map[string]interface{}
		if err := json.Unmarshal(p, &in); err != nil {
			return err
		}

		for k, v := range in {
			mwj.Set(k, v)
		}

		return nil
	}

	var in map[string]json.RawMessage
	if err := json.Unmarshal(p, &in); err != nil {
		return err
	}

	for k, rj := range in {
		v, err := mwj.UnmarshalValueFn(rj)
		if err != nil {
			return err
		}

		mwj.Set(k, v)
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (cm *CMap) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("{")
	if err := cm.ForEach(func(key string, val interface{}) error {
		vj, err := json.Marshal(val)
		if err != nil {
			return err
		}
		kj, _ := json.Marshal(key)
		buf.Write(kj)
		buf.WriteByte(':')
		buf.Write(vj)
		buf.WriteByte(',')
		return nil
	}); err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if out[len(out)-1] == ',' {
		out[len(out)-1] = '}'
	} else {
		out = append(out, '}')
	}
	return out, nil
}
