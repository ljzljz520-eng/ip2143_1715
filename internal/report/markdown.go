package report

import (
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
)

func Markdown(record domain.Record, summary Summary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n编号: %d\n\n状态: %s\n\n", record.Title, record.Number, domain.StatusLabel(record.Status))
	for i, item := range record.Items {
		fmt.Fprintf(&b, "%d. %s\n", i+1, item)
	}
	fmt.Fprintf(&b, "\n审计事件: %d\n", summary.Count)
	return b.String()
}

func Timeline(events []domain.AuditEvent) []string {
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, fmt.Sprintf("%s %s by %s", event.At, event.Action, event.Actor))
	}
	return result
}

func IsComplete(summary Summary) bool {
	for _, action := range summary.Actions {
		if action == "archived" {
			return true
		}
	}
	return false
}
