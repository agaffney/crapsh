package util

import (
	"encoding/json"
	"fmt"
)

func DumpJson(v interface{}, label string) {
	foo, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s%s\n", label, foo)
}

func DumpObject(v interface{}, label string) {
	fmt.Printf("%s%#v\n", label, v)
}
