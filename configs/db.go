package configs

import (
	"automotive/models"
	"automotive/utils"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := utils.GetEnv("DB_USERNAME", "test")
	password := utils.GetEnv("DB_PASSWORD", "test")
	host := utils.GetEnv("DB_HOST", "localhost")
	database := utils.GetEnv("DB_NAME", "automotive")
	schema := utils.GetEnv("DB_SCHEMA", "development")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=5432 sslmode=disable TimeZone=Asia/Jakarta search_path=%v", host, username, password, database, schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.Role{}, &models.User{}, &models.Brand{}, &models.Type{}, &models.Car{}, &models.Cart{})

	// Seed Admin Role
	adminDescription := "admin initial seeder"

	adminRole := models.Role{ID: 1, Name: "admin", Description: &adminDescription, CreatedAt: time.Now(), CreatedBy: "seeder", UpdatedBy: ""}

	var existingRole models.Role
	if err := db.Where("name = ?", adminRole.Name).First(&existingRole).Error; err != nil {
		adminRole.Insert(db)
	}

	// Seed Admin User
	defaultPass := utils.GetEnv("ADMIN_DEFAULT_PASS", "test")

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)
	if errPassword != nil {
		panic("Failed to Generate Password for Admin User")
	}

	adminUser := models.User{ID: 1, Name: "admin", RoleId: 1, Username: "admin", Password: string(hashedPassword), CreatedAt: time.Now(), CreatedBy: "seeder"}

	var existingUser models.User
	if err := db.Where("name = ?", adminUser.Name).First(&existingUser).Error; err != nil {
		db.Create(&adminUser)
	}

	return db
}
