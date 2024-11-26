package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
	"gorm.io/gorm"
)

type CenterAssignment struct {
	CollegeName   string
	CenterName    string
	AssignedSeat  int
	RemainingSeat int
}

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
		// Check the length of the record to ensure it has exactly 5 fields
		if len(record) < 5 {
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

// Helper function to process each record
func ProcessRecord(collegeName string, batchID, programID uint, isCenter bool, capacity int) error {
	// Look up the college_id using the college name
	var college models.College
	if err := initializers.DB.Where("college_name = ?", collegeName).First(&college).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("college not found for name: %s", collegeName)
		}
		return fmt.Errorf("failed to find college: %v", err)
	}

	// Check if a record for this college, batch, and program already exists
	var capacityAndCount models.CapacityAndCount
	result := initializers.DB.Where("college_id = ? AND batch_id = ? AND program_id = ?", college.ID, batchID, programID).First(&capacityAndCount)

	// Only create or update the record if it does not already exist
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create a new CapacityAndCount record
		capacityAndCount = models.CapacityAndCount{
			CollegeID:     college.ID,
			BatchID:       batchID,
			ProgramID:     programID,
			StudentsCount: 0,
			IsCenter:      isCenter,
			Capacity:      capacity,
		}
		if err := initializers.DB.Create(&capacityAndCount).Error; err != nil {
			return fmt.Errorf("failed to create new record: %v", err)
		}
	} else if result.Error == nil {
		// If the record exists, update only if fields need modification
		needsUpdate := false
		if capacityAndCount.IsCenter != isCenter {
			capacityAndCount.IsCenter = isCenter
			needsUpdate = true
		}
		if capacityAndCount.Capacity != capacity {
			capacityAndCount.Capacity = capacity
			needsUpdate = true
		}
		if needsUpdate {
			if err := initializers.DB.Save(&capacityAndCount).Error; err != nil {
				return fmt.Errorf("failed to update record: %v", err)
			}
		}
	}

	return nil
}

func AssignCenters(batchID uint, programID uint) ([]CenterAssignment, error) {
	// Fetch CapacityAndCount entries filtered by batch and program
	var capacityAndCounts []models.CapacityAndCount
	if err := initializers.DB.Preload("College").
		Where("batch_id = ? AND program_id = ?", batchID, programID).
		Find(&capacityAndCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch capacity and count data: %w", err)
	}

	// Identify centers and load their capacity from CapacityAndCount
	var centers []models.CapacityAndCount
	if err := initializers.DB.Preload("College").
		Where("capacity > 0 AND batch_id = ? AND program_id = ?", batchID, programID).
		Find(&centers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch centers: %w", err)
	}

	// Prepare to store assignments and track remaining capacities
	var assignments []CenterAssignment
	remainingCapacities := make(map[uint]int) // key is CollegeID, value is remaining capacity

	// Initialize remaining capacities of each center
	for _, center := range centers {
		remainingCapacities[center.CollegeID] = center.Capacity
	}

	rand.Seed(time.Now().UnixNano())

	// Assign students from each filtered college to available centers
	for _, capCount := range capacityAndCounts {
		college := capCount.College
		remainingStudents := capCount.StudentsCount
		totalAssigned := 0

		// Assign students to centers
		for remainingStudents > 0 {
			// Filter available centers within 50km, excluding self-assignment and circular assignment
			availableCenters := filterCentersWithCyclePrevention(centers, capCount, assignments, remainingCapacities)
			if len(availableCenters) == 0 {
				fmt.Printf("Warning: No available centers found for %s\n", college.CollegeName)
				break
			}

			// Pick a center using weighted randomness
			center := weightedRandom(availableCenters, remainingCapacities)
			if center == nil {
				break
			}

			// Determine how many students to assign
			assignCount := min(remainingStudents, remainingCapacities[center.CollegeID])
			if assignCount > 0 {
				assignments = append(assignments, CenterAssignment{
					CollegeName:   college.CollegeName,
					CenterName:    center.College.CollegeName,
					AssignedSeat:  assignCount,
					RemainingSeat: remainingCapacities[center.CollegeID] - assignCount,
				})

				remainingStudents -= assignCount
				totalAssigned += assignCount
				remainingCapacities[center.CollegeID] -= assignCount
			}
		}

		// Log unassigned students if any remain
		if remainingStudents > 0 {
			fmt.Printf("Warning: %s still has unassigned students: %d\n", college.CollegeName, remainingStudents)
		}
	}

	return assignments, nil
}

func weightedRandom(centers []models.CapacityAndCount, capacities map[uint]int) *models.CapacityAndCount {
	// Calculate total remaining capacity of available centers
	totalCapacity := 0
	for _, center := range centers {
		totalCapacity += capacities[center.CollegeID]
	}

	if totalCapacity == 0 {
		return nil // No available capacity
	}

	// Generate a random number within the total capacity
	randValue := rand.Intn(totalCapacity)
	runningSum := 0

	// Select a center based on the random number
	for _, center := range centers {
		runningSum += capacities[center.CollegeID]
		if randValue < runningSum {
			return &center
		}
	}
	return nil
}

func filterCentersWithCyclePrevention(
	centers []models.CapacityAndCount,
	capCount models.CapacityAndCount,
	assignments []CenterAssignment,
	remainingCapacities map[uint]int,
) []models.CapacityAndCount {
	var availableCenters []models.CapacityAndCount

	for _, center := range centers {
		// Prevent self-assignment
		if capCount.CollegeID == center.CollegeID {
			continue
		}

		// Skip centers with no remaining capacity
		if remainingCapacities[center.CollegeID] <= 0 {
			continue
		}

		// Prevent circular assignments
		if hasAssignment(assignments, center.College.CollegeName, capCount.College.CollegeName) ||
			hasAssignment(assignments, capCount.College.CollegeName, center.College.CollegeName) {
			continue
		}

		// Check distance (within 50km)
		distance := Haversine(capCount.College.Latitude, capCount.College.Longitude, center.College.Latitude, center.College.Longitude)
		if distance > 50 {
			continue
		}

		availableCenters = append(availableCenters, center)
	}

	return availableCenters
}

func hasAssignment(assignments []CenterAssignment, centerName, collegeName string) bool {
	for _, assignment := range assignments {
		if assignment.CenterName == centerName && assignment.CollegeName == collegeName {
			return true
		}
	}
	return false
}

func WriteResultToFile(assignments []CenterAssignment) error {
	// Ensure the 'data' folder exists, create it if not
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create the CSV file inside the 'data' folder
	file, err := os.Create("data/center_assignment_results.csv")
	if err != nil {
		return fmt.Errorf("failed to create result file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	if err := writer.Write([]string{"College Name", "Assigned Center", "Assigned Students"}); err != nil {
		return fmt.Errorf("failed to write header to CSV: %w", err)
	}

	// Write each assignment
	for _, assignment := range assignments {
		if err := writer.Write([]string{
			assignment.CollegeName,
			assignment.CenterName,
			strconv.Itoa(assignment.AssignedSeat),
		}); err != nil {
			return fmt.Errorf("failed to write assignment to file: %w", err)
		}
	}

	return nil
}
