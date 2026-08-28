package store

import (
	"path/filepath"
	"testing"

	"workshopnotice/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persist.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record, err := domain.NewRecord("rec-1", 9, "应急出口", []string{"不得堵塞"}, "fixed")
	if err != nil {
		t.Fatal(err)
	}
	if err = db.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err = db.SaveAudit(domain.NewAudit("audit-1", record.ID, "created", "tester", "fixed", "ok")); err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	loaded, err := db.LoadRecord(record.ID)
	if err != nil || loaded.Title != record.Title || loaded.Number != 9 {
		t.Fatalf("loaded = %#v, err=%v", loaded, err)
	}
	audits, err := db.ListAudits(record.ID)
	if err != nil || len(audits) != 1 {
		t.Fatalf("audits = %d, err=%v", len(audits), err)
	}
}
