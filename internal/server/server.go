package server

import (
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Server struct {
	root          string
	installerPath string
	logger        *slog.Logger
}

func New(root, installerPath string, logger *slog.Logger) *Server {
	return &Server{root: root, installerPath: installerPath, logger: logger}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/healthz" {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
		return
	}
	if r.URL.Path == "/" {
		s.serveLanding(w, r)
		return
	}
	if r.URL.Path == "/watchman/install.sh" {
		s.serveInstaller(w, r)
		return
	}
	if r.URL.Path == "/skills" || r.URL.Path == "/skills/" {
		s.serveSkillsInstructions(w, r)
		return
	}
	if r.URL.Path == "/docs/watchman" || r.URL.Path == "/docs/watchman/" {
		s.serveWatchmanDocs(w, r)
		return
	}

	clean := path.Clean("/" + r.URL.Path)
	if clean == "/" || strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	filename := filepath.Join(s.root, filepath.FromSlash(strings.TrimPrefix(clean, "/")))
	isRelease := strings.Contains(clean, "/releases/")
	if isRelease {
		// Do not let an edge cache retain a not-found response while a new
		// immutable release is rolling out across service replicas.
		w.Header().Set("Cache-Control", "no-store")
	}
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}

	if isRelease {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else if strings.HasSuffix(clean, "/install.sh") || strings.HasSuffix(clean, "/setup.sh") || strings.HasSuffix(clean, "/latest") {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	} else if strings.HasSuffix(clean, ".json") {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	} else if strings.HasSuffix(clean, ".yaml") || strings.HasSuffix(clean, ".yml") {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	} else if strings.HasSuffix(clean, ".md") {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	} else if contentType := mime.TypeByExtension(filepath.Ext(filename)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	s.logger.Info("asset served", "method", r.Method, "path", clean, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}

func (s *Server) serveLanding(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.root, "index.html")
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s.logger.Info("registry landing served", "method", r.Method, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}

func (s *Server) serveWatchmanDocs(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.root, "docs", "watchman", "index.html")
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s.logger.Info("watchman docs served", "method", r.Method, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}

func (s *Server) serveInstaller(w http.ResponseWriter, r *http.Request) {
	info, err := os.Stat(s.installerPath)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s.logger.Info("installer served", "method", r.Method, "bytes", info.Size())
	http.ServeFile(w, r, s.installerPath)
}

func (s *Server) serveSkillsInstructions(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(s.root, "skills.txt")
	info, err := os.Stat(filename)
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s.logger.Info("skills instructions served", "method", r.Method, "bytes", info.Size())
	http.ServeFile(w, r, filename)
}
