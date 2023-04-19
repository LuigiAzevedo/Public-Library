package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/LuigiAzevedo/public-library-v2/config"
	r "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	handler "github.com/LuigiAzevedo/public-library-v2/internal/delivery/http"
	u "github.com/LuigiAzevedo/public-library-v2/internal/domain/usecase"
)

func main() {
	// load env configurations
	config, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load configurations")
	}

	// starts db connection
	db, err := setupDB(config.DbDriver, config.DbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to the database")
	}
	defer db.Close()

	// configurations for the logger middleware
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	httplog.Configure(httplog.Options{Concise: true, TimeFieldFormat: time.DateTime})

	router := chi.NewRouter()

	// middleware
	router.Use(httplog.RequestLogger(log.Logger))
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	// repositories DI
	userRepo := r.NewUserRepository(db)
	bookRepo := r.NewBookRepository(db)
	loanRepo := r.NewLoanRepository(db)

	// usecase DI
	userUC := u.NewUserUseCase(userRepo)
	bookUC := u.NewBookUseCase(bookRepo)
	loanUC := u.NewLoanUseCase(loanRepo, userRepo, bookRepo)

	// HTTP handlers
	handler.NewBookHandler(router, bookUC)
	handler.NewUserHandler(router, userUC)
	handler.NewLoanHandler(router, loanUC)

	server := newServer(config.ServeAddress, router)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// graceful shutdown
	waitForShutdown(server)
}

// waitForShutdown graceful shutdown
func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to gracefully shut down server")
	}
}

// setupDB initiates the database connection
func setupDB(driver, url string) (*sql.DB, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// newServer initiates a http server
func newServer(addr string, r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
