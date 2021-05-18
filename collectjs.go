package collectjs

import (
	"bytes"
	"errors"
	"fmt"

	"git.internal.yunify.com/MDMP2/collectjs/pkg/json/gjson"
	"git.internal.yunify.com/MDMP2/collectjs/pkg/json/jsonparser"
)

var EmptyBytes = []byte("")

type Collect struct {
	path     string
	raw      []byte
	datatype jsonparser.ValueType
	offset   int
	err      error
}

func New(raw string) *Collect {
	return newCollect(Byte(raw))
}

func newCollect(data []byte) *Collect {
	collect := &Collect{}
	value := make([]byte, len(data))
	copy(value, data)
	if _, datatype, _, err := jsonparser.Get(value); err == nil {
		collect.path = ""
		collect.raw = value
		collect.datatype = datatype
		return collect
	} else {
		collect.err = err
		return collect
	}
}

func (cc *Collect) GetRaw() []byte {
	return cc.raw
}

func (cc *Collect) Get(path string) *Collect {
	value := Get(cc.raw, path)
	return newCollect(value)
}

func (cc *Collect) Set(path string, value []byte) {
	newValue := Set(cc.raw, path, value)
	cc.raw = newValue // update collect
}

func (cc *Collect) Append(path string, value []byte) {
	newValue := Append(cc.raw, path, value)
	cc.raw = newValue // update collect
}
func (cc *Collect) Del(path string) {
	newValue := Del(cc.raw, path)
	cc.raw = newValue // update collect
}

func (cc *Collect) Copy() *Collect {
	return newCollect(cc.raw)
}

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

func Del(raw []byte, path string) []byte {
	keys := path2JSONPARSER(path)
	value := jsonparser.Delete(raw, keys...)
	return value
}

type MapHandle func(key []byte, value []byte) []byte

func (cc *Collect) Map(handle MapHandle) {
	ret := cc.Copy()
	cc.foreach(func(key []byte, value []byte) {
		newValue := handle(key, value)
		ret.Set(string(key), newValue)
	})
	cc.raw = ret.raw
	cc.datatype = ret.datatype
}

func (cc *Collect) foreach(fn func(key []byte, value []byte)) {
	if cc.datatype == jsonparser.Object {
		jsonparser.ObjectEach(cc.raw, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			fn(key, warpValue(dataType, value))
			return nil
		})
	}
	if cc.datatype == jsonparser.Array {
		idx := 0
		jsonparser.ArrayEach(cc.raw, func(value []byte, dataType jsonparser.ValueType, offset int) error {
			fn(Byte(fmt.Sprintf("[%d]", idx)), warpValue(dataType, value))
			idx++
			return nil
		})

	}

}

func (cc *Collect) GroupBy(path string) *Collect {
	value := GroupBy(cc.raw, path)
	cc = newCollect(value)
	return cc
}

func (cc *Collect) SortBy(fn func(p1 *Collect, p2 *Collect) bool) {
	if cc.datatype != jsonparser.Array && cc.datatype != jsonparser.Object {
		cc.err = errors.New("datatype is not array or object")
		return
	}
	carr := make([]*Collect, 0)
	cc.foreach(func(key []byte, value []byte) {
		carr = append(carr, newCollect(value))
	})
	By(fn).Sort(carr)

	ret := New("[]")
	for _, c := range carr {
		ret.Append("", c.raw)
	}
	cc.raw = ret.raw
	cc.datatype = ret.datatype
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
	cKey.foreach(func(key []byte, value []byte) {
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
	c.foreach(func(key []byte, value []byte) {
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

func MergeBy(json []byte, path string) []byte {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return []byte("datatype is not array")
	}

	ret := New("{}")
	c.foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if len(keyValue) == 0 {
			return
		}
		oldValue := Get(ret.raw, string(keyValue))
		newValue := Merge(oldValue, value)
		ret.Set(string(keyValue), newValue)
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
	c.foreach(func(key []byte, value []byte) {
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

	mc.foreach(func(key []byte, value []byte) {
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
	c.foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
			keyValue = keyValue[1 : len(keyValue)-1]
		}
		ret, _ = jsonparser.Append(ret, value, string(keyValue))
	})

	return ret
}
