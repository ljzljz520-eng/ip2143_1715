package flow

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type BatchResult struct {
	Created []domain.Record
	Errors  []string
	Total   int
}

func (s *Service) CreateBatch(inputs []NoticeInput, actor string) BatchResult {
	result := BatchResult{Created: make([]domain.Record, 0), Errors: make([]string, 0), Total: len(inputs)}
	for index, input := range inputs {
		record, err := s.CreateRecord(input.Number, input.Title, input.Items, actor)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("item %d: %v", index+1, err))
			continue
		}
		result.Created = append(result.Created, record)
	}
	return result
}

func (s *Service) ReviewPending(actor string) (int, error) {
	records, err := s.SearchRecords(domain.Query{Status: domain.StatusDraft})
	if err != nil {
		return 0, err
	}
	count := 0
	for _, record := range records {
		if _, err = s.ReviewRecord(record.ID, actor); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) ConfirmReviewed(actor string) (int, error) {
	records, err := s.SearchRecords(domain.Query{Status: domain.StatusReviewed})
	if err != nil {
		return 0, err
	}
	count := 0
	for _, record := range records {
		if _, err = s.ConfirmRecord(record.ID, actor); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) ArchiveConfirmed(actor string) (int, error) {
	records, err := s.SearchRecords(domain.Query{Status: domain.StatusConfirmed})
	if err != nil {
		return 0, err
	}
	count := 0
	for _, record := range records {
		if _, err = s.ArchiveRecord(record.ID, actor); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func SortByRisk(records []domain.Record, policy domain.SafetyPolicy) {
	sort.SliceStable(records, func(i, j int) bool {
		left, right := policy.Score(records[i]), policy.Score(records[j])
		if left == right {
			return records[i].Number < records[j].Number
		}
		return left < right
	})
}

func SummarizeBatch(result BatchResult) string {
	return fmt.Sprintf("total=%d created=%d errors=%d", result.Total, len(result.Created), len(result.Errors))
}
