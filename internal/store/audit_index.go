package store

import (
	"sort"
	"strings"

	"workshopnotice/internal/domain"
)

type AuditIndex struct {
	ByAction map[string]int
	ByActor  map[string]int
	ByRecord map[string]int
}

func BuildAuditIndex(events []domain.AuditEvent) AuditIndex {
	index := AuditIndex{ByAction: map[string]int{}, ByActor: map[string]int{}, ByRecord: map[string]int{}}
	for _, event := range events {
		index.ByAction[event.Action]++
		index.ByActor[event.Actor]++
		index.ByRecord[event.RecordID]++
	}
	return index
}

func (index AuditIndex) Actions() []string {
	result := make([]string, 0, len(index.ByAction))
	for action := range index.ByAction {
		result = append(result, action)
	}
	sort.Strings(result)
	return result
}

func (index AuditIndex) FindActor(prefix string) []string {
	result := make([]string, 0)
	needle := strings.ToLower(strings.TrimSpace(prefix))
	for actor := range index.ByActor {
		if needle == "" || strings.HasPrefix(strings.ToLower(actor), needle) {
			result = append(result, actor)
		}
	}
	sort.Strings(result)
	return result
}

func (s *Store) AuditIndex(recordID string) (AuditIndex, error) {
	events, err := s.ListAudits(recordID)
	if err != nil {
		return AuditIndex{}, err
	}
	return BuildAuditIndex(events), nil
}
