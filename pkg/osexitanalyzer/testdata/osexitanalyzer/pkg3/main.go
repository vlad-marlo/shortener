package main

import (
	"os"
)

func main() {
	myExit(1)
}

func myExit(x int) {
	os.Exit(x)
}
