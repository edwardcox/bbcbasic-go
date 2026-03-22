package interpreter

import (
	"strings"
)

func (i *Interpreter) tryAssignment(text string) (bool, error) {
	trimmed := strings.TrimSpace(text)
	upper := strings.ToUpper(trimmed)

	if strings.HasPrefix(upper, "LET ") {
		trimmed = strings.TrimSpace(trimmed[4:])
	}

	eq := strings.Index(trimmed, "=")
	if eq <= 0 {
		return false, nil
	}

	left := strings.TrimSpace(trimmed[:eq])
	right := strings.TrimSpace(trimmed[eq+1:])

	if !isVariableName(left) {
		return false, nil
	}

	value, err := i.evalIntExpression(right)
	if err != nil {
		return true, err
	}

	i.runtime.SetVar(strings.ToUpper(left), value)
	return true, nil
}
