package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"go.etcd.io/bbolt"
	"workshopnotice/internal/domain"
)

var (
	bucketRecords     = []byte("records")
	bucketAudits      = []byte("audits")
	bucketWorkflows   = []byte("workflows")
	bucketAttachments = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		path = filepath.Join(".", "workshopnotice.db")
	}
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{bucketRecords, bucketAudits, bucketWorkflows, bucketAttachments} {
			if _, createErr := tx.CreateBucketIfNotExists(name); createErr != nil {
				return createErr
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty persisted value")
	}
	return json.Unmarshal(data, target)
}

func (s *Store) SaveRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketRecords).Put([]byte(record.ID), data) })
}

func (s *Store) LoadRecord(id string) (domain.Record, error) {
	var record domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error { return decode(tx.Bucket(bucketRecords).Get([]byte(id)), &record) })
	return record, err
}

func (s *Store) ListRecords() ([]domain.Record, error) {
	result := make([]domain.Record, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			result = append(result, record)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].Number < result[j].Number })
	return result, err
}

func (s *Store) SaveAudit(event domain.AuditEvent) error {
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketAudits).Put([]byte(event.ID), data) })
}

func (s *Store) ListAudits(recordID string) ([]domain.AuditEvent, error) {
	result := make([]domain.AuditEvent, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAudits).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if err := decode(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				result = append(result, event)
			}
			return nil
		})
	})
	return result, err
}
