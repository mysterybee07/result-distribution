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

func ParseColleges(filePath string) ([]models.College, error) {
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
		if len(record) < 6 {
			log.Printf("Invalid record length: %+v", record)
			continue
		}

		// Parse latitude and longitude as floats
		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		if err != nil {
			log.Printf("Error parsing latitude for college %s: %v", record[1], err)
			continue
		}

		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		if err != nil {
			log.Printf("Error parsing longitude for college %s: %v", record[1], err)
			continue
		}

		// Parse is_center as a boolean
		isCenter, err := strconv.ParseBool(strings.TrimSpace(record[5]))
		if err != nil {
			log.Printf("Error parsing is_center for college %s: %v", record[1], err)
			continue
		}

		collegeCode := strings.TrimSpace(record[0])

		// Check if college with this CollegeCode already exists in the database
		var existingCollege models.College
		if err := initializers.DB.Where("college_code = ?", collegeCode).First(&existingCollege).Error; err == nil {
			log.Printf("Duplicate college code found, skipping entry: %s", collegeCode)
			continue
		}

		// Create a new College entry if it doesn't exist
		college := models.College{
			CollegeCode: collegeCode,
			CollegeName: strings.TrimSpace(record[1]),
			Address:     strings.TrimSpace(record[2]),
			Latitude:    latitude,
			Longitude:   longitude,
			IsCenter:    isCenter,
		}

		colleges = append(colleges, college)
	}

	// Save parsed colleges to the database if any new colleges are found
	if len(colleges) == 0 {
		return nil, fmt.Errorf("no new colleges to store in the database")
	}

	if err := initializers.DB.Create(&colleges).Error; err != nil {
		return nil, fmt.Errorf("failed to store colleges in database: %w", err)
	}

	return colleges, nil
}
