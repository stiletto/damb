package main

import (
	"fmt"
	"os"
)

func main() {
	tgt := "World"
	if len(os.Args) >= 2 {
		tgt = os.Args[1]
	}

	fmt.Printf("Hello, %s!\n", tgt)
}
