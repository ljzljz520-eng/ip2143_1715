package flow022

import (
	"path/filepath"
	"testing"

	"workshopnotice/internal/flow"
	"workshopnotice/internal/store"
)

func Test2143BusinessRegression(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "regression.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	service := flow.NewService(db, func() string { return "2026-01-01T00:00:00Z" })
	records, err := service.SubmitNoticeBatch([]flow.NoticeInput{{Number: 101, Title: "第一项", Items: []string{"戴好护具"}}, {Number: 202, Title: "第二项", Items: []string{"确认警示牌"}}}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("records = %d", len(records))
	}
	if records[1].Number != 202 {
		t.Fatalf("second notice number = %d, want 202", records[1].Number)
	}
}
