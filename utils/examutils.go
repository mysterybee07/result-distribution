package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

const (
	PREF_DISTANCE_THRESHOLD = 2    // Preferred threshold distance in km
	ABS_DISTANCE_THRESHOLD  = 7    // Absolute threshold distance in km
	MIN_STUDENT_IN_CENTER   = 10   // Min. no of students from a school to be assigned to a center in normal circumstances
	STRETCH_CAPACITY_FACTOR = 0.02 // How much can center capacity be stretched if need arises
	PREF_CUTOFF             = -4   // Do not allocate students with pref score less than cutoff
)

// ParseColleges reads a TSV file and returns a slice of College structs

func ParseColleges(filePath string, batchID uint, programID uint) ([]models.College, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Specify tab delimiter for TSV

	var colleges []models.College
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from file %s: %w", filePath, err)
	}

	// Iterate over records starting from index 1 to skip header row
	for _, record := range records[1:] {
		// Check the length of the record to avoid index out of range errors
		if len(record) < 8 {
			log.Printf("Invalid record length: %+v", record)
			continue
		}

		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64) // Index for latitude
		if err != nil {
			log.Printf("Error parsing latitude for college %s: %v", record[1], err)
			continue
		}

		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64) // Index for longitude
		if err != nil {
			log.Printf("Error parsing longitude for college %s: %v", record[1], err)
			continue
		}

		students, err := strconv.Atoi(strings.TrimSpace(record[5])) // Index for students_count
		if err != nil {
			log.Printf("Error parsing student count for college %s: %v", record[1], err)
			continue
		}

		isCenter, err := strconv.ParseBool(strings.TrimSpace(record[6])) // Index for is_center
		if err != nil {
			log.Printf("Error parsing is_center for college %s: %v", record[1], err)
			continue
		}

		capacity, err := strconv.Atoi(strings.TrimSpace(record[7])) // Index for capacity
		if err != nil {
			log.Printf("Error parsing capacity for college %s: %v", record[1], err)
			continue
		}

		// Create a new College entry
		college := models.College{
			CollegeCode:   strings.TrimSpace(record[0]),
			CollegeName:   strings.TrimSpace(record[1]),
			Address:       strings.TrimSpace(record[2]),
			Latitude:      latitude,
			Longitude:     longitude,
			StudentsCount: students,
			IsCenter:      isCenter,
			Capacity:      capacity,
			BatchID:       batchID,
			ProgramID:     programID,
		}

		colleges = append(colleges, college)
	}

	// Save parsed colleges to the database
	if err := initializers.DB.Create(&colleges).Error; err != nil {
		return nil, fmt.Errorf("failed to store colleges in database: %w", err)
	}

	return colleges, nil
}
