package interpreter

import (
	"fmt"
	"strings"

	"bbcbasic-go/internal/files"
)

func parseQuotedFilename(text string, keyword string) (string, error) {
	rest := strings.TrimSpace(text[len(keyword):])
	if rest == "" {
		return "", fmt.Errorf("missing filename")
	}

	if !strings.HasPrefix(rest, "\"") || !strings.HasSuffix(rest, "\"") || len(rest) < 2 {
		return "", fmt.Errorf("filename must be in quotes")
	}

	filename := strings.TrimSpace(rest[1 : len(rest)-1])
	if filename == "" {
		return "", fmt.Errorf("empty filename")
	}

	return filename, nil
}

func (i *Interpreter) trySave(text string) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if !strings.HasPrefix(upper, "SAVE") {
		return false, nil
	}

	filename, err := parseQuotedFilename(strings.TrimSpace(text), "SAVE")
	if err != nil {
		return true, err
	}

	content := i.program.List()
	if err := files.SaveTextFile(filename, content); err != nil {
		return true, err
	}

	if err := i.host.WriteString("Saved " + filename + "\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (i *Interpreter) tryLoad(text string) (bool, error) {
	upper := strings.ToUpper(strings.TrimSpace(text))
	if !strings.HasPrefix(upper, "LOAD") {
		return false, nil
	}

	filename, err := parseQuotedFilename(strings.TrimSpace(text), "LOAD")
	if err != nil {
		return true, err
	}

	content, err := files.LoadTextFile(filename)
	if err != nil {
		return true, err
	}

	if err := i.program.FromText(content); err != nil {
		return true, err
	}

	if err := i.host.WriteString("Loaded " + filename + "\n"); err != nil {
		return true, err
	}

	return true, nil
}
