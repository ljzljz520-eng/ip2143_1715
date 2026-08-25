package httpapi

import (
	"net/http"
	"strings"
)

func allowedActor(actor string) bool { return strings.TrimSpace(actor) != "" }

func actorOrDefault(actor string) string {
	if !allowedActor(actor) {
		return "anonymous"
	}
	return strings.TrimSpace(actor)
}

func requireJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
