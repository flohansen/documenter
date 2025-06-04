package app

import "context"

type CliConfig struct {
}

type Cli struct {
}

func NewCli(cfg CliConfig) *Cli {
	return &Cli{}
}

func (c *Cli) Run(ctx context.Context) error {
	return nil
}
