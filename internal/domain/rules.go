package domain

import (
	"fmt"
	"sort"
	"strings"
)

type Query struct {
	Number int
	Text   string
	Status Status
}

func MatchQuery(r Record, q Query) bool {
	if q.Number > 0 && r.Number != q.Number {
		return false
	}
	needle := strings.ToLower(strings.TrimSpace(q.Text))
	if needle != "" && !strings.Contains(strings.ToLower(r.Title), needle) {
		matched := false
		for _, item := range r.Items {
			if strings.Contains(strings.ToLower(item), needle) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if q.Status != "" && r.Status != q.Status {
		return false
	}
	return true
}

func SortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Number == records[j].Number {
			return records[i].ID < records[j].ID
		}
		return records[i].Number < records[j].Number
	})
}

func NormalizeItems(items []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func EnsureUniqueNumber(records []Record, number int, except string) error {
	for _, record := range records {
		if record.Number == number && record.ID != except {
			return fmt.Errorf("number %d is already used", number)
		}
	}
	return nil
}

func NextNumber(records []Record) int {
	next := 1
	for _, record := range records {
		if record.Number >= next {
			next = record.Number + 1
		}
	}
	return next
}
