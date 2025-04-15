package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := `z"**\n**"z`
	re := regexp.MustCompile(`"\*\*\\n\*\*"`)
	res := re.ReplaceAllString(text, "X")
	fmt.Println(res)
}
