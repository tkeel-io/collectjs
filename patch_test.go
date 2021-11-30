package collectjs

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestGJson(t *testing.T) {
	res := gjson.Get(`{}`, "xxx")
	t.Log(res.Type.String())
}
