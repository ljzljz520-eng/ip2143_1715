package httpapi

import (
	"encoding/json"
	"net/http"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/report"
)

func (h *Handler) AdminSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	snapshot, err := h.Service.BuildSnapshot()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	statistics := report.CalculateStatistics(snapshot.Records, snapshot.Audits)
	writeJSON(w, http.StatusOK, map[string]any{"snapshot": snapshot, "statistics": statistics})
}

func (h *Handler) AdminExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errMethod)
		return
	}
	id := r.URL.Query().Get("id")
	data, err := h.Service.ExportRecord(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func EncodeRecord(record domain.Record) []byte {
	data, _ := json.Marshal(record)
	return data
}

var errMethod = &methodError{}

type methodError struct{}

func (*methodError) Error() string { return "method not allowed" }
