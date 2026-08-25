package domain

import "sort"

type History struct {
	RecordID string
	Events   []AuditEvent
}

func NewHistory(recordID string, events []AuditEvent) History {
	copyEvents := append([]AuditEvent(nil), events...)
	sort.SliceStable(copyEvents, func(i, j int) bool { return copyEvents[i].At < copyEvents[j].At })
	return History{RecordID: recordID, Events: copyEvents}
}

func (h History) HasAction(action string) bool {
	for _, event := range h.Events {
		if event.Action == action {
			return true
		}
	}
	return false
}

func (h History) Actors() []string {
	actors := make([]string, 0)
	seen := map[string]bool{}
	for _, event := range h.Events {
		if !seen[event.Actor] {
			seen[event.Actor] = true
			actors = append(actors, event.Actor)
		}
	}
	return actors
}

func (h History) LastAction() string {
	if len(h.Events) == 0 {
		return ""
	}
	return h.Events[len(h.Events)-1].Action
}
