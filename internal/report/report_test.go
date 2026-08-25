package report

import (
	"testing"

	"workshopnotice/internal/domain"
)

func TestReportFormatting(t *testing.T) {
	record, _ := domain.NewRecord("r", 3, "有限空间", []string{"监护到位"}, "fixed")
	events := []domain.AuditEvent{domain.NewAudit("a", "r", "created", "alice", "fixed", "ok")}
	summary := GenerateAuditSummary(events)
	if summary.Count != 1 || FormatSummary(summary) == "" {
		t.Fatalf("summary = %#v", summary)
	}
	if !stringsContains(Markdown(record, summary), "有限空间") {
		t.Fatal("markdown omitted title")
	}
}

func stringsContains(value, needle string) bool {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
