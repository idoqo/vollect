package main

import (
	"context"
	"fmt"
	"gitlab.com/idoko/vollect/db"
	"gitlab.com/idoko/vollect/handler"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error occurred: ", err)
	}

	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatal("Could not set up database: ", err)
	}
	defer database.Conn.Close()

	httpHandler := handler.NewHandler(database)
	server := &http.Server{
		Handler: httpHandler,
	}
	go func() {
		server.Serve(listener)
	}()
	defer stop(server)

	log.Printf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server...")
}

func stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not gracefully stop server, forcing shut down...")
		os.Exit(1)
	}
}
