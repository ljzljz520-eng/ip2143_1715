package domain

import "strings"

func MergeRecord(base, incoming Record) Record {
	merged := base
	if strings.TrimSpace(incoming.Title) != "" {
		merged.Title = strings.TrimSpace(incoming.Title)
	}
	if len(incoming.Items) > 0 {
		merged.Items = NormalizeItems(incoming.Items)
	}
	if incoming.Number > 0 {
		merged.Number = incoming.Number
	}
	if incoming.Status != "" {
		merged.Status = incoming.Status
	}
	if incoming.Version > merged.Version {
		merged.Version = incoming.Version
	}
	if incoming.UpdatedAt != "" {
		merged.UpdatedAt = incoming.UpdatedAt
	}
	return merged
}

func CloneRecord(record Record) Record {
	copyRecord := record
	copyRecord.Items = append([]string(nil), record.Items...)
	return copyRecord
}

func CloneRecords(records []Record) []Record {
	result := make([]Record, 0, len(records))
	for _, record := range records {
		result = append(result, CloneRecord(record))
	}
	return result
}
