package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"workshopnotice/internal/domain"
	"workshopnotice/internal/flow"
	"workshopnotice/internal/report"
)

type Handler struct {
	Service *flow.Service
}

func New(service *flow.Service) *Handler { return &Handler{Service: service} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/health" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/records" {
		h.search(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/records" {
		h.create(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/imports" {
		h.importNotice(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/records/") {
		h.recordAction(w, r)
		return
	}
	http.NotFound(w, r)
}

type createRequest struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Items  []string `json:"items"`
	Actor  string   `json:"actor"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request createRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := h.Service.CreateRecord(request.Number, request.Title, request.Items, request.Actor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	number, _ := strconv.Atoi(r.URL.Query().Get("number"))
	records, err := h.Service.SearchRecords(domain.Query{Number: number, Text: r.URL.Query().Get("q"), Status: domain.Status(r.URL.Query().Get("status"))})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (h *Handler) importNotice(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.Service.ImportNotice(string(data), "http")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"record": result.Record, "report": report.GenerateImportReport(result)})
}

func (h *Handler) recordAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	id, action := parts[1], parts[2]
	var body struct {
		Actor string   `json:"actor"`
		Items []string `json:"items"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	var record domain.Record
	var err error
	switch action {
	case "review":
		record, err = h.Service.ReviewRecord(id, body.Actor)
	case "confirm":
		record, err = h.Service.ConfirmRecord(id, body.Actor)
	case "archive":
		record, err = h.Service.ArchiveRecord(id, body.Actor)
	case "update":
		record, err = h.Service.UpdateRecord(id, body.Items, body.Actor)
	case "publish":
		record, err = h.Service.PublishRecord(id, body.Actor)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
