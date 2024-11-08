package utils

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

type CenterAssignment struct {
	CollegeName   string
	CenterName    string
	AssignedSeat  int
	RemainingSeat int
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

		// Find available centers within 50km, excluding self-assignment and circular assignment
		var availableCenters []models.CapacityAndCount
		for _, center := range centers {
			// Prevent self-assignment and ensure remaining capacity
			if capCount.CollegeID == center.CollegeID || remainingCapacities[center.CollegeID] <= 0 {
				continue
			}
			// Prevent circular assignment
			if hasAssignment(assignments, center.College.CollegeName, college.CollegeName) {
				continue
			}
			// Check distance
			distance := Haversine(college.Latitude, college.Longitude, center.College.Latitude, center.College.Longitude)
			if distance < 5 {
				availableCenters = append(availableCenters, center)
			}
		}

		// Warning if no available centers found
		if len(availableCenters) == 0 {
			fmt.Printf("Warning: No available centers found for %s\n", college.CollegeName)
			continue
		}

		// Shuffle available centers to randomize assignment
		rand.Shuffle(len(availableCenters), func(i, j int) {
			availableCenters[i], availableCenters[j] = availableCenters[j], availableCenters[i]
		})

		// Assign students to centers
		for _, center := range availableCenters {
			if remainingStudents <= 0 {
				break
			}

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

// hasAssignment checks if students from centerName are already assigned to collegeName.
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
