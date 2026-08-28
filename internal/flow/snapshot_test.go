package flow

import "testing"

func TestSnapshotAndExport(t *testing.T) {
	service := testService(t)
	record, err := service.CreateRecord(4, "锁具管理", []string{"专人保管"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := service.BuildSnapshot()
	if err != nil || len(snapshot.Records) != 1 {
		t.Fatalf("snapshot: %#v %v", snapshot, err)
	}
	if SnapshotLabel(snapshot) == "" {
		t.Fatal("empty snapshot label")
	}
	data, err := service.ExportRecord(record.ID)
	if err != nil || len(data) == 0 {
		t.Fatalf("export: %v", err)
	}
}
