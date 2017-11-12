package main

import (
	"os"

	"github.com/Songmu/go-httpdate"
)

func main() {
	os.Exit(httpdate.Run(os.Args[1:]))
}
