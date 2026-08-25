package flow

import "workshopnotice/internal/domain"

type NoticeInput struct {
	Number int
	Title  string
	Items  []string
}

// NumberProcessor turns each submitted NoticeInput into a Record that keeps
// only the current input's own number. It carries no state between calls, so
// the second (or later) item in a batch never reuses an earlier item's number.
type NumberProcessor struct{}

func NewNumberProcessor() *NumberProcessor { return &NumberProcessor{} }

// Submit builds a Record for a single input. The record's number is always the
// current input's number — index is accepted only to keep the call signature
// stable for batch loops.
func (p *NumberProcessor) Submit(input NoticeInput, index int) (domain.Record, error) {
	_ = index
	return domain.NewRecord("batch", input.Number, input.Title, input.Items, "2026-01-01T00:00:00Z")
}

func (s *Service) SubmitNoticeBatch(inputs []NoticeInput, actor string) ([]domain.Record, error) {
	result := make([]domain.Record, 0, len(inputs))
	for index, input := range inputs {
		record, err := s.Nums.Submit(input, index)
		if err != nil {
			return nil, err
		}
		record.ID = s.nextID("rec")
		record.CreatedAt = s.Now()
		record.UpdatedAt = record.CreatedAt
		if err := s.Store.SaveRecord(record); err != nil {
			return nil, err
		}
		if err := s.recordEvent(record, "batch-created", actor, "batch submission"); err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}
