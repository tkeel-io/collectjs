package collectjs

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/tkeel-io/collectjs/pkg/json/gjson"
	"github.com/tkeel-io/collectjs/pkg/json/jsonparser"
)

func Get(raw []byte, path string) []byte {
	keys := path2JSONPARSER(path)

	if value, dataType, _, err := jsonparser.Get(raw, keys...); err == nil {
		return warpValue(dataType, value)
	} else {
		path = path2GJSON(path)
		ret := gjson.GetBytes(raw, path)
		return []byte(ret.String())
	}
}

func warpValue(dataType jsonparser.ValueType, value []byte) []byte {
	switch dataType {
	case jsonparser.String:
		return bytes.Join([][]byte{
			[]byte("\""), value, []byte("\""),
		}, []byte{})
	default:
		return value
	}
	return nil
}

func Set(raw []byte, path string, value []byte) []byte {
	keys := path2JSONPARSER(path)
	if value, err := jsonparser.Set(raw, value, keys...); err == nil {
		return value
	} else {
		return []byte(err.Error())
	}
}
func Append(raw []byte, path string, value []byte) []byte {
	keys := path2JSONPARSER(path)
	if value, err := jsonparser.Append(raw, value, keys...); err == nil {
		return value
	} else {
		return []byte(err.Error())
	}
}

func Del(raw []byte, path ...string) []byte {
	for _, v := range path {
		keys := path2JSONPARSER(v)
		raw = jsonparser.Delete(raw, keys...)
	}
	return raw
}

func Combine(key []byte, value []byte) []byte {
	cKey := newCollect(key)
	cValue := newCollect(value)
	if cKey.datatype != jsonparser.Array {
		cKey.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}
	if cValue.datatype != jsonparser.Array {
		cValue.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}
	ret := []byte("{}")
	idx := 0
	cKey.Foreach(func(key []byte, value []byte) {
		//fmt.Println(ret, idx, string(cValue.Get(fmt.Sprintf("[%d]", idx))))
		ret, _ = jsonparser.Set(ret, Get(cValue.raw, fmt.Sprintf("[%d]", idx)), string(value))
		idx++
	})
	return ret
}

func GroupBy(json []byte, path string) []byte {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}

	ret := []byte("{}")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if len(keyValue) == 0 {
			return
		}
		//if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
		//	keyValue = keyValue[1 : len(keyValue)-1]
		//}
		keys := path2JSONPARSER(string(keyValue))
		ret, _ = jsonparser.Append(ret, value, keys...)
	})
	return ret
}

func MergeBy(json []byte, paths ...string) []byte {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}

	ret := New("{}")
	c.Foreach(func(key []byte, value []byte) {
		keys := make([]string, 0, len(paths))
		for _, path := range paths {
			keyValue := Get(value, path)
			if len(keyValue) == 0 {
				break
			}
			keys = append(keys, string(keyValue[1:len(keyValue)-1]))
		}

		if len(keys) == 0 {
			return
		}
		k := append([]byte{byte(34)}, []byte(strings.Join(keys, "+"))...)
		k = append(k, byte(34))
		oldValue := Get(ret.raw, string(k))
		newValue := Merge(oldValue, value)
		ret.Set(string(k), newValue)
	})
	return ret.raw
}

func KeyBy(json []byte, path string) []byte {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}

	ret := []byte("{}")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
			keyValue = keyValue[1 : len(keyValue)-1]
		}
		ret, _ = jsonparser.Set(ret, value, string(keyValue))
	})
	return ret
}

func Merge(oldValue []byte, mergeValue []byte) []byte {
	if len(oldValue) == 0 {
		return mergeValue
	}
	if len(mergeValue) == 0 {
		return oldValue
	}
	cc := newCollect(oldValue)
	if cc.datatype != jsonparser.Object {
		cc.err = errors.New("datatype is not object")
		return []byte("datatype is not object")
	}

	mc := newCollect(mergeValue)
	if mc.datatype != jsonparser.Object {
		mc.err = errors.New("datatype is not object")
		return []byte("datatype is not object")
	}

	mc.Foreach(func(key []byte, value []byte) {
		cc.Set(string(key), value)
	})

	return cc.raw
}

func Sort(json []byte, path string) []byte {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}

	ret := []byte("[]")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
			keyValue = keyValue[1 : len(keyValue)-1]
		}
		ret, _ = jsonparser.Append(ret, value, string(keyValue))
	})

	return ret
}
