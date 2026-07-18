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

func TestServesSkillsPageWithBrowserSafetyHeaders(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "skills")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "index.html"), []byte("<!doctype html><title>Skills</title>"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, path := range []string{"/skills", "/skills/"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("path %q status = %d", path, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
			t.Fatalf("content type = %q", got)
		}
		if got := recorder.Header().Get("Content-Security-Policy"); got == "" {
			t.Fatal("missing content security policy")
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
