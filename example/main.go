package main

import (
	"flag"
	"fmt"
	"os"
)
import "luddite"

type Config struct {
	Service luddite.ServiceConfig
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-c config.yaml]\n", os.Args[0])
}

func main() {
	var cfgFile string

	fs := flag.NewFlagSet("example", flag.ExitOnError)
	fs.StringVar(&cfgFile, "c", "config.yaml", "Path to config file")
	fs.Usage = usage
	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	cfg := Config{}
	if err := luddite.ReadConfig(cfgFile, &cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := luddite.NewService(&cfg.Service)
	s.AddCollectionResource("/users", newUserResource())
	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
