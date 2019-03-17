package main // import "github.com/acoshift/wongnok"

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

	_ "github.com/lib/pq"

	"github.com/acoshift/wongnok/internal/api"
	"github.com/acoshift/wongnok/internal/auth"
	"github.com/acoshift/wongnok/internal/management"
)

func main() {
	fmt.Println("wongnok")
	fmt.Println("version: 1.0.0")

	dataSource := os.Getenv("DB_URL")

	db, err := sql.Open(
		"postgres",
		dataSource,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	server := http.Server{
		Addr: ":8080",
		Handler: api.API{
			Auth:       auth.New(db),
			Management: management.New(db),
		}.Handler(),
	}

	log.Printf("Server listening on %s\n", server.Addr)
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	<-stop
	fmt.Println()
	fmt.Println("^C again to force shutdown")
	go func() {
		<-stop
		fmt.Println()
		fmt.Println("force shutdown")
		os.Exit(0)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Println("can not graceful shutdown")
		return
	}
}
