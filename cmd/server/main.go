package main

import (
	"flag"
	"log"

	"github.com/flohansen/documenter/internal/app"
)

func main() {
	var cfg app.CliConfig
	flag.Parse()

	cli := app.NewCli(cfg)
	if err := cli.Run(app.SignalContext()); err != nil {
		log.Fatalf("cli error: %v", err)
	}
}
