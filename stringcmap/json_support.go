package stringcmap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
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

// WriteTo implements io.WriterTo, outputs the map as a json object.
func (mwj *MapWithJSON) WriteTo(w io.Writer) (n int64, err error) {
	var buf writerWithBytes
	switch w := w.(type) {
	case writerWithBytes:
		buf = w
	default:
		buf = bufio.NewWriter(w)
	}

	_ = buf.WriteByte('{')

	mwj.ForEach(func(key string, val interface{}) bool {
		var vj []byte
		if vj, err = json.Marshal(val); err != nil {
			return false
		}
		if n > 0 {
			_ = buf.WriteByte(',')
		}
		kj, _ := json.Marshal(key)
		kn, _ := buf.Write(kj)
		_ = buf.WriteByte(':')
		vn, _ := buf.Write(vj)
		n += int64(kn + vn + 1)
		return true
	})

	if err != nil {
		n = 0
		return
	}

	_ = buf.WriteByte('}')
	n += 2 // {}

	if buf, ok := buf.(flusher); ok {
		err = buf.Flush()
	}

	return
}

// UnmarshalJSON implements json.Unmarshaler
func (mwj *MapWithJSON) UnmarshalJSON(p []byte) error {
	if mwj.UnmarshalValueFn != nil {
		return mwj.unmarshalJSONTyped(p)
	}

	return mwj.unmarshalJSON(p)
}

func (mwj *MapWithJSON) unmarshalJSON(p []byte) error {
	var in map[string]interface{}
	if err := json.Unmarshal(p, &in); err != nil {
		return err
	}

	if len(in) > 0 && mwj.CMap == nil {
		mwj.CMap = New()
	}

	for k, v := range in {
		mwj.Set(k, v)
	}

	return nil
}

func (mwj *MapWithJSON) unmarshalJSONTyped(p []byte) error {
	var in map[string]json.RawMessage
	if err := json.Unmarshal(p, &in); err != nil {
		return err
	}

	if len(in) > 0 && mwj.CMap == nil {
		mwj.CMap = New()
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
	var buf bytes.Buffer
	if _, err := cm.WithJSON(nil).WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type writerWithBytes interface {
	io.Writer
	io.ByteWriter
}

type flusher interface {
	Flush() error
}
