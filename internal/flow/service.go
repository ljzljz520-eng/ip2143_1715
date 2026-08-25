package flow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/report"
	"workshopnotice/internal/store"
)

type Stamp func() string

type Service struct {
	Store *store.Store
	Now   Stamp
	Seq   int
	Nums  *NumberProcessor
}

func NewService(db *store.Store, now Stamp) *Service {
	if now == nil {
		now = func() string { return "2026-01-01T00:00:00Z" }
	}
	return &Service{Store: db, Now: now, Nums: NewNumberProcessor()}
}

func (s *Service) nextID(prefix string) string {
	s.Seq++
	return fmt.Sprintf("%s-%03d", prefix, s.Seq)
}

func (s *Service) CreateRecord(number int, title string, items []string, actor string) (domain.Record, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return domain.Record{}, err
	}
	if number == 0 {
		number = domain.NextNumber(records)
	}
	if err := domain.EnsureUniqueNumber(records, number, ""); err != nil {
		return domain.Record{}, err
	}
	record, err := domain.NewRecord(s.nextID("rec"), number, title, domain.NormalizeItems(items), s.Now())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err := s.recordEvent(record, "created", actor, "record registered"); err != nil {
		return domain.Record{}, err
	}
	workflow := domain.NewWorkflow(s.nextID("wf"), record.ID, "lifecycle", "registered", s.Now())
	if err := s.Store.SaveWorkflow(workflow); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) ReviewRecord(id, actor string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err = record.Transition(domain.StatusReviewed); err != nil {
		return domain.Record{}, err
	}
	record.UpdatedAt = s.Now()
	if err = s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err = s.recordEvent(record, "reviewed", actor, "review approved"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) ConfirmRecord(id, actor string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err = record.Transition(domain.StatusConfirmed); err != nil {
		return domain.Record{}, err
	}
	record.UpdatedAt = s.Now()
	if err = s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err = s.recordEvent(record, "confirmed", actor, "notice published"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) ArchiveRecord(id, actor string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status != domain.StatusConfirmed {
		return domain.Record{}, domain.ErrTransition
	}
	if err = record.Transition(domain.StatusArchived); err != nil {
		return domain.Record{}, err
	}
	record.UpdatedAt = s.Now()
	if err = s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err = s.recordEvent(record, "archived", actor, "notice archived"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) SearchRecords(query domain.Query) ([]domain.Record, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if domain.MatchQuery(record, query) {
			result = append(result, record)
		}
	}
	domain.SortRecords(result)
	return result, nil
}

func (s *Service) UpdateRecord(id string, items []string, actor string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status == domain.StatusArchived {
		return domain.Record{}, fmt.Errorf("archived record cannot change")
	}
	if err = record.UpdateItems(domain.NormalizeItems(items), s.Now()); err != nil {
		return domain.Record{}, err
	}
	if err = s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err = s.recordEvent(record, "updated", actor, "items changed"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) PublishRecord(id, actor string) (domain.Record, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status == domain.StatusDraft {
		if _, err = s.ReviewRecord(id, actor); err != nil {
			return domain.Record{}, err
		}
		record, err = s.ConfirmRecord(id, actor)
	} else if record.Status == domain.StatusReviewed {
		record, err = s.ConfirmRecord(id, actor)
	} else if record.Status != domain.StatusConfirmed {
		err = fmt.Errorf("record is not publishable")
	}
	return record, err
}

func (s *Service) recordEvent(record domain.Record, action, actor, details string) error {
	event := domain.NewAudit(s.nextID("audit"), record.ID, action, actor, s.Now(), details)
	return s.Store.SaveAudit(event)
}

func Digest(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func (s *Service) AuditSummary(recordID string) (report.Summary, error) {
	events, err := s.Store.ListAudits(recordID)
	if err != nil {
		return report.Summary{}, err
	}
	return report.GenerateAuditSummary(events), nil
}

func JoinItems(items []string) string { return strings.Join(domain.NormalizeItems(items), " | ") }
