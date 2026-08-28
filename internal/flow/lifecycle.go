package flow

import (
	"fmt"

	"workshopnotice/internal/domain"
)

type LifecycleResult struct {
	Record  domain.Record
	Actions []string
}

func (s *Service) CompleteLifecycle(id, actor string) (LifecycleResult, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return LifecycleResult{}, err
	}
	actions := make([]string, 0, 3)
	if record.Status == domain.StatusDraft {
		record, err = s.ReviewRecord(id, actor)
		actions = append(actions, "review")
	}
	if err != nil {
		return LifecycleResult{}, err
	}
	if record.Status == domain.StatusReviewed {
		record, err = s.ConfirmRecord(id, actor)
		actions = append(actions, "confirm")
	}
	if err != nil {
		return LifecycleResult{}, err
	}
	if record.Status == domain.StatusConfirmed {
		record, err = s.ArchiveRecord(id, actor)
		actions = append(actions, "archive")
	}
	if err != nil {
		return LifecycleResult{}, err
	}
	if record.Status != domain.StatusArchived {
		return LifecycleResult{}, fmt.Errorf("lifecycle incomplete")
	}
	return LifecycleResult{Record: record, Actions: actions}, nil
}
