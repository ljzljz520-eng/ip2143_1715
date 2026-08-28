package domain

import (
	"fmt"
	"strings"
)

func StatusLabel(status Status) string {
	switch status {
	case StatusDraft:
		return "登记中"
	case StatusReviewed:
		return "已审核"
	case StatusConfirmed:
		return "已发布"
	case StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}

func (r Record) DisplayLines() []string {
	lines := []string{fmt.Sprintf("编号 %d", r.Number), r.Title, "状态: " + StatusLabel(r.Status)}
	for index, item := range r.Items {
		lines = append(lines, fmt.Sprintf("%d. %s", index+1, item))
	}
	return lines
}

func ParseStatus(value string) (Status, error) {
	status := Status(strings.ToLower(strings.TrimSpace(value)))
	switch status {
	case StatusDraft, StatusReviewed, StatusConfirmed, StatusArchived:
		return status, nil
	default:
		return "", fmt.Errorf("unknown status %q", value)
	}
}

func CanEdit(status Status) bool { return status != StatusArchived }

func CanArchive(status Status) bool { return status == StatusConfirmed }
