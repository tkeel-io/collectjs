package collectjs

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/tkeel-io/collectjs/pkg/json/gjson"
)

var raw = Byte(`{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[{"v":0},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`)
var rawArray = Byte(`[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}]`)
var rawEmptyArray = Byte(`[]`)

func TestCollect_Get(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte
		path string
		want interface{}
	}{
		{"1", raw, "", string(raw)},
		{"2", raw, "cpu", "1"},
		{"2.1", raw, "mem[0]", "\"lo0\""},
		{"3", raw, "a", `[{"v":0},{"v":1},{"v":2}]`},
		{"4", raw, "a[0]", `{"v":0}`},
		{"5", raw, "a[1]", `{"v":1}`},
		{"6", raw, "a[2]", `{"v":2}`},
		{"7", raw, "a[#]", `3`}, // count
		{"8", raw, "a[#].v", `[0,1,2]`},
		{"9", raw, "b[0].v", `{"cv":1}`},
		{"10", raw, "b[1].v", `{"cv":2}`},
		{"11", raw, "b[2].v", `{"cv":3}`},
		{"12", raw, "b[#].v", `[{"cv":1},{"cv":2},{"cv":3}]`},
		{"13", raw, "b[1].v.cv", `2`},
		{"14", raw, "b[#].v.cv", `[1,2,3]`},
		{"14", rawArray, "[0]", `{"v":{"cv":1}}`},
	}
	for _, tt := range tests {
		cc := newCollect(tt.raw)
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.Get(tt.path); !reflect.DeepEqual(string(got.raw), tt.want) {
				t.Errorf("Get() = %v, want %v", string(got.raw), tt.want)
			}
		})
	}
}

