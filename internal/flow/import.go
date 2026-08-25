package flow

import (
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/importer"
	"workshopnotice/internal/report"
)

func (s *Service) ImportNotice(payload string, actor string) (report.ImportResult, error) {
	input, err := importer.Parse(payload)
	if err != nil {
		return report.ImportResult{}, err
	}
	if err = report.ValidateImport(input); err != nil {
		return report.ImportResult{}, err
	}
	record, err := s.CreateRecord(input.Number, input.Title, input.Items, actor)
	if err != nil {
		return report.ImportResult{}, err
	}
	attachment := domain.NewAttachment(s.nextID("att"), record.ID, input.Name, Digest(payload), payload, s.Now())
	if err = s.Store.SaveAttachment(attachment); err != nil {
		return report.ImportResult{}, err
	}
	if err = s.recordEvent(record, "imported", actor, fmt.Sprintf("attachment %s", input.Name)); err != nil {
		return report.ImportResult{}, err
	}
	return report.ImportResult{Record: record, Attachment: attachment, Message: strings.TrimSpace(input.Name) + " imported"}, nil
}
