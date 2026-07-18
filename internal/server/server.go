package server

import (
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	root            string
	installerSource string
	logger          *slog.Logger
	client          *http.Client
}

func New(root, installerSource string, logger *slog.Logger) *Server {
	return &Server{root: root, installerSource: installerSource, logger: logger, client: &http.Client{Timeout: 10 * time.Second}}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/watchman/install.sh" {
		s.proxyInstaller(w, r)
		return
	}
	if r.URL.Path == "/skills" || r.URL.Path == "/skills/" {
		s.serveSkillsPage(w, r)
		return
	}

	clean := path.Clean("/" + r.URL.Path)
	if clean == "/" || strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	filename := filepath.Join(s.root, filepath.FromSlash(strings.TrimPrefix(clean, "/")))
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}

	if strings.Contains(clean, "/releases/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else if contentType := mime.TypeByExtension(filepath.Ext(filename)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	s.logger.Info("asset served", "method", r.Method, "path", clean, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}

func (s *Server) serveSkillsPage(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.root, "skills", "index.html")
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; img-src data:; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
	w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s.logger.Info("skills page served", "method", r.Method, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}

func (s *Server) proxyInstaller(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequestWithContext(r.Context(), r.Method, s.installerSource, nil)
	if err != nil {
		http.Error(w, "installer source is invalid", http.StatusBadGateway)
		return
	}
	response, err := s.client.Do(request)
	if err != nil {
		s.logger.Error("installer proxy failed", "error", err)
		http.Error(w, "installer source is unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		s.logger.Error("installer source returned an error", "status", response.StatusCode)
		http.Error(w, "installer source is unavailable", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(response.StatusCode)
	if r.Method != http.MethodHead {
		_, _ = io.Copy(w, response.Body)
	}
}
