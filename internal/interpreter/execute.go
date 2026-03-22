package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"bbcbasic-go/internal/program"
)

func (i *Interpreter) runProgram() error {
	i.runtime.Reset()

	lines := i.program.SortedLines()
	if len(lines) == 0 {
		return nil
	}

	lineIndexByNumber := make(map[int]int, len(lines))
	for idx, line := range lines {
		lineIndexByNumber[line.Number] = idx
	}

	pc := 0
	for pc < len(lines) {
		line := lines[pc]

		if err := i.executeLine(line); err != nil {
			return fmt.Errorf("line %d: %w", line.Number, err)
		}

		if i.runtime.IsStopped() {
			return nil
		}

		if target, ok := i.runtime.ConsumeJump(); ok {
			nextPC, exists := lineIndexByNumber[target]
			if !exists {
				return fmt.Errorf("target line not found: %d", target)
			}
			pc = nextPC
			continue
		}

		pc++
	}

	return nil
}

func (i *Interpreter) executeLine(line program.Line) error {
	text := strings.TrimSpace(line.Text)
	if text == "" {
		return nil
	}

	upper := strings.ToUpper(text)

	if upper == "END" || upper == "STOP" {
		i.runtime.Stop()
		return nil
	}

	if strings.HasPrefix(upper, "REM") {
		return nil
	}

	if strings.HasPrefix(upper, "PRINT") {
		return i.executePrint(text)
	}

	if handled, err := i.tryIfThen(text); handled {
		return err
	}

	if handled, err := i.tryGoto(text); handled {
		return err
	}

	if handled, err := i.tryInput(text); handled {
		return err
	}

	if handled, err := i.tryAssignment(text); handled {
		return err
	}

	return fmt.Errorf("unsupported statement: %s", text)
}

func (i *Interpreter) executePrint(text string) error {
	rest := strings.TrimSpace(text[len("PRINT"):])
	if rest == "" {
		return i.host.WriteString("\n")
	}

	items, separators, err := splitPrintItems(rest)
	if err != nil {
		return err
	}

	var b strings.Builder

	for idx, item := range items {
		value, err := i.evalPrintItem(item)
		if err != nil {
			return err
		}
		b.WriteString(value)

		if idx < len(separators) {
			switch separators[idx] {
			case ';':
				// no extra spacing
			case ',':
				b.WriteString(" ")
			}
		}
	}

	b.WriteString("\n")
	return i.host.WriteString(b.String())
}

func (i *Interpreter) evalPrintItem(item string) (string, error) {
	item = strings.TrimSpace(item)
	if item == "" {
		return "", nil
	}

	if strings.HasPrefix(item, "\"") && strings.HasSuffix(item, "\"") && len(item) >= 2 {
		return item[1 : len(item)-1], nil
	}

	value, err := i.evalIntExpression(item)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(value), nil
}
