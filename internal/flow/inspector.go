package flow

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type Inspection struct {
	RecordID string
	Score    int
	Band     string
	Issues   []string
	Ready    bool
}

func Inspect(record domain.Record, policy domain.SafetyPolicy) Inspection {
	issues := policy.Validate(record)
	score := policy.Score(record)
	return Inspection{RecordID: record.ID, Score: score, Band: domain.RiskBand(score), Issues: issues, Ready: len(issues) == 0}
}

func InspectAll(records []domain.Record, policy domain.SafetyPolicy) []Inspection {
	result := make([]Inspection, 0, len(records))
	for _, record := range records {
		result = append(result, Inspect(record, policy))
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].Score < result[j].Score })
	return result
}

func InspectionMessage(result Inspection) string {
	if result.Ready {
		return fmt.Sprintf("%s ready at %s risk", result.RecordID, result.Band)
	}
	return fmt.Sprintf("%s blocked: %v", result.RecordID, result.Issues)
}

func (s *Service) InspectRecord(id string, policy domain.SafetyPolicy) (Inspection, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return Inspection{}, err
	}
	return Inspect(record, policy), nil
}
