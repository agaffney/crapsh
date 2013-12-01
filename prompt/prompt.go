package prompt

import (
	"fmt"
)

type Prompt struct {
	Text string
}

func init() {
	fmt.Println("in prompt.init()")
}

func (this *Prompt) Show() {
	fmt.Println(this.Text, ">")
}
