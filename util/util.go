package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func DumpJson(v interface{}, label string) {
	foo, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("failed to dump object to JSON: %s", err.Error())
	}
	fmt.Printf("%s%s\n", label, foo)
}

func DumpObject(v interface{}, label string) {
	fmt.Printf("%s%#v\n", label, v)
}
