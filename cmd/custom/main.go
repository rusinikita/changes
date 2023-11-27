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

func main() {
	repository, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		log.Fatalln("open", err)
	}

	commitIter, err := repository.Log(&git.LogOptions{
		All: false,
	})
	if err != nil {
		log.Fatalln("log", err)
	}

	commitsHandled := 0
	err = commitIter.ForEach(func(commit *object.Commit) error {
		commitsHandled++
		if commitsHandled > 10 {
			return storer.ErrStop
		}

		checkErr := check(commit)
		if checkErr != nil {
			log.Fatalln("check", checkErr)
		}

		return nil
	})
	if err != nil {
		log.Fatalln("for", err)
	}

	log.Println("finished")
}

const (
	maxTitleLen       = 40
	maxFileChangesLen = 400
)

func check(commit *object.Commit) error {
	stats, err := commit.Stats()
	if err != nil {
		return fmt.Errorf("stats %w", err)
	}

	title := strings.SplitN(commit.Message, "\n", 1)[0]

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
