package report

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"workshopnotice/internal/domain"
)

func RecordDigest(record domain.Record) string {
	value := record.ID + "|" + record.Title + "|" + strings.Join(record.Items, "|")
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func CompareRecords(left, right domain.Record) bool { return RecordDigest(left) == RecordDigest(right) }
