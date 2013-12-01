package main

import (
	"fmt"
	"github.com/agaffney/crapsh/prompt"
)

func main() {
	fmt.Println("in main()")
	p := &prompt.Prompt{
		Text: "prompt text",
	}
	p.Show()
}
