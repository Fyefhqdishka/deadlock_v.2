package app

import (
	"context"
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/config"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/handlers"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/service"
	"github.com/Fyefhqdishka/deadlock_v.2/internal/storage/postgres"
	"github.com/Fyefhqdishka/deadlock_v.2/pkg/routes"
	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type App struct {
	db     *pgxpool.Pool
	server *http.Server
}

func (s *App) Run() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %v", err)
	}
	return nil
}

func (s *App) Stop() error {
	s.db.Close()

	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// New creates new instance of application, sets the dependencies and applies migrations
func New(cfg *config.Config) (*App, error) {
	log := initLogging()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := postgres.ConnectDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	initMigrations(dsn)

	storage := postgres.NewStorage(db, log)

	service := service.NewService(storage, log)

	handlers := handlers.NewHandlers(service)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	r := mux.NewRouter()
	routes.RegisterRoutes(r, handlers)
	handler := c.Handler(r)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Info("server starting in port", cfg.Server.Port)

	app := &App{
		db: db,
		server: &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 20,
			Handler:        handler,
			WriteTimeout:   cfg.Server.Timeout,
			ReadTimeout:    cfg.Server.Timeout,
			IdleTimeout:    cfg.Server.IdleTimeout,
		},
	}

	return app, nil
}

func initLogging() *slog.Logger {
	logFileName := "logs/app-" + time.Now().Format("2006-01-02") + ".log"
	logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Не удалось открыть файл для логов: ", err)
	}

	handler := slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(handler)
}

func initMigrations(dsn string) {
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Ошибка при создании объекта миграции: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}

	fmt.Println("Миграции успешно применены!")
}
