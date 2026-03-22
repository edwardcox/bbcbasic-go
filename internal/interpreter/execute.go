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
	for _, line := range lines {
		if err := i.executeLine(line); err != nil {
			return fmt.Errorf("line %d: %w", line.Number, err)
		}
		if i.runtime.IsStopped() {
			return nil
		}
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

	return fmt.Errorf("unsupported statement: %s", text)
}

func (i *Interpreter) executePrint(text string) error {
	rest := strings.TrimSpace(text[len("PRINT"):])
	if rest == "" {
		return i.host.WriteString("\n")
	}

	// Quoted string: PRINT "HELLO"
	if strings.HasPrefix(rest, "\"") && strings.HasSuffix(rest, "\"") && len(rest) >= 2 {
		value := rest[1 : len(rest)-1]
		return i.host.WriteString(value + "\n")
	}

	// Integer number: PRINT 123
	if n, err := strconv.Atoi(rest); err == nil {
		return i.host.WriteString(strconv.Itoa(n) + "\n")
	}

	return fmt.Errorf("PRINT currently supports quoted strings or integers only")
}
