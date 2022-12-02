package main

import (
	"fmt"
	"os"

	"github.com/oxodao/photomaton/cmd"
	"github.com/oxodao/photomaton/config"
	"github.com/oxodao/photomaton/services"
)

func main() {
	if err := config.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := services.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Execute()
}
