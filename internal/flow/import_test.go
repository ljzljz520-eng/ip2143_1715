package flow

import (
	"testing"

	"workshopnotice/internal/report"
)

func TestWorkflowImportReport(t *testing.T) {
	service := testService(t)
	result, err := service.ImportNotice(`{"name":"shift-a.json","number":7,"title":"设备点检","items":["断电检查"]}`, "importer")
	if err != nil {
		t.Fatal(err)
	}
	if result.Attachment.Name != "shift-a.json" || result.Record.Number != 7 {
		t.Fatalf("unexpected import: %#v", result)
	}
	text := report.GenerateImportReport(result)
	if text == "" {
		t.Fatal("empty report")
	}
	attachments, err := service.Store.ListAttachments(result.Record.ID)
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments: %d %v", len(attachments), err)
	}
}
