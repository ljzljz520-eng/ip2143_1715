package store

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type Maintenance struct {
	Records int
	Audits  int
	Oldest  string
	Newest  string
}

func (s *Store) MaintenanceReport() (Maintenance, error) {
	backup, err := s.BuildBackup()
	if err != nil {
		return Maintenance{}, err
	}
	stamps := make([]string, 0, len(backup.Records))
	for _, record := range backup.Records {
		stamps = append(stamps, record.UpdatedAt)
	}
	sort.Strings(stamps)
	result := Maintenance{Records: len(backup.Records), Audits: len(backup.Audits)}
	if len(stamps) > 0 {
		result.Oldest, result.Newest = stamps[0], stamps[len(stamps)-1]
	}
	return result, nil
}

func (m Maintenance) String() string {
	return fmt.Sprintf("records=%d audits=%d oldest=%s newest=%s", m.Records, m.Audits, m.Oldest, m.Newest)
}

func RecordsByStatus(records []domain.Record, status domain.Status) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if status == "" || record.Status == status {
			result = append(result, record)
		}
	}
	return result
}
