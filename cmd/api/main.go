// File: todo/cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"module todo.kegodo.net/internal/data"
)

// configuration struct to hold configuration settings
type config struct {
	port int      //port on which the databased will open on
	env  string   //which environment we are working in (for this quiz it'll be development)
	db   struct { //database limiters and dependencies
		dsn          string //connection to databases
		maxOpenConns int    //limit of open connections
		maxIdleConns int    //limit of idle connections
		MaxIdleTime  string //limit on idle time
	}
}

// application struct is made to facilitate dependency injection
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

// main
func main() {
	//initializing where the configurations are gonna be stored
	var cfg config

	//hardcoding cofigurations since they are required for the quiz
	cfg.port = 4000
	cfg.env = "development"
	cfg.db.dsn = os.Getenv("TODO_DB_DSN")
	cfg.db.maxOpenConns = 25
	cfg.db.maxIdleConns = 25
	cfg.db.MaxIdleTime = "15m"

	//creating logger to log issues or state changes
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//creating connection
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	//ensuring that the connection to the database is closed
	defer db.Close()

	logger.Println("database connection pool established")

	//initializing the app struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	//initializing http server dependencies
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//staring the web server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// OpenDB() function returns a *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//Creating a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}