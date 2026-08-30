// Package ipc provides a lightweight local REST API server that the daemon
// exposes on 127.0.0.1:42069. The GUI connects to this server to read usage data.
package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/mahmud-r-farhan/chronotrace/pkg/storage"
)

const (
	DefaultAddr = "127.0.0.1:42069"
	apiPrefix   = "/api/v1"
)

// Server is the local IPC REST server.
type Server struct {
	db       *storage.DB
	srv      *http.Server
	startedAt time.Time
	version  string
}

// StatusResponse is returned by GET /api/v1/status.
type StatusResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	UptimeSec int64  `json:"uptime_seconds"`
	Addr      string `json:"addr"`
}

// New creates a new IPC server bound to addr (e.g. "127.0.0.1:42069").
func New(addr string, db *storage.DB, version string) *Server {
	if addr == "" {
		addr = DefaultAddr
	}
	s := &Server{
		db:        db,
		startedAt: time.Now(),
		version:   version,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiPrefix+"/status", s.handleStatus)
	mux.HandleFunc(apiPrefix+"/usage/today", s.handleToday)
	mux.HandleFunc(apiPrefix+"/usage/week", s.handleWeek)
	mux.HandleFunc(apiPrefix+"/usage/month", s.handleMonth)
	mux.HandleFunc(apiPrefix+"/usage/timeline", s.handleTimeline)
	mux.HandleFunc(apiPrefix+"/usage/summary", s.handleSummary)

	s.srv = &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s
}

// Start begins serving. It blocks until the server is stopped.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("ipc: listen on %s: %w", s.srv.Addr, err)
	}
	log.Printf("[ipc] listening on http://%s", s.srv.Addr)
	return s.srv.Serve(ln)
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() {
	_ = s.srv.Close()
}

// --- Handlers ---

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, StatusResponse{
		Status:    "ok",
		Version:   s.version,
		UptimeSec: int64(time.Since(s.startedAt).Seconds()),
		Addr:      s.srv.Addr,
	})
}

func (s *Server) handleToday(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	apps, err := s.db.GetUsageByDay(date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, apps)
}

func (s *Server) handleWeek(w http.ResponseWriter, r *http.Request) {
	apps, err := s.db.GetUsageByWeek()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, apps)
}

func (s *Server) handleMonth(w http.ResponseWriter, r *http.Request) {
	apps, err := s.db.GetUsageByMonth()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, apps)
}

func (s *Server) handleTimeline(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	slots, err := s.db.GetTimeline(date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, slots)
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	summary, err := s.db.GetDaySummary(date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, summary)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[ipc] encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// corsMiddleware adds CORS headers so the Wails WebView can call the daemon.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
