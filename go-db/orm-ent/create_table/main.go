package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// migration
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("connected")
}
