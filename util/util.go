package util

import (
	"encoding/json"
	"fmt"
)

func DumpJson(v interface{}) {
	foo, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", foo)
}
