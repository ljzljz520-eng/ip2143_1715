package domain

import (
	"fmt"
	"strings"
)

type Checklist struct {
	Name     string
	Required []string
	Optional []string
}

func NewChecklist(name string, required, optional []string) Checklist {
	return Checklist{Name: strings.TrimSpace(name), Required: NormalizeItems(required), Optional: NormalizeItems(optional)}
}

func (c Checklist) Validate(items []string) []string {
	provided := map[string]bool{}
	for _, item := range items {
		provided[strings.ToLower(strings.TrimSpace(item))] = true
	}
	missing := make([]string, 0)
	for _, item := range c.Required {
		if !provided[strings.ToLower(item)] {
			missing = append(missing, item)
		}
	}
	return missing
}

func (c Checklist) Completion(items []string) float64 {
	total := len(c.Required) + len(c.Optional)
	if total == 0 {
		return 1
	}
	missing := len(c.Validate(items))
	return float64(total-missing) / float64(total)
}

func (c Checklist) Message(items []string) string {
	missing := c.Validate(items)
	if len(missing) == 0 {
		return fmt.Sprintf("%s complete", c.Name)
	}
	return fmt.Sprintf("%s missing: %s", c.Name, strings.Join(missing, ", "))
}
