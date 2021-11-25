package collectjs

import (
	"errors"
	"fmt"

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
	newValue := Set(cc.raw, path, value)
	cc.raw = newValue // update collect
}

func (cc *Collect) Append(path string, value []byte) {
	newValue := Append(cc.raw, path, value)
	cc.raw = newValue // update collect
}

func (cc *Collect) Del(path ...string) {
	newValue := Del(cc.raw, path...)
	cc.raw = newValue // update collect
}

func (cc *Collect) Copy() *Collect {
	return newCollect(cc.raw)
}

type MapHandle func(key []byte, value []byte) []byte

func (cc *Collect) Map(handle MapHandle) {
	ret := cc.Copy()
	cc.Foreach(func(key []byte, value []byte) {
		newValue := handle(key, value)
		ret.Set(string(key), newValue)
	})
	cc.raw = ret.raw
	cc.datatype = ret.datatype
}

func (cc *Collect) Foreach(fn func(key []byte, value []byte)) {
	// dispose object.
	if cc.datatype == jsonparser.Object {
		jsonparser.ObjectEach(cc.raw, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			fn(key, warpValue(dataType, value))
			return nil
		})
	}

	// dispose array.
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
