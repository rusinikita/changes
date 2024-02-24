package errors

import (
	"strings"

	"golang.org/x/exp/slices"
)

const (
	CommitGroup    = "Commits"
	FilesDiffGroup = "Changes"
)

type OutputTree []Node

type Node struct {
	Name     string
	Groups   []Node
	Messages []string
}

func PrepareOutput(err error) OutputTree {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *multiErr:
		return multiToTree(*e)
	case *pathErr:
		return multiToTree(multiErr{errs: []pathErr{*e}})
	}

	return OutputTree{{
		Name:     "Error",
		Messages: []string{err.Error()},
	}}
}

func multiToTree(e multiErr) (result OutputTree) {
	tree := map[string]map[string][]string{}

	for _, err := range e.errs {
		first, second, err := err.cutRoot()

		if tree[first] == nil {
			tree[first] = map[string][]string{}
		}

		tree[first][second] = append(tree[first][second], err.Error())
	}

	for root, groups := range tree {
		node := Node{
			Name: root,
		}

		for name, issues := range groups {
			slices.Sort(issues)

			if name == "" {
				node.Messages = issues
				continue
			}

			node.Groups = append(node.Groups, Node{
				Name:     name,
				Messages: issues,
			})
		}

		slices.SortFunc(node.Groups, compare)

		result = append(result, node)
	}

	slices.SortFunc(result, compare)

	return result
}

func compare(a, b Node) int {
	return strings.Compare(a.Name, b.Name)
}
