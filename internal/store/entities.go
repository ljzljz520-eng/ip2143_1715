package store

import (
	"sort"

	"go.etcd.io/bbolt"
	"workshopnotice/internal/domain"
)

func (s *Store) SaveWorkflow(workflow domain.Workflow) error {
	data, err := encode(workflow)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketWorkflows).Put([]byte(workflow.ID), data) })
}

func (s *Store) ListWorkflows(recordID string) ([]domain.Workflow, error) {
	result := make([]domain.Workflow, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketWorkflows).ForEach(func(_, value []byte) error {
			var workflow domain.Workflow
			if err := decode(value, &workflow); err != nil {
				return err
			}
			if recordID == "" || workflow.RecordID == recordID {
				result = append(result, workflow)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) SaveAttachment(attachment domain.Attachment) error {
	data, err := encode(attachment)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketAttachments).Put([]byte(attachment.ID), data) })
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	result := make([]domain.Attachment, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketAttachments).ForEach(func(_, value []byte) error {
			var attachment domain.Attachment
			if err := decode(value, &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				result = append(result, attachment)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) CountAll() (map[string]int, error) {
	counts := map[string]int{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for name, bucket := range map[string][]byte{"records": bucketRecords, "audits": bucketAudits, "workflows": bucketWorkflows, "attachments": bucketAttachments} {
			count := 0
			if err := tx.Bucket(bucket).ForEach(func(_, _ []byte) error { count++; return nil }); err != nil {
				return err
			}
			counts[name] = count
		}
		return nil
	})
	return counts, err
}
