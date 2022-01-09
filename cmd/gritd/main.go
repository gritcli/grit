package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gritcli/grit/server"
)

// version string, automatically set during build process.
var version = "0.0.0"

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := server.Run(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
