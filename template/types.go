package template

import "strings"

type Label = string

type Parameters = []uint64

type TransactionBody string

type TransactionEdit struct {
	PrepareBlock      string
	ExecuteBlock      string
	FieldDeclarations string
}

type Template interface {
	Name() string
	Label() Label

	// Cardinality returns the number of parameters that this template has
	Cardinality() uint

	TransactionEditFunc(parameters Parameters) (TransactionEdit, error)
	InitialParameters() Parameters
}

type Registry interface {
	Get(label Label) (Template, error)
	AllLabels() []Label
}

func TrimAndReplaceIndentation(s string, indentation int) string {
	// convert all tabs to spaces
	s = strings.ReplaceAll(s, "\t", "    ")

	// split into lines
	lines := strings.Split(s, "\n")
	// remove leading and trailing lines with only whitespace
	for len(lines) > 0 && len(strings.TrimSpace(lines[0])) == 0 {
		lines = lines[1:]
	}
	for len(lines) > 0 && len(strings.TrimSpace(lines[len(lines)-1])) == 0 {
		lines = lines[:len(lines)-1]
	}

	// get minimum spaces in all non-empty lines
	minSpaces := -1
	for _, line := range lines {
		spaces := 0
		for _, c := range line {
			if c == ' ' {
				spaces++
			} else {
				break
			}
		}
		if minSpaces == -1 || spaces < minSpaces {
			minSpaces = spaces
		}

	}

	// replace minSpaces from all non-empty lines with indentation spaces
	indentationStr := strings.Repeat(" ", indentation)
	for i, line := range lines {
		if len(line) > 0 {
			lines[i] = indentationStr + line[minSpaces:]
		}
	}
	return strings.Join(lines, "\n") + "\n"
}
