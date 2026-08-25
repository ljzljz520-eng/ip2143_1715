package report

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type Statistics struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	ByActor  map[string]int `json:"by_actor"`
}

func CalculateStatistics(records []domain.Record, events []domain.AuditEvent) Statistics {
	result := Statistics{Total: len(records), ByStatus: map[string]int{}, ByActor: map[string]int{}}
	for _, record := range records {
		result.ByStatus[string(record.Status)]++
	}
	for _, event := range events {
		result.ByActor[event.Actor]++
	}
	return result
}

func (s Statistics) Statuses() []string {
	keys := make([]string, 0, len(s.ByStatus))
	for key := range s.ByStatus {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s Statistics) String() string {
	return fmt.Sprintf("total=%d statuses=%v actors=%d", s.Total, s.Statuses(), len(s.ByActor))
}

func CompletionRate(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	complete := 0
	for _, record := range records {
		if record.Status == domain.StatusArchived {
			complete++
		}
	}
	return float64(complete) / float64(len(records))
}
