package interpreter

import (
	"bbcbasic-go/internal/files"
)

func (i *Interpreter) LoadProgramFromFile(filename string) error {
	content, err := files.LoadTextFile(filename)
	if err != nil {
		return err
	}

	return i.program.FromText(content)
}

func (i *Interpreter) RunLoadedProgram() error {
	if !i.program.HasLines() {
		return nil
	}

	return i.runProgram()
}
