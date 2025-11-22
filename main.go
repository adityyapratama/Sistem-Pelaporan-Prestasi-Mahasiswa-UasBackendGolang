package main

import (
	"context"
	"log"
	"os"
	"uas-pelaporan-prestasi-mahasiswa/database"

	// Sesuaikan dengan nama modul kamu di go.mod

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	// 2. Connect PostgreSQL
	pgDB, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal("❌ Gagal konek Postgres:", err)
	}
	defer pgDB.Close() // Tutup koneksi kalau aplikasi mati

	// 3. Connect MongoDB
	mongoClient, mongoDB, err := database.ConnectMongo()
	if err != nil {
		log.Fatal("❌ Gagal konek Mongo:", err)
	}
	// Disconnect Mongo kalau aplikasi mati
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnect mongo:", err)
		}
	}()

	// (Opsional) Print nama database mongo biar yakin
	log.Println("📂 Menggunakan Mongo Database:", mongoDB.Name())

	// 4. Init Fiber App
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server Running & Connected to DBs!",
			"status":  "OK",
		})
	})

	// 5. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}