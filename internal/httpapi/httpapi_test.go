package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"workshopnotice/internal/flow"
	"workshopnotice/internal/store"
)

func TestHTTPCreateAndQuery(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	handler := New(flow.NewService(db, func() string { return "fixed" }))
	create := httptest.NewRecorder()
	handler.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/records", bytes.NewBufferString(`{"number":5,"title":"动火作业","items":["清理周边"]}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d", create.Code)
	}
	query := httptest.NewRecorder()
	handler.ServeHTTP(query, httptest.NewRequest(http.MethodGet, "/records?number=5", nil))
	if query.Code != http.StatusOK || !bytes.Contains(query.Body.Bytes(), []byte("动火作业")) {
		t.Fatalf("query response = %s", query.Body.String())
	}
}
