package main

import (
	"fmt"

	"github.com/yong/Dictionary/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	word := "hello"

	dictionary.Add(word, "Gretting")
	dictionary.Delete(word)

	err := dictionary.Delete(word)
	if err != nil {
		fmt.Println(err)
	}
}
