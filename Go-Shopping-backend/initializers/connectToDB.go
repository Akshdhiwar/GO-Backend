package initializers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectToDB() {
	var dbUser string

	if os.Getenv("RAILS_ENVIRONMENT") == "LOCAL" {
		dbUser = os.Getenv("RAILS_DATABASE_USER")
	} else {
		dbUser = os.Getenv("RAILS_DATABASE_USER_PROD")
	}

	log.Println(dbUser)

	var dbPassword string

	if os.Getenv("RAILS_ENVIRONMENT") == "LOCAL" {
		dbPassword = os.Getenv("RAILS_DATABASE_PASSWORD")
	} else {
		dbPassword = os.Getenv("RAILS_DATABASE_PASSWORD_PROD")
	}
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
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

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
		count INT,
		price_id TEXT
	)`)

	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	// // Create the users table
	// _, err = DB.Exec(context.Background(), `
	//     CREATE TABLE IF NOT EXISTS users (
	//         id SERIAL PRIMARY KEY,
	//         email TEXT UNIQUE NOT NULL,
	//         password TEXT NOT NULL,
	// 		first_name TEXT NOT NULL,
	// 		last_name TEXT NOT NULL,
	//         role INTEGER DEFAULT 2,
	//         cart_id UUID
	//     )
	// `)

	// if err != nil {
	// 	log.Fatalf("Failed to create users table: %v", err)
	// }

	// Create the carts table
	_, err = DB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS carts (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            created_at TIMESTAMPTZ DEFAULT now(),
            updated_at TIMESTAMPTZ DEFAULT now(),
            deleted_at TIMESTAMPTZ,
            user_id UUID,
            products JSONB[] DEFAULT '{}'::JSONB[]
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create carts table: %v", err)
	}

	_, err = DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ DEFAULT now(),
    		email TEXT NOT NULL,
    		products JSONB[] DEFAULT '{}'::JSONB[]
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}

	fmt.Println("All migrations executed successfully")
}
