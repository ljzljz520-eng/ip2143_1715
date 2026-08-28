package store

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type BucketStats struct {
	Name  string
	Count int
	Bytes int
}

func (s *Store) BucketStats() ([]BucketStats, error) {
	result := make([]BucketStats, 0, 4)
	err := s.db.View(func(tx *bbolt.Tx) error {
		buckets := map[string][]byte{"records": bucketRecords, "audits": bucketAudits, "workflows": bucketWorkflows, "attachments": bucketAttachments}
		for name, bucketName := range buckets {
			bucket := tx.Bucket(bucketName)
			stats := BucketStats{Name: name}
			err := bucket.ForEach(func(key, value []byte) error {
				if value == nil {
					return nil
				}
				stats.Count++
				stats.Bytes += len(key) + len(value)
				return nil
			})
			if err != nil {
				return err
			}
			result = append(result, stats)
		}
		return nil
	})
	return result, err
}

func (s *Store) EnsureRecord(id string) error {
	_, err := s.LoadRecord(id)
	if err != nil {
		return fmt.Errorf("record %s unavailable: %w", id, err)
	}
	return nil
}
