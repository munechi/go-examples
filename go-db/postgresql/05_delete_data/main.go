package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func deleteData(ctx context.Context, db *pgxpool.Pool) error {
	sql := `
		DELETE FROM users
		WHERE id = $1
	`

	tag, err := db.Exec(ctx, sql, 1)

	log.Printf("deleted rows: %d", tag.RowsAffected())
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

	if err := deleteData(ctx, db); err != nil {
		log.Fatal(err)
	}
}
