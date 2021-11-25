package collectjs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tkeel-io/collectjs/pkg/json/jsonparser"
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

func ByteNew(raw []byte) *Collect {
	return newCollect(raw)
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

// GetRaw returns raw data.
func (cc *Collect) GetRaw() []byte {
	return cc.raw
}

// GetError returns collect error.
func (cc *Collect) GetError() error {
	return cc.err
}

// GetDataType returns collect data type.
func (cc *Collect) GetDataType() string {
	return cc.datatype.String()
}

func (cc *Collect) Get(path string) *Collect {
	value := Get(cc.raw, path)
	return newCollect(value)
}

func (cc *Collect) Set(path string, value []byte) {
	cc.raw, cc.err = Set(cc.raw, path, value)
}

func (cc *Collect) Append(path string, value []byte) {
	cc.raw, cc.err = Append(cc.raw, path, value)
}

func (cc *Collect) Del(path ...string) {
	cc.raw = Del(cc.raw, path...)
}

func (cc *Collect) Copy() *Collect {
	return newCollect(cc.raw)
}

type MapHandle func(key []byte, value []byte) []byte

func (cc *Collect) Foreach(fn func(key []byte, value []byte)) {
	cc.raw = ForEach(cc.raw, cc.datatype, fn)
}

func (cc *Collect) Map(handle MapHandle) {
	ret := cc.Copy()
	cc.Foreach(func(key []byte, value []byte) {
		newValue := handle(key, value)
		ret.Set(string(key), newValue)
	})
	cc.raw, cc.datatype = ret.raw, ret.datatype
}

func (cc *Collect) GroupBy(path string) *Collect {
	value, err := GroupBy(cc.raw, path)
	cc = newCollect(value)
	cc.err = err
	return cc
}

func (cc *Collect) SortBy(fn func(p1 *Collect, p2 *Collect) bool) {
	if cc.datatype != jsonparser.Array && cc.datatype != jsonparser.Object {
		cc.err = errors.New("datatype is not array or object")
		return
	}
	carr := make([]*Collect, 0)
	cc.Foreach(func(key []byte, value []byte) {
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

func Combine(key []byte, value []byte) ([]byte, error) {
	cKey := newCollect(key)
	cValue := newCollect(value)
	if cKey.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	} else if cValue.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	}

	var (
		idx int
		err error
		ret = []byte("{}")
	)

	cKey.Foreach(func(key []byte, value []byte) {
		if ret, err = jsonparser.Set(ret, Get(cValue.raw, fmt.Sprintf("[%d]", idx)), string(value)); nil != err {
			cKey.err = err
		}
		idx++
	})
	return ret, cKey.err
}

func GroupBy(json []byte, path string) ([]byte, error) {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	}

	var err error
	ret := []byte("{}")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if len(keyValue) == 0 {
			return
		}
		keys := path2JSONPARSER(string(keyValue))
		if ret, err = jsonparser.Append(ret, value, keys...); nil != err {
			c.err = err
		}
	})
	return ret, c.err
}

func MergeBy(json []byte, paths ...string) ([]byte, error) {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	}

	var err error
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

		var newValue []byte
		k := append([]byte{byte(34)}, []byte(strings.Join(keys, "+"))...)
		k = append(k, byte(34))
		oldValue := Get(ret.raw, string(k))
		if newValue, err = Merge(oldValue, value); nil != err {
			ret.err = err
		}
		ret.Set(string(k), newValue)
	})
	return ret.raw, ret.err
}

func KeyBy(json []byte, path string) ([]byte, error) {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	}

	var err error
	ret := []byte("{}")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
			keyValue = keyValue[1 : len(keyValue)-1]
		}
		ret, err = jsonparser.Set(ret, value, string(keyValue))
	})
	return ret, err
}

func Merge(oldValue []byte, mergeValue []byte) ([]byte, error) {
	if len(oldValue) == 0 {
		return mergeValue, nil
	} else if len(mergeValue) == 0 {
		return oldValue, nil
	}

	cc := newCollect(oldValue)
	mc := newCollect(mergeValue)
	if cc.datatype != jsonparser.Object {
		return nil, errors.New("datatype is not object")
	} else if mc.datatype != jsonparser.Object {
		return nil, errors.New("datatype is not object")
	}

	mc.Foreach(func(key []byte, value []byte) {
		cc.Set(string(key), value)
	})

	return cc.raw, cc.err
}

func Sort(json []byte, path string) ([]byte, error) {
	c := newCollect(json)
	if c.datatype != jsonparser.Array {
		return nil, errors.New("datatype is not array")
	}

	var err error
	ret := []byte("[]")
	c.Foreach(func(key []byte, value []byte) {
		keyValue := Get(value, path)
		if keyValue[0] == '"' && keyValue[len(keyValue)-1] == '"' {
			keyValue = keyValue[1 : len(keyValue)-1]
		}
		if ret, err = jsonparser.Append(ret, value, string(keyValue)); nil != err {
			c.err = err
		}
	})

	return ret, c.err
}
