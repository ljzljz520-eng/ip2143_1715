package store

import (
	"strings"

	"go.etcd.io/bbolt"
	"workshopnotice/internal/domain"
)

func (s *Store) FindByTitle(text string) ([]domain.Record, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(strings.TrimSpace(text))
	result := make([]domain.Record, 0)
	for _, record := range records {
		if needle == "" || strings.Contains(strings.ToLower(record.Title), needle) {
			result = append(result, record)
		}
	}
	return result, nil
}

func (s *Store) ReplaceRecord(record domain.Record) error {
	return s.SaveRecord(record)
}

func (s *Store) RemoveRecord(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketRecords).Delete([]byte(id)) })
}
