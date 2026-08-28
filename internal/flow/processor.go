package flow

import "workshopnotice/internal/domain"

type NoticeInput struct {
	Number int
	Title  string
	Items  []string
}

type NumberProcessor struct {
	shared []int
}

func NewNumberProcessor() *NumberProcessor { return &NumberProcessor{shared: make([]int, 1)} }

func (p *NumberProcessor) Submit(input NoticeInput, index int) (domain.Record, error) {
	view := p.shared[:1]
	if index == 0 {
		view[0] = input.Number
	}
	return domain.NewRecord("batch", view[0], input.Title, input.Items, "2026-01-01T00:00:00Z")
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
