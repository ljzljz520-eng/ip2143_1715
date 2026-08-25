package importer

import (
	"crypto/sha256"
	"encoding/hex"
)

func Checksum(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func VerifyChecksum(data, expected string) bool { return Checksum(data) == expected }
