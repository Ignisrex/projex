package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"project-tracker/server/internal/config"
	httplayer "project-tracker/server/internal/http"
	"project-tracker/server/internal/repositories"
	"project-tracker/server/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connection established")

	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	projectService := services.NewProjectService(projectRepo, taskRepo)
	taskService := services.NewTaskService(taskRepo)

	router := httplayer.NewRouter(cfg, projectService, taskService)

	slog.Info("server listening", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
