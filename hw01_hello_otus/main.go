package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func printReversed(input string) {
	reversed := stringutil.Reverse(input)
	fmt.Println(reversed)
}

func main() {
	printReversed("Hello, OTUS!")
}
