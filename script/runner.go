package script

import (
	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/conf"
	"github.com/rusinikita/changes/errors"
)

type Runner interface {
	Run(commits []commit.Commit, diff []change.Change) error
}

const scriptRulesKey = "custom-rules"

type runner struct {
	scripts []*Script
}

func GetScriptsRunner(config conf.Conf, props value.Properties) (Runner, error) {
	var scriptConfigs []Config

	err := config.Unmarshal(scriptRulesKey, &scriptConfigs)
	if err != nil {
		return nil, errors.Prefix(err, scriptRulesKey)
	}

	env, err := createEnv(props)
	if err != nil {
		return nil, errors.Prefix(err, scriptRulesKey, "script env")
	}

	var scripts []*Script

	for _, c := range scriptConfigs {
		script, scriptErr := New(c, env)
		err = errors.Add(err, scriptErr, errors.StrToPathPrefix(c.Message))

		if script != nil {
			scripts = append(scripts, script)
		}
	}

	return &runner{scripts: scripts}, errors.Prefix(err, scriptRulesKey)
}

func (r runner) Run(commits []commit.Commit, changes []change.Change) (err error) {
	data := newData(commits, changes)

	for _, script := range r.scripts {
		result, runErr := script.Run(data)

		err = errors.Add(err, runErr, "script", errors.StrToPathPrefix(script.message))

		err = addValidationResult(err, script.message, result)
	}

	return err
}

func addValidationResult(err error, msg string, result any) error {
	switch v := result.(type) {
	case bool:
		if !v {
			return err
		}

		err = errors.Add(err, errors.New(msg))
	case []commit.Commit:
		for _, c := range v {
			err = errors.Add(err, errors.New(msg), c.ErrPrefix())
		}
	case []change.Change:
		for _, c := range v {
			err = errors.Add(err, errors.New(msg), c.ErrPrefix())
		}
	}

	return err
}
