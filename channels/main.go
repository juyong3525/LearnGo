package main

import (
	"fmt"
)

func main() {
	c := make(chan string)
	people := [5]string{"yong", "eun", "sam", "USguy", "Japanguy"}

	for _, person := range people {
		go isSexy(person, c)
	}

	for i := 0; i < len(people); i++ {
		fmt.Println(<-c)
	}
}

func isSexy(person string, c chan string) {
	c <- person + " is sexy"
}
