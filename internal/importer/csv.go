package importer

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func ParseLines(data string) ([]Input, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	lineNumber := 0
	result := make([]Input, 0)
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			return nil, fmt.Errorf("line %d requires name|number|title|items", lineNumber)
		}
		number, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("line %d number: %w", lineNumber, err)
		}
		items := make([]string, 0)
		for _, item := range strings.Split(parts[3], ";") {
			if text := strings.TrimSpace(item); text != "" {
				items = append(items, text)
			}
		}
		input := Normalize(Input{Name: strings.TrimSpace(parts[0]), Number: number, Title: strings.TrimSpace(parts[2]), Items: items})
		if err = Validate(input); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		result = append(result, input)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func RenderLines(inputs []Input) string {
	lines := make([]string, 0, len(inputs))
	for _, input := range inputs {
		lines = append(lines, fmt.Sprintf("%s|%d|%s|%s", input.Name, input.Number, input.Title, strings.Join(input.Items, ";")))
	}
	return strings.Join(lines, "\n")
}

func MergeInputs(groups ...[]Input) []Input {
	result := make([]Input, 0)
	for _, group := range groups {
		for _, input := range group {
			result = append(result, Normalize(input))
		}
	}
	return result
}
