package interpreter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type exprParser struct {
	interp *Interpreter
	input  string
	pos    int
}

func (i *Interpreter) evalIntExpression(expr string) (int, error) {
	parser := &exprParser{
		interp: i,
		input:  expr,
		pos:    0,
	}

	value, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}

	parser.skipSpaces()
	if !parser.isAtEnd() {
		return 0, fmt.Errorf("unexpected input near: %s", parser.remaining())
	}

	return value, nil
}

func (p *exprParser) parseExpression() (int, error) {
	value, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for {
		p.skipSpaces()

		if p.match('+') {
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			value += right
			continue
		}

		if p.match('-') {
			right, err := p.parseTerm()
			if err != nil {
				return 0, err
			}
			value -= right
			continue
		}

		break
	}

	return value, nil
}

func (p *exprParser) parseTerm() (int, error) {
	value, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	for {
		p.skipSpaces()

		if p.match('*') {
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			value *= right
			continue
		}

		if p.match('/') {
			right, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			value /= right
			continue
		}

		break
	}

	return value, nil
}

func (p *exprParser) parseFactor() (int, error) {
	p.skipSpaces()

	if p.match('-') {
		value, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -value, nil
	}

	if p.match('(') {
		value, err := p.parseExpression()
		if err != nil {
			return 0, err
		}

		p.skipSpaces()
		if !p.match(')') {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		return value, nil
	}

	if p.isAtEnd() {
		return 0, fmt.Errorf("unexpected end of expression")
	}

	ch := p.peek()

	if unicode.IsDigit(ch) {
		return p.parseNumber()
	}

	if unicode.IsLetter(ch) {
		return p.parseVariable()
	}

	return 0, fmt.Errorf("unexpected character: %c", ch)
}

func (p *exprParser) parseNumber() (int, error) {
	start := p.pos

	for !p.isAtEnd() && unicode.IsDigit(p.peek()) {
		p.pos++
	}

	text := p.input[start:p.pos]
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", text)
	}

	return value, nil
}

func (p *exprParser) parseVariable() (int, error) {
	start := p.pos

	for !p.isAtEnd() {
		ch := p.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '%' || ch == '$' {
			p.pos++
			continue
		}
		break
	}

	name := strings.ToUpper(strings.TrimSpace(p.input[start:p.pos]))
	if name == "" {
		return 0, fmt.Errorf("invalid variable")
	}

	value, ok := p.interp.runtime.GetVar(name)
	if !ok {
		return 0, fmt.Errorf("unknown variable: %s", name)
	}

	return value, nil
}

func (p *exprParser) skipSpaces() {
	for !p.isAtEnd() && unicode.IsSpace(p.peek()) {
		p.pos++
	}
}

func (p *exprParser) match(ch rune) bool {
	p.skipSpaces()
	if p.isAtEnd() || p.peek() != ch {
		return false
	}
	p.pos++
	return true
}

func (p *exprParser) peek() rune {
	return rune(p.input[p.pos])
}

func (p *exprParser) isAtEnd() bool {
	return p.pos >= len(p.input)
}

func (p *exprParser) remaining() string {
	if p.isAtEnd() {
		return ""
	}
	return strings.TrimSpace(p.input[p.pos:])
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
