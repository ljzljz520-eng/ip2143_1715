package flow

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/report"
)

type Snapshot struct {
	Records   []domain.Record
	Audits    []domain.AuditEvent
	Workflows []domain.Workflow
}

func (s *Service) BuildSnapshot() (Snapshot, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.Store.ListAudits("")
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.Store.ListWorkflows("")
	if err != nil {
		return Snapshot{}, err
	}
	sort.Slice(records, func(i, j int) bool { return records[i].UpdatedAt > records[j].UpdatedAt })
	return Snapshot{Records: records, Audits: audits, Workflows: workflows}, nil
}

func SnapshotLabel(snapshot Snapshot) string {
	if len(snapshot.Records) == 0 {
		return "empty workshop"
	}
	return fmt.Sprintf("%d records, %d audits, %d workflows", len(snapshot.Records), len(snapshot.Audits), len(snapshot.Workflows))
}

func (s *Service) ExportRecord(id string) ([]byte, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return nil, err
	}
	audits, err := s.Store.ListAudits(id)
	if err != nil {
		return nil, err
	}
	return report.MarshalExport(report.BuildExport(record, audits))
}
