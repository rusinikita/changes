// custom PR changes checker for dogfooding and example
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rusinikita/changes/git"
)

const (
	maxTitleLen       = 50
	maxFileChangesLen = 180
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}

	log.Println("ok")
}

func run() error {
	change, err := git.GetChange()
	if err != nil {
		return err
	}

	return changesCheck(change)
}

func changesCheck(c git.Change) error {
	commits, err := c.Commits()
	if err != nil {
		return err
	}

	for _, commit := range commits {
		err = checkCommit(commit)
		if err != nil {
			return err
		}
	}

	diff, err := c.FilesDiff()
	if err != nil {
		return err
	}

	return checkDiff(diff)
}

func checkCommit(commit git.Commit) error {
	title := strings.SplitN(commit.Message, "\n", 2)[0] //nolint:revive

	log.Println(title)

	if strings.HasPrefix(title, "Merge") {
		return nil
	}

	if len(title) > maxTitleLen {
		return fmt.Errorf("too long title '%s'", title)
	}

	if !strings.HasSuffix(commit.Author.Email, "gmail.com") {
		return fmt.Errorf("commit from work email '%s'", title)
	}

	return nil
}

func checkDiff(diff []git.FileChange) error {
	for _, file := range diff {
		if !strings.HasSuffix(file.Path, ".go") {
			continue
		}

		if file.Stats().Additions > maxFileChangesLen {
			return fmt.Errorf("file '%s' has to many changes", file.Path)
		}
	}

	return nil
}