func TestCollect_Set(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		value string
		want  interface{}
	}{
		{"2", "cpu", "2", `{"cpu":2,"mem": ["lo0", "eth1", "eth2"],"a":[{"v":0},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"3", "a", `{"v":0}`, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":{"v":0},"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"4", "a[0]", `0`, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[0,{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"5", "a[0].v", `{"v":0}`, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[{"v":{"v":0}},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := newCollect(raw)
			cc.Set(tt.path, []byte(tt.value))
			if got := cc.raw; !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Get() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCollect_Append(t *testing.T) {
	tests := []struct {
		name  string
		raw   []byte
		path  string
		value string
		want  interface{}
	}{
		{"1", raw, "cpu", "2", `Unknown value type`},
		{"2", raw, "mem", "2", `{"cpu":1,"mem": ["lo0", "eth1", "eth2",2],"a":[{"v":0},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"3", raw, "a", `{"v":11}`, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[{"v":0},{"v":1},{"v":2},{"v":11}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"4", raw, "a[0]", `0`, `Unknown value type`},
		{"5", raw, "a[0].v", `{"v":0}`, `Unknown value type`},
		{"5", rawEmptyArray, "", `{"v":0}`, `[{"v":0}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := newCollect(tt.raw)
			cc.Append(tt.path, []byte(tt.value))
			if got := cc.raw; !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Get() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCollect_Del(t *testing.T) {
	tests := []struct {
		name string
		path []string
		want interface{}
	}{
		{"2", []string{"cpu"}, `{"mem": ["lo0", "eth1", "eth2"],"a":[{"v":0},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"3", []string{"a", "b", "cpu"}, `{"mem": ["lo0", "eth1", "eth2"],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"4", []string{"a[0]"}, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
		{"5", []string{"a[0].v"}, `{"cpu":1,"mem": ["lo0", "eth1", "eth2"],"a":[{},{"v":1},{"v":2}],"b":[{"v":{"cv":1}},{"v":{"cv":2}},{"v":{"cv":3}}],"where": 10,"metadata": {"name": "Light1", "price": 11.05}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := newCollect(raw)
			cc.Del(tt.path...)
			if got := cc.raw; !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Get() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func Example_Combine() {
	collection := []byte(`["name", "number"]`)
	collection2 := []byte(`["Mohamed Salah", 11]`)
	combine, _ := Combine(collection, collection2)
	fmt.Println(string(combine))

	// Output:
	// {"name":"Mohamed Salah","number":11}
}

var rawGroup = Byte(`[{"count": "1","product": "Chair","manufacturer": "IKEA"},{"sum": "10","product": "Desk","manufacturer": "IKEA"},{"product": "Chair","manufacturer": "Herman Miller"}]`)

func Example_GroupBy() {
	ret, _ := GroupBy(rawGroup, "manufacturer") //node_memory_MemTotal_bytes
	fmt.Println(string(ret))

	// Output:
	// {"IKEA":[{"count": "1","product": "Chair","manufacturer": "IKEA"},{"sum": "10","product": "Desk","manufacturer": "IKEA"}],"Herman Miller":[{"product": "Chair","manufacturer": "Herman Miller"}]}
}

func Example_MergeBy() {
	ret, _ := MergeBy(rawGroup, "product", "manufacturer") //node_memory_MemTotal_bytes
	fmt.Println(string(ret))

	// Output:
	// {"Chair+IKEA":{"count": "1","product": "Chair","manufacturer": "IKEA"},"Desk+IKEA":{"sum": "10","product": "Desk","manufacturer": "IKEA"},"Chair+Herman Miller":{"product": "Chair","manufacturer": "Herman Miller"}}
}

func Example_KeyBy() {
	ret, _ := KeyBy(rawGroup, "manufacturer") //node_memory_MemTotal_bytes
	fmt.Println(string(ret))

	// Output:
	// {"IKEA":{"sum": "10","product": "Desk","manufacturer": "IKEA"},"Herman Miller":{"product": "Chair","manufacturer": "Herman Miller"}}
}

func Example_Merge() {
	var rawObject1 = Byte(`{"id": 1,"price": 29,}`)
	var rawObject2 = Byte(`{"price": "229","discount": false}`)
	ret, _ := Merge(rawObject1, rawObject2)
	fmt.Println(string(ret))

	// Output:
	// {"id": 1,"price": "229",,"discount":false}
}

func Example_Demo() {
	collection1 := New(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"6519189504"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"1787977728"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"5775802368"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"19626115072"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3252543488"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"8202657792"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"25112969216"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3972988928"]}]}}`)
	result := collection1.Get("data.result")
	result.Map(func(key []byte, bytes []byte) []byte {
		val := New("{}")
		val.Set("timestamp", Get(bytes, "value[0]"))
		val.Set("value", Get(bytes, "value[1]"))
		ret := New("{}")
		ret.Set(string(Get(bytes, "metric.__name__")), val.raw)
		ret.Set("instance", Get(bytes, "metric.instance"))
		return ret.raw
	})
	ret, _ := GroupBy(result.raw, "instance") //node_memory_MemTotal_bytes

	metricValue := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("[0]").raw, p2.Get("[0]").raw) > 0
	}

	newCollect(ret).SortBy(metricValue)
	fmt.Println(string(result.raw))

	// Output:
	// [{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"6519189504"},"instance":"192.168.14.102:9100"},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"1787977728"},"instance":"192.168.14.146:9100"},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"5775802368"},"instance":"192.168.21.163:9100"},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"19626115072"},"instance":"192.168.21.174:9100"},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"3252543488"},"instance":"localhost:9100"},{"node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8203091968"},"instance":"192.168.14.102:9100"},{"node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8203091968"},"instance":"192.168.14.146:9100"},{"node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8202657792"},"instance":"192.168.21.163:9100"},{"node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"25112969216"},"instance":"192.168.21.174:9100"},{"node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"3972988928"},"instance":"localhost:9100"}]
}

func Example_Demo2() {
	collection1 := New(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"1"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"3"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"2"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"19626115072"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3252543488"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"8202657792"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"25112969216"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3972988928"]}]}}`)
	result := collection1.Get("data.result")
	result.Map(func(key []byte, bytes []byte) []byte {
		val := New("{}")
		val.Set("timestamp", Get(bytes, "value[0]"))
		val.Set("value", Get(bytes, "value[1]"))
		ret := New("{}")
		ret.Set(string(Get(bytes, "metric.__name__")), val.raw)
		ret.Set("instance", Get(bytes, "metric.instance"))
		return ret.raw
	})

	ret, _ := MergeBy(result.raw, "instance") //node_memory_MemTotal_bytes

	MemAvailable := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("node_memory_MemAvailable_bytes.value").raw, p2.Get("node_memory_MemAvailable_bytes.value").raw) > 0
	}
	MemTotal := func(p1, p2 *Collect) bool {
		return bytes.Compare(p1.Get("node_memory_MemTotal_bytes.value").raw, p2.Get("node_memory_MemTotal_bytes.value").raw) < 0
	}

	sorted := newCollect(ret)
	sorted.SortBy(MemTotal)
	sorted.SortBy(MemAvailable)
	fmt.Println(string(sorted.raw))

	// Output:
	// [{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"3252543488"},"instance":"localhost:9100","node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"3972988928"}},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"3"},"instance":"192.168.14.146:9100","node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8203091968"}},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"2"},"instance":"192.168.21.163:9100","node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8202657792"}},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"19626115072"},"instance":"192.168.21.174:9100","node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"25112969216"}},{"node_memory_MemAvailable_bytes":{"timestamp":1620999810.899,"value":"1"},"instance":"192.168.14.102:9100","node_memory_MemTotal_bytes":{"timestamp":1620999810.899,"value":"8203091968"}}]
}

func Example_AAA() {

	fmt.Println(gjson.Get(string(rawArray), "0"))
	fmt.Println(gjson.Get(string(`["Mohamed Salah", 11]`), "0"))

	// Output:
	// {"v":{"cv":1}}
	//Mohamed Salah
}
