package output

import (
	"strings"

	"github.com/rusinikita/changes/errors"
)

func MarkdownOutput(tree errors.OutputTree) string {
	sb := strings.Builder{}

	sb.WriteString("# Changes report\n")

	for _, root := range tree {
		sb.WriteString("## " + root.Name)

		sb.WriteString(messages(root.Messages) + "\n")

		for _, group := range root.Groups {
			sb.WriteString("**" + group.Name + "**")

			sb.WriteString(messages(group.Messages) + "\n")
		}
	}

	if len(tree) == 0 {
		sb.WriteString("OK")
	}

	return strings.TrimSpace(sb.String())
}
