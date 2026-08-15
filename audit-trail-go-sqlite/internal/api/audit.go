package api

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"example.com/template-app/internal/store"
)

// init plugs the read-only audit endpoint into the mux through the api
// socket. There is no write route: entries are appended from the code
// that performs the change, via store.Audit.Record.
func init() {
	Register(func(mux *http.ServeMux, s *Server) {
		repo := store.Audit{DB: s.DB}
		mux.HandleFunc("GET /api/audit", func(w http.ResponseWriter, r *http.Request) {
			limit, offset := 100, 0
			if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
				limit = v
			}
			if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v > 0 {
				offset = v
			}
			entries, err := repo.List(r.URL.Query().Get("object"), limit, offset)
			if err != nil {
				Fail(w, err)
				return
			}
			WriteJSON(w, http.StatusOK, entries)
		})
	})
}

// ClientIP is the caller's address, honouring the first proxy hop.
func ClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
