package report

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type Alert struct {
	RecordID string
	Level    string
	Message  string
}

func BuildAlerts(records []domain.Record, policy domain.SafetyPolicy) []Alert {
	alerts := make([]Alert, 0)
	for _, record := range records {
		issues := policy.Validate(record)
		if len(issues) > 0 {
			alerts = append(alerts, Alert{RecordID: record.ID, Level: "error", Message: fmt.Sprintf("%s: %v", record.Title, issues)})
			continue
		}
		score := policy.Score(record)
		if domain.RiskBand(score) == "high" {
			alerts = append(alerts, Alert{RecordID: record.ID, Level: "warning", Message: fmt.Sprintf("%s requires attention", record.Title)})
		}
	}
	sort.SliceStable(alerts, func(i, j int) bool { return alerts[i].RecordID < alerts[j].RecordID })
	return alerts
}

func AlertCounts(alerts []Alert) map[string]int {
	counts := map[string]int{}
	for _, alert := range alerts {
		counts[alert.Level]++
	}
	return counts
}

func Critical(alerts []Alert) []Alert {
	result := make([]Alert, 0)
	for _, alert := range alerts {
		if alert.Level == "error" {
			result = append(result, alert)
		}
	}
	return result
}
