// Package ipc provides a lightweight local REST API server that the daemon
// exposes on 127.0.0.1:42069. The GUI connects to this server to read usage data.
package ipc

import (
	"context"
	"encoding/json"
	"errors"
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

	// shutdownTimeout bounds how long Stop waits for in-flight requests.
	shutdownTimeout = 3 * time.Second
)

// Server is the local IPC REST server.
type Server struct {
	db        *storage.DB
	srv       *http.Server
	startedAt time.Time
	version   string
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
	mux.HandleFunc(apiPrefix+"/status", s.getOnly(s.handleStatus))
	mux.HandleFunc(apiPrefix+"/usage/today", s.getOnly(s.handleToday))
	mux.HandleFunc(apiPrefix+"/usage/week", s.getOnly(s.handleWeek))
	mux.HandleFunc(apiPrefix+"/usage/month", s.getOnly(s.handleMonth))
	mux.HandleFunc(apiPrefix+"/usage/timeline", s.getOnly(s.handleTimeline))
	mux.HandleFunc(apiPrefix+"/usage/summary", s.getOnly(s.handleSummary))
	mux.HandleFunc("/", s.handleNotFound)

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
// A graceful Stop is not reported as an error.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("ipc: listen on %s: %w", s.srv.Addr, err)
	}
	log.Printf("[ipc] listening on http://%s", s.srv.Addr)
	if err := s.srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ipc: serve: %w", err)
	}
	return nil
}

// Stop gracefully shuts down the server, waiting briefly for in-flight
// requests before closing the listener.
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Printf("[ipc] graceful shutdown failed, forcing close: %v", err)
		_ = s.srv.Close()
	}
}

// getOnly restricts a handler to GET requests, answering anything else with a
// JSON 405. (OPTIONS preflight requests are answered by corsMiddleware before
// they reach the router.)
func (s *Server) getOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeErrorStatus(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		next(w, r)
	}
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
	date, ok := queryDate(w, r)
	if !ok {
		return
	}
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
	date, ok := queryDate(w, r)
	if !ok {
		return
	}
	slots, err := s.db.GetTimeline(date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, slots)
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	date, ok := queryDate(w, r)
	if !ok {
		return
	}
	summary, err := s.db.GetDaySummary(date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, summary)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeErrorStatus(w, http.StatusNotFound, fmt.Sprintf("unknown endpoint %q", r.URL.Path))
}

// --- Helpers ---

// queryDate extracts and validates the optional ?date=YYYY-MM-DD parameter.
// It writes a JSON 400 response and returns false when the value is invalid.
func queryDate(w http.ResponseWriter, r *http.Request) (string, bool) {
	date := r.URL.Query().Get("date")
	if date == "" {
		return "", true
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		writeErrorStatus(w, http.StatusBadRequest, "invalid date parameter, want YYYY-MM-DD")
		return "", false
	}
	return date, true
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[ipc] encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, http.StatusInternalServerError, err.Error())
}

func writeErrorStatus(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
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
