package main

import (
	"os"
)

func main() {
	os.Exit(1) // want "found call os.Exit() in main pkg"
}
