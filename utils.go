package collectjs

import (
	"strings"
)

func Byte(raw string) []byte {
	return []byte(raw)
}

func path2JSONPARSER(path string) ([]string) {
	keys := []string{}
	if len(path) > 0 {
		if path[0] == '"' && path[len(path)-1] == '"' {
			return []string{path[1 : len(path)-1]}
		}
		path = strings.Replace(path, "[", ".[", -1)
		keys = strings.Split(path, ".")
	}
	if len(keys) >0 && keys[0] == ""{
		return keys[1:]
	}
	return keys
}

func path2GJSON(path string) string {
	path = strings.Replace(path, "[", ".", -1)
	path = strings.Replace(path, "]", "", -1)
	if len(path) >0 && path[0] == '.'{
		return path[1:]
	}
	return path
}

