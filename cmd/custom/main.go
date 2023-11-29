// custom PR changes checker for dogfooding and example
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

const (
	maxTitleLen       = 40
	maxFileChangesLen = 400
)

func check(commit *object.Commit, stats object.FileStats) error {
	title := strings.SplitN(commit.Message, "\n", 2)[0] //nolint:revive

	if strings.HasPrefix(title, "Merge") {
		return nil
	}

	if len(title) > maxTitleLen {
		return fmt.Errorf("too long title '%s'", title)
	}

	for _, stat := range stats {
		if !strings.HasSuffix(stat.Name, ".go") {
			continue
		}

		if stat.Addition > maxFileChangesLen {
			return fmt.Errorf("file '%s' has to many changes in commit '%s'", stat.Name, title)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}

	log.Println("ok")
}

func run() error {
	rep, _ := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})

	return repositoryChecks(rep)
}

type repository interface {
	Log(o *git.LogOptions) (object.CommitIter, error)
}

func repositoryChecks(r repository) error {
	commitIter, err := r.Log(&git.LogOptions{
		All: false,
	})
	if err != nil {
		return err
	}

	commitsHandled := 0
	err = commitIter.ForEach(func(commit *object.Commit) error {
		commitsHandled++
		if commitsHandled > 10 {
			return storer.ErrStop
		}

		stats, statsErr := commit.Stats()
		if statsErr != nil {
			return fmt.Errorf("stats %w", statsErr)
		}

		checkErr := check(commit, stats)
		if checkErr != nil {
			return checkErr
		}

		return nil
	})

	return err
}
