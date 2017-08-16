package stringcmap

import (
	"bytes"
	"encoding/json"
)

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
