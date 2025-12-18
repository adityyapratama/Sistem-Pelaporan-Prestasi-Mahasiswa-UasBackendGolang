package main

import (
	"context"
	"log"
	"os"
	"uas-pelaporan-prestasi-mahasiswa/apps/repository"
	"uas-pelaporan-prestasi-mahasiswa/apps/service"
	"uas-pelaporan-prestasi-mahasiswa/config"
	"uas-pelaporan-prestasi-mahasiswa/database"

	"github.com/joho/godotenv"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	
	pgDB, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer pgDB.Close()

	mongoClient, _, err := database.ConnectMongo()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Println("Error disconnect mongo:", err)
		}
	}()

	mongoDbInstance := mongoClient.Database("db_prestasi")

	userRepo := repository.NewUserRepository(pgDB)
	permissionRepo :=repository.NewPostgresPermissionRepository(pgDB)
	studentRepo :=repository.NewStudentRepository(pgDB)
	lectureRepo :=repository.NewPostgresLectureRepository(pgDB)
	achievementRepo := repository.NewAchievementRepo(pgDB,mongoDbInstance)

	authService := service.NewAuthService(userRepo, permissionRepo)
	permService := service.NewPermissionService(permissionRepo)
	studentService := service.NewStudentService(studentRepo)
	lectureService :=service.NewLectureService(lectureRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo)


	
	app := config.NewApp(authService,permService,studentService,lectureService,achievementService)
	app.Static("/uploads", "./uploads")
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
	
	
}