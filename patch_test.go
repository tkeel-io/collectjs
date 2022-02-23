package collectjs

import (
	"testing"

	"github.com/tidwall/gjson"
	"github.com/tkeel-io/collectjs/pkg/json/jsonparser"
)

func TestGJson(t *testing.T) {
	res := gjson.Get(`{"a":[1]}`, "a.[0]")
	t.Log(res.Type.String(), res.Raw)
}

func TestSet(t *testing.T) {
	res, err := Append([]byte(`[1,2]`), "", []byte(`20202`))
	t.Log(err)
	t.Log(string(res))
}

func TestJsonparser_Get(t *testing.T) {
	value, dataT, offset, err := jsonparser.Get([]byte(`{"a":[11]`), "a", "[0]")
	t.Log(string(value))
	t.Log(dataT)
	t.Log(offset)
	t.Log(err)
}

func TestGet(t *testing.T) {
	res, _, err := Get([]byte(`{"a":{"b":12}}`), "a.b")
	t.Log(err)
	t.Log(string(res))
}
