package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusReviewed  Status = "reviewed"
	StatusConfirmed Status = "confirmed"
	StatusArchived  Status = "archived"
)

type Record struct {
	ID        string   `json:"id"`
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	Items     []string `json:"items"`
	Status    Status   `json:"status"`
	Version   int      `json:"version"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	At       string `json:"at"`
	Details  string `json:"details"`
}

type Workflow struct {
	ID          string `json:"id"`
	RecordID    string `json:"record_id"`
	Kind        string `json:"kind"`
	Stage       string `json:"stage"`
	State       string `json:"state"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

var (
	ErrInvalidNumber = errors.New("notice number must be positive")
	ErrInvalidTitle  = errors.New("notice title is required")
	ErrNoItems       = errors.New("notice requires at least one item")
	ErrTransition    = errors.New("invalid notice status transition")
)

func NewRecord(id string, number int, title string, items []string, stamp string) (Record, error) {
	if number <= 0 {
		return Record{}, ErrInvalidNumber
	}
	if strings.TrimSpace(title) == "" {
		return Record{}, ErrInvalidTitle
	}
	clean := make([]string, 0, len(items))
	for _, item := range items {
		if text := strings.TrimSpace(item); text != "" {
			clean = append(clean, text)
		}
	}
	if len(clean) == 0 {
		return Record{}, ErrNoItems
	}
	return Record{ID: id, Number: number, Title: strings.TrimSpace(title), Items: clean, Status: StatusDraft, Version: 1, CreatedAt: stamp, UpdatedAt: stamp}, nil
}

func (r Record) Validate() error {
	if r.Number <= 0 {
		return ErrInvalidNumber
	}
	if strings.TrimSpace(r.Title) == "" {
		return ErrInvalidTitle
	}
	if len(r.Items) == 0 {
		return ErrNoItems
	}
	if r.Version <= 0 {
		return fmt.Errorf("version must be positive")
	}
	return nil
}

func (r *Record) Transition(next Status) error {
	valid := (r.Status == StatusDraft && next == StatusReviewed) ||
		(r.Status == StatusReviewed && next == StatusConfirmed) ||
		(r.Status == StatusConfirmed && next == StatusArchived)
	if !valid {
		return ErrTransition
	}
	r.Status = next
	r.Version++
	return nil
}

func (r *Record) UpdateItems(items []string, stamp string) error {
	if len(items) == 0 {
		return ErrNoItems
	}
	clean := make([]string, 0, len(items))
	for _, item := range items {
		if text := strings.TrimSpace(item); text != "" {
			clean = append(clean, text)
		}
	}
	if len(clean) == 0 {
		return ErrNoItems
	}
	r.Items = clean
	r.UpdatedAt = stamp
	r.Version++
	return nil
}

func (r Record) Summary() string {
	return fmt.Sprintf("%d: %s (%s, v%d)", r.Number, r.Title, r.Status, r.Version)
}

func NewAudit(id, recordID, action, actor, at, details string) AuditEvent {
	return AuditEvent{ID: id, RecordID: recordID, Action: action, Actor: actor, At: at, Details: details}
}

func NewWorkflow(id, recordID, kind, stage, stamp string) Workflow {
	return Workflow{ID: id, RecordID: recordID, Kind: kind, Stage: stage, State: "active", StartedAt: stamp}
}

func CompleteWorkflow(w *Workflow, stamp string) {
	w.State = "completed"
	w.CompletedAt = stamp
}

func NewAttachment(id, recordID, name, digest, content, stamp string) Attachment {
	return Attachment{ID: id, RecordID: recordID, Name: name, Digest: digest, Content: content, CreatedAt: stamp}
}
