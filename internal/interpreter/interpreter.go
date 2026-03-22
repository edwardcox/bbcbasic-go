package interpreter

import (
	"strings"

	"bbcbasic-go/internal/host"
	"bbcbasic-go/internal/program"
	"bbcbasic-go/internal/runtime"
)

type Interpreter struct {
	host    host.Host
	program *program.Program
	runtime *runtime.Runtime
}

func New(h host.Host) *Interpreter {
	return &Interpreter{
		host:    h,
		program: program.New(),
		runtime: runtime.New(),
	}
}

func (i *Interpreter) Run() error {
	if err := i.host.Init(); err != nil {
		return err
	}

	if err := i.host.WriteString("BBC BASIC for macOS (Go prototype)\n"); err != nil {
		return err
	}
	if err := i.host.WriteString("Type QUIT to exit. Use LIST, NEW, RUN, or numbered lines.\n"); err != nil {
		return err
	}

	for {
		line, err := i.host.ReadLine("> ")
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		if trimmed == "" {
			continue
		}

		lineNumber, lineText, isNumberedLine, err := program.ParseNumberedLine(trimmed)
		if err != nil {
			if err := i.host.WriteString("Error: " + err.Error() + "\n"); err != nil {
				return err
			}
			continue
		}

		if isNumberedLine {
			i.program.SetLine(lineNumber, lineText)
			continue
		}

		switch upper {
		case "QUIT", "EXIT":
			if err := i.host.WriteString("Bye.\n"); err != nil {
				return err
			}
			return nil

		case "CLS", "CLEAR":
			if err := i.host.ClearScreen(); err != nil {
				return err
			}
			continue

		case "LIST":
			listing := i.program.List()
			if listing == "" {
				continue
			}
			if err := i.host.WriteString(listing); err != nil {
				return err
			}
			continue

		case "NEW":
			i.program.Clear()
			continue

		case "RUN":
			if !i.program.HasLines() {
				if err := i.host.WriteString("No program\n"); err != nil {
					return err
				}
				continue
			}
			if err := i.runProgram(); err != nil {
				if err := i.host.WriteString("Error: " + err.Error() + "\n"); err != nil {
					return err
				}
			}
			continue

		default:
			if err := i.host.WriteString("You typed: " + line + "\n"); err != nil {
				return err
			}
		}
	}
}
