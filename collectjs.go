package collectjs

import (
	"errors"
	"fmt"
	"git.internal.yunify.com/MDMP2/collectjs/pkg/json/gjson"
	"git.internal.yunify.com/MDMP2/collectjs/pkg/json/jsonparser"
	"strings"
)

var EMPTY = []byte("")

type Collect struct {
	raw      []byte
	datatype jsonparser.ValueType
	offset   int
	err      error
}

func New(raw string) *Collect {
	return new([]byte(raw))
}

func NewFromByte(raw []byte) *Collect {
	return new(raw)
}

func new(data []byte) *Collect {
	collect := &Collect{}
	if value, datatype, offset, err := jsonparser.Get(data); err == nil {
		collect.raw = value
		collect.datatype = datatype
		collect.offset = offset
		return collect
	} else {
		collect.err = err
		return collect
	}
}

func (c *Collect) Go(path string) *Collect {
	//c.result = gjson.GetBytes(c.raw, path)
	keys := strings.Split(path, ".")
	if value, datatype, offset, err := jsonparser.Get(c.raw, keys...); err == nil {
		c.raw = value
		c.datatype = datatype
		c.offset = offset
		return c
	} else {
		c.err = err
		return c
	}
}

func (c *Collect) Get2(path string) *Collect {
	keys := strings.Split(path, ".")
	collect := &Collect{}
	if value, datatype, offset, err := jsonparser.Get(c.raw, keys...); err == nil {
		collect.raw = value
		collect.datatype = datatype
		collect.offset = offset
		return collect
	} else {
		collect.err = err
		return collect
	}
}

func (c *Collect) Get(path string) *Collect {
	path = strings.Replace(path, "[", ".", -1)
	path = strings.Replace(path, "]", "", -1)
	ret := gjson.GetBytes(c.raw, path)
	return new([]byte(ret.Raw))
}

func (c *Collect) GetByte(path string) []byte {
	keys := strings.Split(path, ".")
	if value, _, _, err := jsonparser.Get(c.raw, keys...); err == nil {
		return value
	} else {
		c.err = err
		return []byte{}
	}
}

func (c *Collect) GroupBy(path string) *Collect {
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return c
	}

	ret := []byte("{}")
	jsonparser.ArrayEach(c.raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		key := new(value).GetByte(path)
		ret, _ = jsonparser.Append(ret, value, string(key))
	})
	return new(ret)
}

func (c *Collect) Combine(c2 *Collect) *Collect {
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return c
	}
	ret := []byte("{}")
	idx := 0
	jsonparser.ArrayEach(c.raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ret, _ = jsonparser.Set(ret, c2.GetByte(fmt.Sprintf("[%d]", idx)), string(value))
	})
	c.raw = ret
	return c
}

type MapHandle func([]byte) []byte

func (c *Collect) Map(handle MapHandle) *Collect {
	if c.datatype != jsonparser.Array {
		c.err = errors.New("datatype is not array")
		return c
	}
	ret := []byte("[]")
	jsonparser.ArrayEach(c.raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ret, _ = jsonparser.Append(ret, handle(value))
	})
	c.raw = ret
	return c
}

func (c *Collect) All() []byte {
	return c.raw
}
