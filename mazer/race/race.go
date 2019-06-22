package main

import (
	"noops/mazer"
)

func main() {
	err := mazer.RunRace("WillAbides")
	if err != nil {
		panic(err)
	}
}
