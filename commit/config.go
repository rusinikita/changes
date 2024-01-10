package commit

import (
	"github.com/rusinikita/changes/commit/subject"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/conf"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
)

const (
	conventionalCommitFormat = "(type)((context))?(!)?: (title)"
	formatKey                = "commit.subject"
)

type Parser struct {
	value.Properties
	*subject.Parser
}

func GetParser(config conf.Conf) (*Parser, error) {
	properties, err := value.Get(config)

	format := config.GetString(formatKey, conventionalCommitFormat)

	parser, formatErr := subject.NewParser(format, properties)

	err = errors.Add(err, errors.Prefix(formatErr, formatKey))
	if err != nil {
		return nil, err
	}

	return &Parser{
		Properties: properties,
		Parser:     parser,
	}, nil
}

func (p Parser) Parse(commits []git.Commit) (enrichedCommits []Commit, err error) {
	for _, c := range commits {
		values, commitErr := p.Parser.Parse(c.Subject())

		enriched := Commit{
			Commit: c,
			Values: values,
		}

		err = errors.Add(err, commitErr, enriched.ErrPrefix())

		enrichedCommits = append(enrichedCommits, enriched)
	}

	return enrichedCommits, err
}
