package program

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	Number int
	Text   string
}

type Program struct {
	lines map[int]string
}

func New() *Program {
	return &Program{
		lines: make(map[int]string),
	}
}

func (p *Program) Clear() {
	p.lines = make(map[int]string)
}

func (p *Program) SetLine(number int, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		delete(p.lines, number)
		return
	}
	p.lines[number] = text
}

func (p *Program) DeleteLine(number int) {
	delete(p.lines, number)
}

func (p *Program) HasLines() bool {
	return len(p.lines) > 0
}

func (p *Program) SortedLines() []Line {
	numbers := make([]int, 0, len(p.lines))
	for n := range p.lines {
		numbers = append(numbers, n)
	}
	sort.Ints(numbers)

	result := make([]Line, 0, len(numbers))
	for _, n := range numbers {
		result = append(result, Line{
			Number: n,
			Text:   p.lines[n],
		})
	}
	return result
}

func (p *Program) List() string {
	var b strings.Builder

	for _, line := range p.SortedLines() {
		b.WriteString(strconv.Itoa(line.Number))
		if line.Text != "" {
			b.WriteString(" ")
			b.WriteString(line.Text)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func ParseNumberedLine(input string) (lineNumber int, lineText string, ok bool, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, "", false, nil
	}

	i := 0
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		i++
	}

	if i == 0 {
		return 0, "", false, nil
	}

	numberText := input[:i]
	number, convErr := strconv.Atoi(numberText)
	if convErr != nil {
		return 0, "", false, fmt.Errorf("invalid line number: %w", convErr)
	}
	if number <= 0 {
		return 0, "", false, fmt.Errorf("line number must be greater than zero")
	}

	rest := strings.TrimSpace(input[i:])
	return number, rest, true, nil
}
