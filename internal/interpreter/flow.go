package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *Interpreter) tryGoto(text string) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if !strings.HasPrefix(upper, "GOTO") {
		return false, nil
	}

	rest := strings.TrimSpace(text[len("GOTO"):])
	if rest == "" {
		return true, fmt.Errorf("missing line number for GOTO")
	}

	target, err := strconv.Atoi(rest)
	if err != nil || target <= 0 {
		return true, fmt.Errorf("invalid GOTO target: %s", rest)
	}

	i.runtime.SetJump(target)
	return true, nil
}

func (i *Interpreter) tryIfThen(text string) (bool, error) {
	upper := strings.ToUpper(text)
	if !strings.HasPrefix(strings.TrimSpace(upper), "IF ") {
		return false, nil
	}

	thenIdx := strings.Index(upper, "THEN")
	if thenIdx < 0 {
		return true, fmt.Errorf("missing THEN")
	}

	conditionText := strings.TrimSpace(text[:thenIdx])
	thenText := strings.TrimSpace(text[thenIdx+len("THEN"):])

	if len(conditionText) < 2 {
		return true, fmt.Errorf("invalid IF statement")
	}

	conditionExpr := strings.TrimSpace(conditionText[len("IF"):])
	if conditionExpr == "" {
		return true, fmt.Errorf("missing IF condition")
	}

	if thenText == "" {
		return true, fmt.Errorf("missing THEN target")
	}

	ok, err := i.evalCondition(conditionExpr)
	if err != nil {
		return true, err
	}

	if !ok {
		return true, nil
	}

	target, err := strconv.Atoi(thenText)
	if err != nil || target <= 0 {
		return true, fmt.Errorf("invalid THEN target: %s", thenText)
	}

	i.runtime.SetJump(target)
	return true, nil
}
