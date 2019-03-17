package main // import "github.com/acoshift/wongnok"

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
