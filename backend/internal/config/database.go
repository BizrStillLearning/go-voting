package config

import (
	"fmt"
	"log"

	"go-voting/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=root123 dbname=go-voting port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke PostgreSQL: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Admin{},
		&models.Poll{},
		&models.Option{},
		&models.VoterToken{},
		&models.Vote{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	fmt.Println("Database PostgreSQL Berhasil Terhubung & Dimigrasi!")
	seedAdmin()
}

func seedAdmin() {
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)

	if adminCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.Admin{
			Username: "admin",
			Password: string(hashedPassword),
		}
		DB.Create(&admin)
		fmt.Println("Seeder: Akun admin default berhasil dibuat di Postgres!")
	}
}
