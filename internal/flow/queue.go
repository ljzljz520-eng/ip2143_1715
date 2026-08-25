package flow

import (
	"fmt"
	"sort"

	"workshopnotice/internal/domain"
)

type QueueItem struct {
	RecordID string
	Priority int
	Reason   string
}

type ReviewQueue struct {
	Items []QueueItem
}

func BuildReviewQueue(records []domain.Record, policy domain.SafetyPolicy) ReviewQueue {
	items := make([]QueueItem, 0)
	for _, record := range records {
		if record.Status == domain.StatusDraft || record.Status == domain.StatusReviewed {
			items = append(items, QueueItem{RecordID: record.ID, Priority: policy.Score(record), Reason: domain.ReviewMessage(record, policy)})
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Priority < items[j].Priority })
	return ReviewQueue{Items: items}
}

func (q ReviewQueue) Next() (QueueItem, bool) {
	if len(q.Items) == 0 {
		return QueueItem{}, false
	}
	return q.Items[0], true
}

func (q *ReviewQueue) Remove(id string) {
	for index, item := range q.Items {
		if item.RecordID == id {
			q.Items = append(q.Items[:index], q.Items[index+1:]...)
			return
		}
	}
}

func (q ReviewQueue) String() string {
	if len(q.Items) == 0 {
		return "review queue empty"
	}
	return fmt.Sprintf("review queue: %d items", len(q.Items))
}
