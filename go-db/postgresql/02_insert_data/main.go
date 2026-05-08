package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func insertData(ctx context.Context, db *pgxpool.Pool) error {
	sql := `
	INSERT INTO users(name, age)
	VALUES($1, $2)
	`

	_, err := db.Exec(ctx, sql, "田中", 20)
	return err
}

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := insertData(ctx, db); err != nil {
		log.Fatal(err)
	}
}
