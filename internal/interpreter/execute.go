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
				return fmt.Errorf("line %d: target line not found: %d", line.Number, target)
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

	if upper == "END" {
		i.runtime.Stop()
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

	if strings.HasPrefix(rest, "\"") && strings.HasSuffix(rest, "\"") && len(rest) >= 2 {
		value := rest[1 : len(rest)-1]
		return i.host.WriteString(value + "\n")
	}

	value, err := i.evalIntExpression(rest)
	if err != nil {
		return err
	}

	return i.host.WriteString(strconv.Itoa(value) + "\n")
}
