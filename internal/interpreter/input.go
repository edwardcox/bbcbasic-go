package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *Interpreter) tryInput(text string) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if !strings.HasPrefix(upper, "INPUT") {
		return false, nil
	}

	rest := strings.TrimSpace(text[len("INPUT"):])
	if rest == "" {
		return true, fmt.Errorf("missing variable for INPUT")
	}

	if !isVariableName(rest) {
		return true, fmt.Errorf("invalid INPUT variable: %s", rest)
	}

	if err := i.host.WriteString("? "); err != nil {
		return true, err
	}

	line, err := i.host.ReadLineNoPrompt()
	if err != nil {
		return true, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return true, fmt.Errorf("empty INPUT")
	}

	value, err := strconv.Atoi(line)
	if err != nil {
		return true, fmt.Errorf("INPUT currently supports integers only")
	}

	i.runtime.SetVar(strings.ToUpper(rest), value)
	return true, nil
}
