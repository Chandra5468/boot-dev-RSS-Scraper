package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*
	Purpose of project :

		Aggregate data from RSS Feeds

		RSS Feeds :
			RSS(Really Simple Syndication) Feeds are an easy way to stay up to date with your favorite websites, such as blogs or online magazines.
			If a site offers an RSS feed, you get notified whenever a post goes up, and then you can read a summary or the whole post

		RSS Scraper :
			An RSS scraper is a tool or script that automatically extracts data from RSS feeds,
			which are web feeds that deliver regularly updated content in a standardized, computer-readable form

		Needed things :

			a) Postgresql

			b) VSCode

			c) Golang

			d) chi router : light weight router
*/

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type dbStruct struct {
	db *sql.DB
}

func main() {
	fmt.Println("Namaste")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("port not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DBUrl not found in the environment")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("postgres connection failed")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("error pinging db ")
	}

	pdb := &dbStruct{
		db: db,
	}

	router := chi.NewRouter()

	router.Use(cors)

	// Versioning routers
	v1Router := chi.NewRouter()
	v1Router.HandleFunc("GET /healthz", handlerReadiness)
	v1Router.HandleFunc("GET /err", handlerErr)
	v1Router.HandleFunc("POST /create/user", pdb.handlerCreateUser)
	v1Router.HandleFunc("GET /user/info", pdb.middlewareAuth(pdb.handlerGetUser))
	v1Router.HandleFunc("POST /create/feeds", pdb.middlewareAuth(pdb.handlerCreateFeed))
	v1Router.HandleFunc("GET /all/feeds", pdb.handlerGetFeeds)
	v1Router.HandleFunc("POST /feed_follows", pdb.middlewareAuth(pdb.handlerCreateFeedFollow))
	v1Router.HandleFunc("GET /feed_follows", pdb.middlewareAuth(pdb.handlerGetFeedFollows))
	v1Router.HandleFunc("DELETE /feed_follows/{feedFollowId}", pdb.middlewareAuth(pdb.handlerDeleteFeedFollow))
	router.Mount("/v1", v1Router)
	//------------------------------

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Println("server starting on port :", portString)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Println("Server not listening ......", err)
		}
	}()

	<-done

	slog.Info("shutting down server")
	close(done)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	err = pdb.db.Close()
	if err != nil {
		slog.Error("error while closing psql connection pool", "message", err.Error())
	} else {
		slog.Info("database connectivity closed")
	}
	err = srv.Shutdown(ctx)
	if err != nil {
		slog.Error("server failed to shutdown", "message", err.Error())
	} else {
		slog.Info("server shutdown successful")
	}
}
