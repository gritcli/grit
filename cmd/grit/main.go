package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gritcli/grit/cli"
)

// version string, automatically set during build process.
var version = "0.0.0"

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := cli.Run(version); err != nil {
		os.Exit(1)
	}
}
