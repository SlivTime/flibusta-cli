package main

import (
	"github.com/slivtime/flibusta-cli/cmd/app-cli"
	"log"
)

func main() {
	appCli := app_cli.FlibustaCLI{}
	err := appCli.Start()
	if err != nil {
		log.Fatal(err)
	}
}
