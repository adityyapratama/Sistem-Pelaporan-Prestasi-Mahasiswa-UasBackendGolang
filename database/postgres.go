package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" 
)

func ConnectPostgres() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN tidak ditemukan di .env")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test Ping
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Setting Connection Pool (Biar performa bagus)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Berhasil terhubung ke PostgreSQL")
	return db, nil
}