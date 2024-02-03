package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/conf"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
	"github.com/rusinikita/changes/script"
)

var (
	rootCmd = &cobra.Command{
		Use:   "changes",
		Short: "Changes is command line tool for Code Review automation and changelog generation",
		Long: `A fast, flexible, yet simple tool for code review automation and changelog generation, 
				built with love by Nikita Rusin and friends in Go.
                Complete documentation is available at https://rusinikita.github.io/changes/`,
		Run: func(cmd *cobra.Command, args []string) {},
	}
	checkCmd = &cobra.Command{
		Use:     "check",
		Aliases: []string{"c", "ch"},
		Short:   "Checks commit messages and files diff using rules from a config file",
		Run: func(cmd *cobra.Command, args []string) {
			err := check(config)
			if err == nil {
				return
			}

			cmd.Println(err)
			os.Exit(1) //nolint: revive
		},
	}
)

func check(config conf.Conf) (err error) {
	commitParser, commitErr := commit.GetParser(config)
	err = errors.Add(err, commitErr)

	runner, scriptsErr := script.GetScriptsRunner(config, commitParser.Properties)
	err = errors.Add(err, scriptsErr)

	if err != nil {
		return errors.Prefix(err, "config validation")
	}

	gitChange, err := git.GetChange()
	if err != nil {
		return errors.Prefix(err, "git call")
	}

	commits, commitErr := gitChange.Commits()
	err = errors.Add(err, commitErr, "commit")

	diff, diffErr := gitChange.FilesDiff()
	err = errors.Add(err, diffErr, "diff")

	if err != nil {
		return errors.Prefix(err, "git call")
	}

	enrichedCommits, commitErr := commitParser.Parse(commits)
	err = errors.Add(err, commitErr)

	err = errors.Add(err, runner.Run(enrichedCommits, change.Changes(diff)))

	return errors.Prefix(err, "changes validation")
}
