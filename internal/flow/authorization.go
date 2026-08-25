package flow

import (
	"fmt"

	"workshopnotice/internal/domain"
)

func Authorize(role, action string) error {
	parsed, err := domain.ParseRole(role)
	if err != nil {
		return err
	}
	if !domain.ActionAllowed(parsed, action) {
		return fmt.Errorf("%s", domain.PermissionMessage(parsed, action))
	}
	return nil
}

func (s *Service) ExecuteAs(role string, command Command) (domain.Record, error) {
	if err := Authorize(role, command.Action); err != nil {
		return domain.Record{}, err
	}
	return s.Execute(command)
}

func (s *Service) CreateAs(role string, number int, title string, items []string) (domain.Record, error) {
	if err := Authorize(role, "create"); err != nil {
		return domain.Record{}, err
	}
	return s.CreateRecord(number, title, items, role)
}

func (s *Service) ArchiveAs(role, id string) (domain.Record, error) {
	if err := Authorize(role, "archive"); err != nil {
		return domain.Record{}, err
	}
	return s.ArchiveRecord(id, role)
}
