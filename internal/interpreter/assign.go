package interpreter

import (
	"strings"
)

func (i *Interpreter) tryAssignment(text string) (bool, error) {
	eq := strings.Index(text, "=")
	if eq <= 0 {
		return false, nil
	}

	left := strings.TrimSpace(text[:eq])
	right := strings.TrimSpace(text[eq+1:])

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
