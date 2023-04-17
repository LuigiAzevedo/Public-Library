package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/LuigiAzevedo/public-library-v2/config"
	r "github.com/LuigiAzevedo/public-library-v2/internal/database/repository"
	handler "github.com/LuigiAzevedo/public-library-v2/internal/delivery/http"
	u "github.com/LuigiAzevedo/public-library-v2/internal/domain/usecase"
)

func main() {
	config, err := config.LoadAppConfig(".")
	if err != nil {
		log.Fatalf("unable to load configurations: %s", err)
	}

	db, err := setupDB(config.DbDriver, config.DbURL)
	if err != nil {
		log.Fatalf("unable to connect to the database: %s", err)
	}

	defer db.Close()

	router := chi.NewRouter()

	userRepo := r.NewUserRepository(db)
	userUC := u.NewUserService(userRepo)

	bookRepo := r.NewBookRepository(db)
	bookUC := u.NewBookService(bookRepo)

	loanRepo := r.NewLoanRepository(db)
	loanUC := u.NewLoanService(loanRepo, userRepo, bookRepo)

	handler.NewBookHandler(router, bookUC)
	handler.NewUserHandler(router, userUC)
	handler.NewLoanHandler(router, loanUC)

	server := newServer(config.ServeAddress, router)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err)
		}
	}()
	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to gracefully shut down server: %s", err)
	}
}

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

func newServer(addr string, r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
