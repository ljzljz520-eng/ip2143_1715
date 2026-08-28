package importer

import (
	"fmt"
	"strings"
)

type Schema struct {
	Name     string
	Required []string
	Version  int
}

func DefaultSchema() Schema {
	return Schema{Name: "workshop-notice", Required: []string{"name", "number", "title", "items"}, Version: 1}
}

func (s Schema) ValidateFields(input Input) []string {
	missing := make([]string, 0)
	for _, field := range s.Required {
		present := true
		switch field {
		case "name":
			present = strings.TrimSpace(input.Name) != ""
		case "number":
			present = input.Number > 0
		case "title":
			present = strings.TrimSpace(input.Title) != ""
		case "items":
			present = len(input.Items) > 0
		default:
			present = false
		}
		if !present {
			missing = append(missing, field)
		}
	}
	return missing
}

func (s Schema) Message(input Input) string {
	missing := s.ValidateFields(input)
	if len(missing) == 0 {
		return fmt.Sprintf("%s v%d valid", s.Name, s.Version)
	}
	return fmt.Sprintf("%s missing %s", s.Name, strings.Join(missing, ","))
}

func Compatible(version int) bool { return version == DefaultSchema().Version }
