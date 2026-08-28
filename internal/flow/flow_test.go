package flow

import (
	"path/filepath"
	"testing"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/store"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "notice.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewService(db, func() string { return "2026-01-01T00:00:00Z" })
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	service := testService(t)
	record, err := service.CreateRecord(10, "高处作业", []string{"佩戴安全带"}, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.ReviewRecord(record.ID, "bob"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.ConfirmRecord(record.ID, "bob"); err != nil {
		t.Fatal(err)
	}
	archived, err := service.ArchiveRecord(record.ID, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != domain.StatusArchived {
		t.Fatalf("status = %s", archived.Status)
	}
}
