package store

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
	"workshopnotice/internal/domain"
)

type Backup struct {
	Records     []domain.Record
	Audits      []domain.AuditEvent
	Workflows   []domain.Workflow
	Attachments []domain.Attachment
}

func (s *Store) BuildBackup() (Backup, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Backup{}, err
	}
	audits, err := s.ListAudits("")
	if err != nil {
		return Backup{}, err
	}
	workflows, err := s.ListWorkflows("")
	if err != nil {
		return Backup{}, err
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return Backup{}, err
	}
	return Backup{Records: records, Audits: audits, Workflows: workflows, Attachments: attachments}, nil
}

func (s *Store) MarshalBackup() ([]byte, error) {
	backup, err := s.BuildBackup()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(backup, "", "  ")
}

func (s *Store) RestoreBackup(data []byte) error {
	var backup Backup
	if err := json.Unmarshal(data, &backup); err != nil {
		return fmt.Errorf("decode backup: %w", err)
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, record := range backup.Records {
			value, err := encode(record)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketRecords).Put([]byte(record.ID), value); err != nil {
				return err
			}
		}
		for _, event := range backup.Audits {
			value, err := encode(event)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketAudits).Put([]byte(event.ID), value); err != nil {
				return err
			}
		}
		for _, workflow := range backup.Workflows {
			value, err := encode(workflow)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketWorkflows).Put([]byte(workflow.ID), value); err != nil {
				return err
			}
		}
		for _, attachment := range backup.Attachments {
			value, err := encode(attachment)
			if err != nil {
				return err
			}
			if err = tx.Bucket(bucketAttachments).Put([]byte(attachment.ID), value); err != nil {
				return err
			}
		}
		return nil
	})
}
