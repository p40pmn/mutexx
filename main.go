package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	ctx := context.Background()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	e := echo.New()

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(fmt.Sprintf(":%s", getEnv("PORT", "3003")))
	}()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err)
		}
		log.Println("server shutdown gracefully")

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(ctx, time.Second*15)
		defer cancel()

		log.Println("shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown server: %s", err)
		}
		log.Println("shutdown server gracefully")
	}

}
