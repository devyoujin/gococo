package main

import (
	"fmt"
	"os"

	"github.com/devyoujin/gococo/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %+v", err)
		os.Exit(1)
	}
}
