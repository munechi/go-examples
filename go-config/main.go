package main

import (
	"context"
	"fmt"
	"go-config/config"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := config.Load()
	dsn := cfg.DB.DSN()

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	fmt.Println("connected!")
}
