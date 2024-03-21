package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectToDB() {
	dbUser := os.Getenv("RAILS_DATABASE_USER")
	dbPassword := os.Getenv("RAILS_DATABASE_PASSWORD")
	dbName := os.Getenv("RAILS_DATABASE_NAME")
	dbHost := os.Getenv("RAILS_DATABASE_HOST")
	dbPort := os.Getenv("RAILS_DATABASE_PORT")

	// Construct DSN
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	if dsn == "" {
		log.Fatal("RAILS_DB environment variable is empty")
	}

	var err error
	DB, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected to database")

	migration()
}

func migration() {
	var err error
	_, err = DB.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		deleted_at TIMESTAMPTZ,
		title TEXT UNIQUE,
		price NUMERIC,
		description TEXT,
		category TEXT,
		image TEXT,
		rating REAL,
		count INT
	)`)

	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	// Create the users table
	_, err = DB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role INTEGER DEFAULT 2,
            cart_id UUID
        )
    `)

	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create the carts table
	_, err = DB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS carts (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            created_at TIMESTAMPTZ DEFAULT now(),
            updated_at TIMESTAMPTZ DEFAULT now(),
            deleted_at TIMESTAMPTZ,
            user_id INTEGER,
            products TEXT[] DEFAULT '{}'
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create carts table: %v", err)
	}

	fmt.Println("All migrations executed successfully")
}
