package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gritcli/grit/daemon"
)

// version string, automatically set during build process.
var version = "0.0.0"

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := daemon.Run(version); err != nil {
		fmt.Fprintln(os.Stderr, err) // TODO: make responsibility of daemon package.
		os.Exit(1)
	}
}
