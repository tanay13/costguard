package main

import (
	"os"

	"github.com/tanay13/costguard/cmd/costguard/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
