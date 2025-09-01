package main

import (
	_ "embed"

	"github.com/connordoman/pos/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {

	}
}
