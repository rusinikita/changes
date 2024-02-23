package output

import (
	"strings"

	"github.com/rusinikita/changes/errors"
)

func MarkdownOutput(tree errors.OutputTree) string {
	sb := strings.Builder{}

	for _, root := range tree {
		sb.WriteString("## " + root.Name)

		sb.WriteString(messages(root.Messages) + "\n")

		for _, group := range root.Groups {
			sb.WriteString(group.Name)

			sb.WriteString(messages(group.Messages) + "\n")
		}
	}

	return strings.TrimSpace(sb.String())
}
