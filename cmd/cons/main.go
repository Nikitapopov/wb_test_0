package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	order_repo_cache "wb_test_1/internal/cache/order"
	"wb_test_1/internal/consumer"
	"wb_test_1/internal/logger"
	"wb_test_1/internal/manager"
	"wb_test_1/internal/order_service"
	pg_order_repo "wb_test_1/internal/pg/order"
	"wb_test_1/pkg/go_cache"
	pg "wb_test_1/pkg/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	webPort    = "8080"
	pgHost     = "localhost"
	pgPort     = 5433
	pgUser     = "postgres"
	pgPassword = "postgres"
	pgDbname   = "orders"
)

func main() {
	logger := logger.NewLogger()
	pgClient := startPg()
	dbOrderRepo := pg_order_repo.NewRepo(pgClient, logger)

	goCacheClient := go_cache.NewGoCacher()
	cacheOrderRepo := order_repo_cache.NewRepo(goCacheClient, logger)

	orderService := order_service.NewService(dbOrderRepo, cacheOrderRepo, logger)
	managerHandler := manager.NewHandler(orderService)
	go consumer.Start(orderService, cacheOrderRepo, logger)

	routes1 := routes()
	managerHandler.Register(routes1)

	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%s", webPort),
		Handler: routes1,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func startPg() *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPassword, pgDbname)
	return pg.NewPg(dsn)
}

func routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	return mux
}
