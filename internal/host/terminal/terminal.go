package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"bbcbasic-go/internal/host"
)

type TerminalHost struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func New() host.Host {
	return &TerminalHost{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (t *TerminalHost) Init() error {
	return nil
}

func (t *TerminalHost) Reset() error {
	return nil
}

func (t *TerminalHost) WriteChar(b byte) error {
	if err := t.writer.WriteByte(b); err != nil {
		return err
	}
	return t.writer.Flush()
}

func (t *TerminalHost) WriteString(s string) error {
	if _, err := t.writer.WriteString(s); err != nil {
		return err
	}
	return t.writer.Flush()
}

func (t *TerminalHost) ReadChar() (byte, error) {
	return t.reader.ReadByte()
}

func (t *TerminalHost) ReadLine(prompt string) (string, error) {
	if err := t.WriteString(prompt); err != nil {
		return "", err
	}

	line, err := t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimRight(line, "\r\n")
	return line, nil
}

func (t *TerminalHost) ClearScreen() error {
	_, err := fmt.Fprint(t.writer, "\033[2J\033[H")
	if err != nil {
		return err
	}
	return t.writer.Flush()
}
