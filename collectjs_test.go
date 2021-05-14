package collectjs

import (
	"fmt"
	"git.internal.yunify.com/MDMP2/collectjs/pkg/json/jsonparser"
)

func Example_Combine() {
	collection := New(`["name", "number"]`)
	collection2 := New(`["Mohamed Salah", 11]`)
	combine := collection.Combine(collection2)
	fmt.Println(string(combine.All()))

	// Output:
	// {"name":Mohamed Salah,"number":Mohamed Salah}
}

func Example_Get() {
	collection := New(`{
		"cpu":1,
		"mem": ["lo0", "eth1", "eth2"],
		"a":[{"v":0},{"v":1},{"v":2}],
		"b":[{"v":{"cv":1}},{"v":{"cv":1}},{"v":{"cv":1}}],
		"where": 10,
		"metadata": {"name": "Light1", "price": 11.05}
	}`)
	combine := collection.Get("a[#].v")
	fmt.Println(string(combine.All()))

	// Output:
	// [0,1,2]
}

func Example_GroupBy() {
	collection := New(`[
		  {
			"product": "Chair",
			"manufacturer": "IKEA"
		  },
		  {
			"product": "Desk",
			"manufacturer": "IKEA"
		  },
		  {
			"product": "Chair",
			"manufacturer": "Herman Miller"
		  }
		]`)
	combine := collection.GroupBy("manufacturer")
	fmt.Println(string(combine.All()))

	// Output:
	// [0,1,2]
}

func Example_Map() {
	collection := New(`[
		  {
			"product": "Chair",
			"manufacturer": "IKEA"
		  },
		  {
			"product": "Desk",
			"manufacturer": "IKEA"
		  },
		  {
			"product": "Chair",
			"manufacturer": "Herman Miller"
		  }
		]`)
	combine := collection.Map(func(bytes []byte) []byte {
		fmt.Println("++",string(bytes))
		return bytes
	})
	fmt.Println(string(combine.All()))

	// Output:
	// [0,1,2]
}

func Example_Demo() {
	collection1 := New(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"6519189504"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"1787977728"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"5775802368"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"19626115072"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3252543488"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620999810.899,"8203091968"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620999810.899,"8202657792"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620999810.899,"25112969216"]},{"metric":{"__name__":"node_memory_MemTotal_bytes","instance":"localhost:9100","job":"linux"},"value":[1620999810.899,"3972988928"]}]}}`)
	//collection2 := New(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.102:9100","job":"linux"},"value":[1620986447.013,"6506049536"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.14.146:9100","job":"linux"},"value":[1620986447.013,"1795796992"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.163:9100","job":"linux"},"value":[1620986447.013,"5764272128"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"192.168.21.174:9100","job":"linux"},"value":[1620986447.013,"19644997632"]},{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"localhost:9100","job":"linux"},"value":[1620986447.013,"3249033216"]}]}}`)

	result := collection1.Get("data.result")
	ret := result.GroupBy("metric.instance") //node_memory_MemTotal_bytes
	fmt.Println(string(ret.All()))

	// Output:
	// [0,1,2]
}

func Example_Demo2() {
	collection := New(`{
		"cpu":1,
		"mem": ["lo0", "eth1", "eth2"],
		"a":[{"v":0},{"v":1},{"v":2}],
		"b":[{"v":{"cv":1}},{"v":{"cv":1}},{"v":{"cv":1}}],
		"where": 10,
		"metadata": {"name": "Light1", "price": 11.05}
	}`)
	z, e := jsonparser.Append(collection.raw, []byte("11"), "bbb")
	z, e = jsonparser.Append(z, []byte("11"), "bbb")
	z, e = jsonparser.Append(z, []byte("11"), "bbb")
	z, e = jsonparser.Append(z, []byte("11"), "bbb")
	fmt.Println(string(z), e)

	// Output:
}
