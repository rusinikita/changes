// custom PR changes checker for dogfooding and example
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rusinikita/changes/commit/subject"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
)

const (
	maxFileChangesLen = 180
	format            = `(type)((context))?: (title)`
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

	commitParser, err := subject.NewParser(format, value.DefaultProperties)
	if err != nil {
		return err
	}

	t := tools{commitParser}

	return changesCheck(change, t)
}

type tools struct {
	*subject.Parser
}

func changesCheck(c git.Change, t tools) error {
	commits, err := c.Commits()
	if err != nil {
		return err
	}

	for _, commit := range commits {
		commitErr := checkCommit(commit, t)

		err = errors.Add(err, commitErr, commit.Subject())
	}

	if err != nil {
		return err
	}

	diff, err := c.FilesDiff()
	if err != nil {
		return err
	}

	return checkDiff(diff)
}

func checkCommit(commit git.Commit, t tools) error {
	commitSubject := commit.Subject()

	log.Println(commitSubject)

	if strings.HasPrefix(commitSubject, "Merge") {
		return nil
	}

	_, err := t.Parse(commitSubject)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(commit.Author.Email, "gmail.com") {
		return fmt.Errorf("commit from work email '%s'", commitSubject)
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
