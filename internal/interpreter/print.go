package interpreter

import (
	"fmt"
	"strings"
)

func splitPrintItems(input string) ([]string, []rune, error) {
	var items []string
	var separators []rune

	var current strings.Builder
	inString := false

	flush := func() {
		items = append(items, strings.TrimSpace(current.String()))
		current.Reset()
	}

	for _, r := range input {
		switch r {
		case '"':
			inString = !inString
			current.WriteRune(r)

		case ';', ',':
			if inString {
				current.WriteRune(r)
				continue
			}
			flush()
			separators = append(separators, r)

		default:
			current.WriteRune(r)
		}
	}

	if inString {
		return nil, nil, fmt.Errorf("unterminated string in PRINT")
	}

	flush()

	return items, separators, nil
}
