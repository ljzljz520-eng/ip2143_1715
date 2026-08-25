package flow

import (
	"fmt"
	"sort"
	"strings"

	"workshopnotice/internal/domain"
)

type RetentionRule struct {
	ArchivedAfterVersion int
	KeepAuditActions     []string
}

func DefaultRetentionRule() RetentionRule {
	return RetentionRule{ArchivedAfterVersion: 1, KeepAuditActions: []string{"created", "confirmed", "archived"}}
}

func (rule RetentionRule) KeepRecord(record domain.Record) bool {
	if record.Status != domain.StatusArchived {
		return true
	}
	return record.Version <= rule.ArchivedAfterVersion+3
}

func (rule RetentionRule) KeepAudit(event domain.AuditEvent) bool {
	for _, action := range rule.KeepAuditActions {
		if event.Action == action {
			return true
		}
	}
	return false
}

func FilterRetained(records []domain.Record, rule RetentionRule) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if rule.KeepRecord(record) {
			result = append(result, record)
		}
	}
	return result
}

func FilterRetainedAudits(events []domain.AuditEvent, rule RetentionRule) []domain.AuditEvent {
	result := make([]domain.AuditEvent, 0)
	for _, event := range events {
		if rule.KeepAudit(event) {
			result = append(result, event)
		}
	}
	return result
}

func RetentionSummary(records []domain.Record, events []domain.AuditEvent, rule RetentionRule) string {
	keptRecords := FilterRetained(records, rule)
	keptEvents := FilterRetainedAudits(events, rule)
	return fmt.Sprintf("retained records=%d audits=%d", len(keptRecords), len(keptEvents))
}

func SortForArchive(records []domain.Record) {
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].Status == records[j].Status {
			return strings.Compare(records[i].UpdatedAt, records[j].UpdatedAt) < 0
		}
		return records[i].Status == domain.StatusArchived
	})
}
