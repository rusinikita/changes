package output

import (
	"strings"

	"github.com/rusinikita/changes/errors"
)

func TerminalOutput(tree errors.OutputTree) string {
	sb := strings.Builder{}

	for _, root := range tree {
		sb.WriteString(root.Name)
		sb.WriteString("\n-------")

		sb.WriteString(messages(root.Messages) + "\n")

		for _, group := range root.Groups {
			sb.WriteString(group.Name)

			sb.WriteString(messages(group.Messages) + "\n")
		}
	}

	return strings.TrimSpace(sb.String())
}

func messages(mm []string) string {
	if len(mm) == 0 {
		return ""
	}

	return "\n- " + strings.Join(mm, "\n- ") + "\n"
}
