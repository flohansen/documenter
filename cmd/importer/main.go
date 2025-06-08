package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/flohansen/documenter/internal/app"
	"gopkg.in/yaml.v3"
)

type flags struct {
	ConfigPath string
}

func main() {
	var flags flags
	flag.StringVar(&flags.ConfigPath, "config", "documenter.config.yaml", "The path to the configuration file")
	flag.Parse()

	config, err := readConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("could not read config: %v", err)
	}

	cli := app.NewImporter(nil, config)
	if err := cli.Run(app.SignalContext()); err != nil {
		log.Fatalf("cli error: %v", err)
	}
}

func readConfig(name string) (app.Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return app.Config{}, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	var config app.Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return app.Config{}, fmt.Errorf("json decode error: %w", err)
	}

	return config, nil
}
