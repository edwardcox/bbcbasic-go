package interpreter

import (
	"fmt"
	"strings"
)

func (i *Interpreter) evalCondition(expr string) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false, fmt.Errorf("empty condition")
	}

	operators := []string{"<>", "<=", ">=", "=", "<", ">"}

	for _, op := range operators {
		if idx := strings.Index(expr, op); idx >= 0 {
			leftText := strings.TrimSpace(expr[:idx])
			rightText := strings.TrimSpace(expr[idx+len(op):])

			if leftText == "" || rightText == "" {
				return false, fmt.Errorf("invalid condition")
			}

			left, err := i.evalIntExpression(leftText)
			if err != nil {
				return false, err
			}

			right, err := i.evalIntExpression(rightText)
			if err != nil {
				return false, err
			}

			switch op {
			case "=":
				return left == right, nil
			case "<>":
				return left != right, nil
			case "<":
				return left < right, nil
			case "<=":
				return left <= right, nil
			case ">":
				return left > right, nil
			case ">=":
				return left >= right, nil
			}
		}
	}

	value, err := i.evalIntExpression(expr)
	if err != nil {
		return false, err
	}

	return value != 0, nil
}
