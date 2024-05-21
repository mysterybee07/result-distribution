package initializers

import (
	"log"
	"os"

	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	LoadEnvironment()
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not set in environment variable")
	}
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Migration in database
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Batch{},
		&models.Program{},
		&models.Semester{},
		&models.Subject{},
		&models.Student{},
	); err != nil {
		log.Fatalf("Error migrating database to database")
	}
}
