package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServesRegistryLandingPage(t *testing.T) {
	root := t.TempDir()
	const landing = "<!doctype html><title>tools</title><main>agent skills and MCP servers</main>"
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte(landing), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(method, "/", nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("method %s status = %d", method, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
			t.Fatalf("content type = %q", got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
			t.Fatalf("cache control = %q", got)
		}
		if method == http.MethodGet && recorder.Body.String() != landing {
			t.Fatalf("body = %q", recorder.Body.String())
		}
	}
}

func TestServesWatchmanDocumentation(t *testing.T) {
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs", "watchman")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	const rendered = "<!doctype html><title>Watchman docs</title>"
	const markdown = "# Watchman documentation\n"
	if err := os.WriteFile(filepath.Join(docsDir, "index.html"), []byte(rendered), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "reference.md"), []byte(markdown), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))

	for _, route := range []string{"/docs/watchman", "/docs/watchman/"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, route, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("route %q status = %d", route, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
			t.Fatalf("route %q content type = %q", route, got)
		}
		if recorder.Body.String() != rendered {
			t.Fatalf("route %q body = %q", route, recorder.Body.String())
		}
	}

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/docs/watchman/reference.md", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("markdown status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/markdown; charset=utf-8" {
		t.Fatalf("markdown content type = %q", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("markdown cache control = %q", got)
	}
	if recorder.Body.String() != markdown {
		t.Fatalf("markdown body = %q", recorder.Body.String())
	}
}

func TestServesPackagedInstallerWithSafeHeaders(t *testing.T) {
	installer := filepath.Join(t.TempDir(), "install.sh")
	if err := os.WriteFile(installer, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(t.TempDir(), installer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/watchman/install.sh", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("cache control = %q", got)
	}
}

func TestServesPlainTextSkillsInstructions(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "skills.txt"), []byte("npx skills add owner/repo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, path := range []string{"/skills", "/skills/"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("path %q status = %d", path, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
			t.Fatalf("content type = %q", got)
		}
		if body := recorder.Body.String(); body != "npx skills add owner/repo\n" {
			t.Fatalf("body = %q", body)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("x-content-type-options = %q", got)
		}
	}
}

func TestServesSetupScriptsWithoutCaching(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sandbox"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sandbox", "setup.sh"), []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sandbox/setup.sh", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("cache control = %q", got)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("x-content-type-options = %q", got)
	}
}

func TestRejectsDirectoriesAndWriteMethods(t *testing.T) {
	server := New(t.TempDir(), filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/sandbox/", nil),
		httptest.NewRequest(http.MethodPost, "/watchman/install.sh", nil),
		httptest.NewRequest(http.MethodPost, "/healthz", nil),
	} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusNotFound && recorder.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d", recorder.Code)
		}
	}
}

func TestMissingReleaseAssetIsNotCached(t *testing.T) {
	server := New(t.TempDir(), filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/watchman/releases/v9.9.9/missing.tar.gz", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache control = %q", got)
	}
}
