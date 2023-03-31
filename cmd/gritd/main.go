package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dogmatiq/ferrite"
	"github.com/gritcli/grit/daemon"
)

// version string, automatically set during build process.
var version = "0.0.0"

func main() {
	ferrite.Init()
	rand.Seed(time.Now().UnixNano())

	if err := daemon.Run(version); err != nil {
		fmt.Fprintln(os.Stderr, err) // TODO: make responsibility of daemon package.
		os.Exit(1)
	}
}
