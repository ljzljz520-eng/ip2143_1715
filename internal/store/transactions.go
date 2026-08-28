package store

import (
	"fmt"

	"go.etcd.io/bbolt"
	"workshopnotice/internal/domain"
)

func (s *Store) SaveRecordWithAudit(record domain.Record, event domain.AuditEvent) error {
	if err := record.Validate(); err != nil {
		return err
	}
	recordData, err := encode(record)
	if err != nil {
		return err
	}
	eventData, err := encode(event)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(bucketRecords).Put([]byte(record.ID), recordData); err != nil {
			return fmt.Errorf("record transaction: %w", err)
		}
		if err := tx.Bucket(bucketAudits).Put([]byte(event.ID), eventData); err != nil {
			return fmt.Errorf("audit transaction: %w", err)
		}
		return nil
	})
}

func (s *Store) DeleteRecordAndChildren(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(bucketRecords).Delete([]byte(id)); err != nil {
			return err
		}
		for _, bucketName := range [][]byte{bucketAudits, bucketWorkflows, bucketAttachments} {
			bucket := tx.Bucket(bucketName)
			keys := make([][]byte, 0)
			if err := bucket.ForEach(func(key, value []byte) error {
				if value != nil {
					return nil
				}
				keys = append(keys, append([]byte(nil), key...))
				return nil
			}); err != nil {
				return err
			}
			for _, key := range keys {
				if err := bucket.Delete(key); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
