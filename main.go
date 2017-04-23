package main

import (
	"fmt"
	"os"

	"github.com/Nify/pack/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// handle.AddAllZip("t.zip")
}
