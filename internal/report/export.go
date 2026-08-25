package report

import (
	"encoding/json"
	"sort"

	"workshopnotice/internal/domain"
)

type Export struct {
	Record  domain.Record
	Audits  []domain.AuditEvent
	Summary Summary
}

func BuildExport(record domain.Record, audits []domain.AuditEvent) Export {
	ordered := append([]domain.AuditEvent(nil), audits...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].At < ordered[j].At })
	return Export{Record: record, Audits: ordered, Summary: GenerateAuditSummary(ordered)}
}

func MarshalExport(export Export) ([]byte, error) { return json.MarshalIndent(export, "", "  ") }

func CountByAction(events []domain.AuditEvent) map[string]int {
	counts := map[string]int{}
	for _, event := range events {
		counts[event.Action]++
	}
	return counts
}
