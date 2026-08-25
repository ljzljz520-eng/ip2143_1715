package importer

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Input struct {
	Name   string   `json:"name"`
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Items  []string `json:"items"`
}

func Parse(payload string) (Input, error) {
	var input Input
	if strings.TrimSpace(payload) == "" {
		return input, fmt.Errorf("import payload is empty")
	}
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		return input, fmt.Errorf("parse import payload: %w", err)
	}
	return input, nil
}

func Canonical(input Input) string {
	data, _ := json.Marshal(input)
	return string(data)
}

func Fields(input Input) map[string]string {
	return map[string]string{"name": input.Name, "number": fmt.Sprint(input.Number), "title": input.Title, "items": fmt.Sprint(len(input.Items))}
}
