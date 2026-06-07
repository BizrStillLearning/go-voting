package config

import (
	"fmt"
	"log"

	"go-voting/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go-voting?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Admin{},
		&models.Poll{},
		&models.Option{},
		&models.VoterToken{},
		&models.Vote{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", err)
	}

	fmt.Println("Database berhasil diinisialisasi!")
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
		fmt.Println("Seeder: Akun admin default berhasil dibuat (Username: admin | Pass: admin123)")
	}
}
