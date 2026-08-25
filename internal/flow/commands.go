package flow

import (
	"fmt"
	"strings"

	"workshopnotice/internal/domain"
)

type Command struct {
	Action   string
	RecordID string
	Actor    string
	Items    []string
}

func (s *Service) Execute(command Command) (domain.Record, error) {
	actor := strings.TrimSpace(command.Actor)
	if actor == "" {
		actor = "system"
	}
	switch strings.ToLower(command.Action) {
	case "review":
		return s.ReviewRecord(command.RecordID, actor)
	case "confirm":
		return s.ConfirmRecord(command.RecordID, actor)
	case "archive":
		return s.ArchiveRecord(command.RecordID, actor)
	case "publish":
		return s.PublishRecord(command.RecordID, actor)
	case "update":
		return s.UpdateRecord(command.RecordID, command.Items, actor)
	default:
		return domain.Record{}, fmt.Errorf("unknown command %q", command.Action)
	}
}

func SupportedCommands() []string {
	return []string{"review", "confirm", "archive", "publish", "update"}
}
