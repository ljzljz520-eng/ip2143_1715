package domain

import "testing"

func TestRecordRules(t *testing.T) {
	record, err := NewRecord("r", 1, "标题", []string{"一", "", "一", "二"}, "fixed")
	if err != nil {
		t.Fatal(err)
	}
	if len(record.Items) != 3 {
		t.Fatalf("constructor preserves input items: %d", len(record.Items))
	}
	if err = record.Transition(StatusConfirmed); err == nil {
		t.Fatal("invalid transition accepted")
	}
	if err = record.Transition(StatusReviewed); err != nil {
		t.Fatal(err)
	}
	if !CanEdit(record.Status) || CanArchive(record.Status) {
		t.Fatal("unexpected edit rules")
	}
}
