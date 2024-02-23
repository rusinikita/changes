package errors

const (
	CommitGroup    = "commit"
	FilesDiffGroup = "diff"
)

type OutputTree map[string]map[string][]string

func multiToTree(e multiErr) OutputTree {
	r := OutputTree{}

	for _, err := range e.errs {
		first, second, err := err.cutRoot()

		if r[first] == nil {
			r[first] = map[string][]string{}
		}

		r[first][second] = append(r[first][second], err.Error())
	}

	return r
}
