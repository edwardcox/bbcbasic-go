package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func (i *Interpreter) evalIntExpression(expr string) (int, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, fmt.Errorf("empty expression")
	}

	tokens := tokenizeSimpleExpr(expr)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}

	value, err := i.evalSimpleToken(tokens[0])
	if err != nil {
		return 0, err
	}

	pos := 1
	for pos < len(tokens) {
		if pos+1 >= len(tokens) {
			return 0, fmt.Errorf("incomplete expression")
		}

		op := tokens[pos]
		right, err := i.evalSimpleToken(tokens[pos+1])
		if err != nil {
			return 0, err
		}

		switch op {
		case "+":
			value += right
		case "-":
			value -= right
		case "*":
			value *= right
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			value /= right
		default:
			return 0, fmt.Errorf("unsupported operator: %s", op)
		}

		pos += 2
	}

	return value, nil
}

func (i *Interpreter) evalSimpleToken(token string) (int, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, fmt.Errorf("empty token")
	}

	if n, err := strconv.Atoi(token); err == nil {
		return n, nil
	}

	if isVariableName(token) {
		if value, ok := i.runtime.GetVar(strings.ToUpper(token)); ok {
			return value, nil
		}
		return 0, fmt.Errorf("unknown variable: %s", token)
	}

	return 0, fmt.Errorf("invalid token: %s", token)
}

func tokenizeSimpleExpr(expr string) []string {
	var tokens []string
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}

	for _, r := range expr {
		switch {
		case unicode.IsSpace(r):
			flush()

		case r == '+' || r == '-' || r == '*' || r == '/':
			flush()
			tokens = append(tokens, string(r))

		default:
			current.WriteRune(r)
		}
	}

	flush()
	return tokens
}

func isVariableName(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return false
			}
			continue
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '%' || r == '$' {
			continue
		}
		return false
	}

	return true
}
