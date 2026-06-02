package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/yellowkode-academy/linkvault/internal/api"
	"github.com/yellowkode-academy/linkvault/internal/middleware"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "linkvault.db"
	}

	repo, err := storage.NewSQLiteRepository(dsn)
	if err != nil {
		slog.Error("falha ao abrir banco", "err", err)
		os.Exit(1)
	}
	defer repo.Close()

	mux := http.NewServeMux()
	api.NewLinkHandler(repo).RegisterRoutes(mux)

	handler := middleware.Logger(middleware.CORS(mux))

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	} else {
		addr = ":" + addr
	}

	slog.Info("servidor iniciado", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		slog.Error("falha ao iniciar servidor", "err", err)
		os.Exit(1)
	}
}
