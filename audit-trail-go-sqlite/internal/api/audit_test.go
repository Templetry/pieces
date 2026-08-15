package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/template-app/internal/store"
)

// newAuditServer boots the API over a throwaway in-memory database and
// hands back the handle, so the test can append entries the way
// application code does.
func newAuditServer(t *testing.T) (*httptest.Server, *sql.DB) {
	t.Helper()
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(db); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewMux(db))
	t.Cleanup(func() { srv.Close(); db.Close() })
	return srv, db
}

func TestAuditRecordsAndLists(t *testing.T) {
	srv, db := newAuditServer(t)
	repo := store.Audit{DB: db}

	if _, err := repo.Record(store.AuditEntry{
		Actor: "ada", Action: "updated", ObjectName: "note", ObjectID: "7", Detail: "title",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Record(store.AuditEntry{
		Actor: "ada", Action: "deleted", ObjectName: "invoice", ObjectID: "8",
	}); err != nil {
		t.Fatal(err)
	}

	res, body := do(t, "GET", srv.URL+"/api/audit", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("audit = %d %s", res.StatusCode, body)
	}
	var all []store.AuditEntry
	if err := json.Unmarshal(body, &all); err != nil || len(all) != 2 {
		t.Fatalf("entries = %s (%v)", body, err)
	}
	if all[0].ObjectName != "invoice" || all[0].At == "" {
		t.Errorf("newest first, with a timestamp: %+v", all[0])
	}

	_, body = do(t, "GET", srv.URL+"/api/audit?object=note", nil)
	var filtered []store.AuditEntry
	if err := json.Unmarshal(body, &filtered); err != nil || len(filtered) != 1 {
		t.Fatalf("filtered = %s (%v)", body, err)
	}
	if filtered[0].Detail != "title" || filtered[0].Actor != "ada" {
		t.Errorf("entry lost data: %+v", filtered[0])
	}
}

func TestAuditHasNoWriteRoutes(t *testing.T) {
	srv, _ := newAuditServer(t)
	for _, method := range []string{"POST", "PUT", "DELETE"} {
		res, _ := do(t, method, srv.URL+"/api/audit", map[string]string{"action": "forged"})
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s /api/audit = %d, want 405: the trail must not be writable over HTTP",
				method, res.StatusCode)
		}
	}
}
