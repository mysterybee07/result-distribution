package initializers

import (
	"log"

	"github.com/mysterybee07/result-distribution-system/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedBatches() {
	var count int64

	DB.Model(&models.Batch{}).Count(&count)

	if count == 0 {
		batches := []models.Batch{
			{Batch: 2020},
		}
		if err := DB.Create(&batches).Error; err != nil {
			log.Println("Failed to seed Batches:", err)
		} else {
			log.Println("Batches seeded successfully")
		}
	}
}

func SeedProgramsAndSemesters() {
	var count int64
	DB.Model(&models.Program{}).Count(&count)

	if count == 0 {
		// Define programs
		programs := []models.Program{
			{ProgramName: "Computer Science"},
			{ProgramName: "Information Technology"},
		}

		// Insert programs
		if err := DB.Create(&programs).Error; err != nil {
			log.Println("Failed to seed programs:", err)
			return
		}
		log.Println("Programs seeded successfully")
		// Seed semesters for each program
		for _, program := range programs {
			semesters := []models.Semester{
				{Name: 1, ProgramID: program.ID},
				{Name: 2, ProgramID: program.ID},
				{Name: 3, ProgramID: program.ID},
				{Name: 4, ProgramID: program.ID},
			}

			if err := DB.Create(&semesters).Error; err != nil {
				log.Println("Failed to seed semesters for program:", program.ProgramName, err)
			} else {
				log.Println("Semesters seeded successfully for program:", program.ProgramName)
			}
		}
	}
}
func SeedUsers() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		// Fetch existing Batch and Program to associate them with the User
		var batch models.Batch
		var program models.Program

		// Get the first batch
		if err := DB.First(&batch).Error; err != nil {
			log.Println("Failed to fetch the first batch:", err)
			return
		}
		// Get the first program
		if err := DB.First(&program).Error; err != nil {
			log.Println("Failed to fetch the first program:", err)
			return
		}

		log.Println("Using Batch ID:", batch.ID)
		log.Println("Using Program ID:", program.ID)

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password:", err)
			return
		}

		users := []models.User{
			{
				Symbol:   "SYM001",
				Email:    "admin@example.com",
				Password: string(hashedPassword),
				Role:     "admin",
				ImageURL: "/static/images/uploads/default.png",
				// No batch and program for admin
			},
			{
				BatchID:      &batch.ID,   // Should not be 0
				ProgramID:    &program.ID, // Should not be 0
				Symbol:       "SYM002",
				Registration: "REG002",
				Email:        "user1@example.com",
				Password:     string(hashedPassword),
				Role:         "user",
				ImageURL:     "/static/images/uploads/default.png",
			},
		}

		if err := DB.Create(&users).Error; err != nil {
			log.Println("Failed to seed users:", err)
		} else {
			log.Println("Users table seeded successfully")
		}
	}
}

func SeedStudents() {
	var count int64
	DB.Model(&models.Student{}).Count(&count)

	if count == 0 {

		var batch models.Batch
		var program models.Program

		// Get the first batch
		if err := DB.First(&batch).Error; err != nil {
			log.Println("Failed to fetch the first batch:", err)
			return
		}
		// Get the first program
		if err := DB.First(&program).Error; err != nil {
			log.Println("Failed to fetch the first program:", err)
			return
		}

		students := []models.Student{
			{
				BatchID:         batch.ID,
				ProgramID:       program.ID,
				SymbolNumber:    "SYM001",
				Registration:    "REG001",
				Fullname:        "Biraj Pudasaini",
				CurrentSemester: 1,
				Status:          "active",
			},
		}
		if err := DB.Create(&students).Error; err != nil {
			log.Println("Failed to seed users:", err)
		} else {
			log.Println("Students table seeded successfully")
		}
	}
}
