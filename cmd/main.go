package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/codebyaadi/rss-feed-agg/config"
	"github.com/codebyaadi/rss-feed-agg/internal/database"
	"github.com/codebyaadi/rss-feed-agg/internal/handlers"
	"github.com/codebyaadi/rss-feed-agg/internal/redis"
	"github.com/codebyaadi/rss-feed-agg/internal/utils"
)

func main() {
	log.Print("starting server...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	dbUrl := os.Getenv("POSTGRES_URL")
	if dbUrl == "" {
		log.Fatal("POSTGRES_URL must be set")
	}

	if err := redis.InitRedis(); err != nil {
		log.Fatalf("can't connect to Redis: %v", err)
	}
	defer redis.CloseRedis()

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("can't connect to postgres database %v", err)
	}
	defer conn.Close()

	db := database.New(conn)
	apiCfg := &config.ApiConfig{
		DB: db,
	}

	go utils.RSSFeedScrapper(db, 10, time.Minute)

	handler := &handlers.Handler{ApiConfig: apiCfg}

	if err := conn.Ping(); err != nil {
		log.Fatalf("can't reach postgres database %v", err)
	}

	log.Println("successfully connected to database")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlerHealth)
	mux.HandleFunc("GET /error", handlerErr)
	mux.HandleFunc(("POST /users/create"), handler.CreateUser)
	mux.HandleFunc("POST /users/login", handler.LoginUser)
	mux.HandleFunc(("GET /users"), handler.AuthMiddleware(handler.GetUserByAPIKey))
	mux.HandleFunc(("POST /feeds/create"), handler.AuthMiddleware(handler.CreateFeed))
	mux.HandleFunc(("GET /feeds"), handler.GetAllFeeds)
	mux.HandleFunc(("POST /feeds/follow"), handler.AuthMiddleware(handler.CreateFeedFollow))
	mux.HandleFunc(("GET /feeds/follow"), handler.AuthMiddleware(handler.GetAllFeedFollows))
	mux.HandleFunc(("DELETE /feeds/follow/{feedFollowID}"), handler.AuthMiddleware(handler.DeleteFeedFollow))
	mux.HandleFunc(("GET /posts"), handler.AuthMiddleware(handler.GetPostsForUser))

	addr := ":" + port
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", port, err)
		}
	}()

	<-shutdownCh
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("server stopped")
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
}
