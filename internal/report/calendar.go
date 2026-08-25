package report

import (
	"fmt"
	"sort"
	"strings"

	"workshopnotice/internal/domain"
)

type DaySummary struct {
	Day     string
	Created int
	Updated int
	Closed  int
}

func BuildCalendar(events []domain.AuditEvent) []DaySummary {
	byDay := map[string]*DaySummary{}
	for _, event := range events {
		day := event.At
		if len(day) >= 10 {
			day = day[:10]
		}
		summary := byDay[day]
		if summary == nil {
			summary = &DaySummary{Day: day}
			byDay[day] = summary
		}
		switch event.Action {
		case "created", "imported", "batch-created":
			summary.Created++
		case "updated":
			summary.Updated++
		case "archived":
			summary.Closed++
		}
	}
	result := make([]DaySummary, 0, len(byDay))
	for _, summary := range byDay {
		result = append(result, *summary)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Day < result[j].Day })
	return result
}

func CalendarLine(summary DaySummary) string {
	return fmt.Sprintf("%s created=%d updated=%d closed=%d", summary.Day, summary.Created, summary.Updated, summary.Closed)
}

func SearchCalendar(summaries []DaySummary, prefix string) []DaySummary {
	result := make([]DaySummary, 0)
	for _, summary := range summaries {
		if strings.HasPrefix(summary.Day, prefix) {
			result = append(result, summary)
		}
	}
	return result
}

func ActionTotal(summary DaySummary) int { return summary.Created + summary.Updated + summary.Closed }
