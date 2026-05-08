package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func updateData(ctx context.Context, db *pgxpool.Pool) error {
	sql := `
		UPDATE users
		SET age = $2
		WHERE id = $1
	`

	tag, err := db.Exec(ctx, sql, 1, 20)

	log.Printf("updated rows: %d", tag.RowsAffected())
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

	if err := updateData(ctx, db); err != nil {
		log.Fatal(err)
	}
}
