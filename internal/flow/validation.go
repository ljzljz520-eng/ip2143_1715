package flow

import (
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
)

type ValidationResult struct {
	Valid  bool
	Errors []string
	Notes  []string
}

func ValidateRecord(record domain.Record) ValidationResult {
	result := ValidationResult{Valid: true, Errors: make([]string, 0), Notes: make([]string, 0)}
	if err := record.Validate(); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, err.Error())
	}
	if strings.TrimSpace(record.Title) != record.Title {
		result.Notes = append(result.Notes, "title has surrounding whitespace")
	}
	if len(record.Items) > 10 {
		result.Notes = append(result.Notes, "large checklist")
	}
	return result
}

func ValidateBatch(inputs []NoticeInput) []ValidationResult {
	result := make([]ValidationResult, 0, len(inputs))
	for _, input := range inputs {
		record, err := domain.NewRecord("validation", input.Number, input.Title, input.Items, "fixed")
		if err != nil {
			result = append(result, ValidationResult{Valid: false, Errors: []string{err.Error()}, Notes: []string{input.Title}})
			continue
		}
		result = append(result, ValidateRecord(record))
	}
	return result
}

func ValidationSummary(results []ValidationResult) string {
	valid, invalid := 0, 0
	for _, result := range results {
		if result.Valid {
			valid++
		} else {
			invalid++
		}
	}
	return fmt.Sprintf("valid=%d invalid=%d", valid, invalid)
}

func RequireValid(results []ValidationResult) error {
	for index, result := range results {
		if !result.Valid {
			return fmt.Errorf("batch item %d invalid: %s", index+1, strings.Join(result.Errors, "; "))
		}
	}
	return nil
}

func (s *Service) ValidateStored(id string, policy domain.SafetyPolicy) (ValidationResult, error) {
	record, err := s.Store.LoadRecord(id)
	if err != nil {
		return ValidationResult{}, err
	}
	result := ValidateRecord(record)
	for _, issue := range policy.Validate(record) {
		result.Valid = false
		result.Errors = append(result.Errors, issue)
	}
	return result, nil
}
