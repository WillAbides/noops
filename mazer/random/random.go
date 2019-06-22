package main

import (
	"fmt"

	"noops/mazer"
)

func main() {
	result, err := mazer.DoRandomMaze(200, 200)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
