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

func TestRejectsDirectoriesAndWriteMethods(t *testing.T) {
	server := New(t.TempDir(), filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, request := range []*http.Request{
		httptest.NewRequest(http.MethodGet, "/", nil),
		httptest.NewRequest(http.MethodPost, "/watchman/install.sh", nil),
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
