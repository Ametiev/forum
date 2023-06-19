package main

import (
	"dyelesho/forum/internal/dbs"
	"dyelesho/forum/internal/handlers"
	"dyelesho/forum/internal/models"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := dbs.OpenDB()
	if err != nil {
		errorLog.Fatal(err)
	}

	dbs.CreatePosts(db)
	dbs.CreateTables(db)
	defer db.Close()
	templateCache, err := handlers.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &handlers.Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Posts:         &models.Model{DB: db},
		TemplateCache: templateCache,
		Users:         &models.UserModel{DB: db},
		Reactions:     &models.ReactionModel{DB: db},
	}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
