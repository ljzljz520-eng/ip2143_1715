package report

import (
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/importer"
)

type Summary struct {
	RecordID string
	Count    int
	Actions  []string
	Latest   string
}

type ImportResult struct {
	Record     domain.Record
	Attachment domain.Attachment
	Message    string
}

func GenerateAuditSummary(events []domain.AuditEvent) Summary {
	summary := Summary{Actions: make([]string, 0, len(events))}
	for _, event := range events {
		if summary.RecordID == "" {
			summary.RecordID = event.RecordID
		}
		summary.Count++
		summary.Actions = append(summary.Actions, event.Action)
		summary.Latest = event.At
	}
	return summary
}

func GenerateImportReport(result ImportResult) string {
	return fmt.Sprintf("%s #%d %s: %d items", result.Message, result.Record.Number, result.Record.Title, len(result.Record.Items))
}

func ValidateImport(input importer.Input) error { return importer.Validate(importer.Normalize(input)) }

func FormatSummary(summary Summary) string {
	if summary.Count == 0 {
		return "no audit events"
	}
	return fmt.Sprintf("%s: %d events (%s), latest %s", summary.RecordID, summary.Count, strings.Join(summary.Actions, ","), summary.Latest)
}

func FilterActions(events []domain.AuditEvent, action string) []domain.AuditEvent {
	result := make([]domain.AuditEvent, 0)
	for _, event := range events {
		if action == "" || event.Action == action {
			result = append(result, event)
		}
	}
	return result
}
