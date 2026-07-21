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

func TestServesSandboxDocumentationForHumansAndAgents(t *testing.T) {
	root := t.TempDir()
	docsDir := filepath.Join(root, "docs", "sandbox")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]struct {
		body        string
		contentType string
	}{
		"index.html": {body: "<!doctype html><title>Sandbox docs</title>", contentType: "text/html; charset=utf-8"},
		"index.md":   {body: "# Sandbox documentation\n", contentType: "text/markdown; charset=utf-8"},
		"llms.txt":   {body: "# Sandbox\n", contentType: "text/plain; charset=utf-8"},
	}
	for name, file := range files {
		if err := os.WriteFile(filepath.Join(docsDir, name), []byte(file.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, route := range []string{"/docs/sandbox", "/docs/sandbox/"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, route, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("route %q status = %d", route, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != files["index.html"].contentType {
			t.Fatalf("route %q content type = %q", route, got)
		}
		if recorder.Body.String() != files["index.html"].body {
			t.Fatalf("route %q body = %q", route, recorder.Body.String())
		}
	}

	for _, name := range []string{"index.md", "llms.txt"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/docs/sandbox/"+name, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("file %q status = %d", name, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != files[name].contentType {
			t.Fatalf("file %q content type = %q", name, got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
			t.Fatalf("file %q cache control = %q", name, got)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("file %q x-content-type-options = %q", name, got)
		}
		if recorder.Body.String() != files[name].body {
			t.Fatalf("file %q body = %q", name, recorder.Body.String())
		}
	}
}

func TestServesSharedFrontendAssetsWithoutCaching(t *testing.T) {
	root := t.TempDir()
	assets := map[string]string{
		"theme.css":     "text/css; charset=utf-8",
		"tools.css":     "text/css; charset=utf-8",
		"tools.js":      "text/javascript; charset=utf-8",
		"docs/docs.css": "text/css; charset=utf-8",
		"docs/docs.js":  "text/javascript; charset=utf-8",
	}
	for name := range assets {
		filename := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, []byte("asset\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for name, contentType := range assets {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/"+name, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("asset %q status = %d", name, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != contentType {
			t.Fatalf("asset %q content type = %q", name, got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
			t.Fatalf("asset %q cache control = %q", name, got)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("asset %q x-content-type-options = %q", name, got)
		}
	}
}

func TestServesRegistryStylesheetAlias(t *testing.T) {
	root := t.TempDir()
	const stylesheet = ".registry-page { color: inherit; }\n"
	if err := os.WriteFile(filepath.Join(root, "tools.css"), []byte(stylesheet), 0o644); err != nil {
		t.Fatal(err)
	}

	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/registry.css", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/css; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("cache control = %q", got)
	}
	if recorder.Body.String() != stylesheet {
		t.Fatalf("body = %q", recorder.Body.String())
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

func TestServesInstallableSkillAssets(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "public")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(base, "skills", "objects-storage", "agents")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]struct {
		body        string
		contentType string
	}{
		filepath.Join(base, "skills", "objects-storage", "SKILL.md"): {
			body:        "---\nname: objects-storage\n---\n",
			contentType: "text/markdown; charset=utf-8",
		},
		filepath.Join(skillDir, "openai.yaml"): {
			body:        "interface:\n  display_name: Objects Storage\n",
			contentType: "application/yaml; charset=utf-8",
		},
	}
	for filename, file := range files {
		if err := os.WriteFile(filename, []byte(file.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for route, expected := range map[string]struct {
		body        string
		contentType string
	}{
		"/skills/objects-storage/SKILL.md": {
			body:        files[filepath.Join(base, "skills", "objects-storage", "SKILL.md")].body,
			contentType: files[filepath.Join(base, "skills", "objects-storage", "SKILL.md")].contentType,
		},
		"/skills/objects-storage/agents/openai.yaml": {
			body:        files[filepath.Join(skillDir, "openai.yaml")].body,
			contentType: files[filepath.Join(skillDir, "openai.yaml")].contentType,
		},
	} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, route, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("route %q status = %d", route, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != expected.contentType {
			t.Fatalf("route %q content type = %q", route, got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
			t.Fatalf("route %q cache control = %q", route, got)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("route %q x-content-type-options = %q", route, got)
		}
		if got := recorder.Body.String(); got != expected.body {
			t.Fatalf("route %q body = %q", route, got)
		}
	}

	for _, route := range []string{"/skills/objects-storage/", "/skills/objects-storage/../SKILL.md"} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, route, nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("route %q status = %d", route, recorder.Code)
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

func TestServesMCPClientTemplatesWithoutCaching(t *testing.T) {
	root := t.TempDir()
	clientsDir := filepath.Join(root, "sandbox", "clients")
	if err := os.MkdirAll(clientsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"opencode.json": "application/json; charset=utf-8",
		"goose.yaml":    "application/yaml; charset=utf-8",
		"index.md":      "text/markdown; charset=utf-8",
	}
	for name := range files {
		if err := os.WriteFile(filepath.Join(clientsDir, name), []byte("template\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	server := New(root, filepath.Join(t.TempDir(), "install.sh"), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for name, contentType := range files {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sandbox/clients/"+name, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("template %q status = %d", name, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != contentType {
			t.Fatalf("template %q content type = %q", name, got)
		}
		if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
			t.Fatalf("template %q cache control = %q", name, got)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("template %q x-content-type-options = %q", name, got)
		}
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
