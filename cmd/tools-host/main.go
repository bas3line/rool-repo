package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bas3line/tools-host/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	root := os.Getenv("TOOLS_ROOT")
	if root == "" {
		root = "./public"
	}
	const installerPath = "./install.sh"

	handler := server.New(root, installerPath, slog.Default())
	slog.Info("tools host listening", "address", "0.0.0.0:"+port, "root", root, "installer_path", installerPath)
	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		slog.Error("tools host stopped", "error", err)
		os.Exit(1)
	}
}
