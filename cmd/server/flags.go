package main

import (
	goflag "flag"

	"gameserver/internal/run"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

func ParseFlags(cfg *run.Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	address := flag.StringP("address", "a", "localhost:8080", "Address and port to run server")
	database := flag.StringP("database", "d", "", "Database connection string")

	flag.Parse()

	cfg.Address = *address
	cfg.DataBaseDSN = *database
}

func SetEnv(cfg *run.Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
