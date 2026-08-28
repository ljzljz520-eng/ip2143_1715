package domain

import (
	"fmt"
	"strings"
)

type SafetyPolicy struct {
	RequiredWords []string
	MinimumItems  int
	AllowDrafts   bool
}

func DefaultPolicy() SafetyPolicy {
	return SafetyPolicy{RequiredWords: []string{"安全", "检查", "确认"}, MinimumItems: 1, AllowDrafts: true}
}

func (p SafetyPolicy) Validate(record Record) []string {
	issues := make([]string, 0)
	if record.Number <= 0 {
		issues = append(issues, "编号必须为正数")
	}
	if strings.TrimSpace(record.Title) == "" {
		issues = append(issues, "标题不能为空")
	}
	if len(record.Items) < p.MinimumItems {
		issues = append(issues, fmt.Sprintf("至少需要%d项", p.MinimumItems))
	}
	if !p.AllowDrafts && record.Status == StatusDraft {
		issues = append(issues, "当前策略不允许草稿")
	}
	return issues
}

func (p SafetyPolicy) Score(record Record) int {
	score := 0
	text := strings.ToLower(record.Title + " " + strings.Join(record.Items, " "))
	for _, word := range p.RequiredWords {
		if strings.Contains(text, strings.ToLower(word)) {
			score++
		}
	}
	if record.Status == StatusConfirmed {
		score++
	}
	return score
}

func RiskBand(score int) string {
	switch {
	case score <= 1:
		return "high"
	case score == 2:
		return "medium"
	default:
		return "low"
	}
}

func ReviewMessage(record Record, policy SafetyPolicy) string {
	issues := policy.Validate(record)
	if len(issues) > 0 {
		return strings.Join(issues, "; ")
	}
	return fmt.Sprintf("%s risk=%s", record.Summary(), RiskBand(policy.Score(record)))
}
