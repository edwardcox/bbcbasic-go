package interpreter

import (
	"strings"

	"bbcbasic-go/internal/host"
)

type Interpreter struct {
	host host.Host
}

func New(h host.Host) *Interpreter {
	return &Interpreter{
		host: h,
	}
}

func (i *Interpreter) Run() error {
	if err := i.host.Init(); err != nil {
		return err
	}

	if err := i.host.WriteString("BBC BASIC for macOS (Go prototype)\n"); err != nil {
		return err
	}
	if err := i.host.WriteString("Type QUIT to exit\n"); err != nil {
		return err
	}

	for {
		line, err := i.host.ReadLine("> ")
		if err != nil {
			return err
		}

		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		switch upper {
		case "":
			continue

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

		default:
			if err := i.host.WriteString("You typed: " + line + "\n"); err != nil {
				return err
			}
		}
	}
}
