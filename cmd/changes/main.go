package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/rusinikita/changes/conf"
)

var (
	cfgFile string
	config  conf.Conf
)

func main() {
	initConfig := func() {
		var err error

		config, err = conf.New(cfgFile)
		if err != nil {
			log.Println("File config not applied:", err)
		}
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config", "",
		"path to config file (default is .changes.[yaml,toml,json])",
	)
	rootCmd.AddCommand(checkCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
