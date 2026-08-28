package importer

import (
	"fmt"
	"strings"
)

func Validate(input Input) error {
	if input.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("title is required")
	}
	validItems := 0
	for _, item := range input.Items {
		if strings.TrimSpace(item) != "" {
			validItems++
		}
	}
	if validItems == 0 {
		return fmt.Errorf("items are required")
	}
	return nil
}

func Normalize(input Input) Input {
	input.Name = strings.TrimSpace(input.Name)
	input.Title = strings.TrimSpace(input.Title)
	result := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		if text := strings.TrimSpace(item); text != "" {
			result = append(result, text)
		}
	}
	input.Items = result
	return input
}
