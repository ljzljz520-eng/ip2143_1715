package flow

import (
	"testing"

	"workshopnotice/internal/domain"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	service := testService(t)
	a, err := service.CreateRecord(2, "消防通道", []string{"保持畅通"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = service.CreateRecord(3, "叉车作业", []string{"减速"}, "operator"); err != nil {
		t.Fatal(err)
	}
	found, err := service.SearchRecords(domain.Query{Text: "消防"})
	if err != nil || len(found) != 1 || found[0].ID != a.ID {
		t.Fatalf("unexpected search: %#v %v", found, err)
	}
	updated, err := service.UpdateRecord(a.ID, []string{"保持畅通", "禁止堆物"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 {
		t.Fatalf("version = %d", updated.Version)
	}
	published, err := service.PublishRecord(a.ID, "reviewer")
	if err != nil || published.Status != domain.StatusConfirmed {
		t.Fatalf("publish: %#v %v", published, err)
	}
}
