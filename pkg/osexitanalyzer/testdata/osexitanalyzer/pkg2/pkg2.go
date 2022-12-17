package main

import (
	"os"
)

func main() {
	os.Exit(2) // want "found unexpected call in main func of main pkg"
}
