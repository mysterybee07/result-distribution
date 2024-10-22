package utils

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/mysterybee07/result-distribution-system/initializers"
	models "github.com/mysterybee07/result-distribution-system/models/exam"
)

type CenterAssignment struct {
	CollegeName   string
	CenterName    string
	AssignedSeat  int
	RemainingSeat int // Changed to int if it should represent remaining capacity
}

func WriteResultToFile(assignments []CenterAssignment) error {
	file, err := os.Create("center_assignment_results.csv")
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

// AssignCenters assigns colleges to centers based on capacity and distance
// AssignCenters assigns colleges to centers based on capacity and distance
func AssignCenters() ([]CenterAssignment, error) {
	var colleges []models.College
	if err := initializers.DB.Find(&colleges).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch colleges: %w", err)
	}

	// Identify the colleges that are centers
	var centers []models.College
	for _, college := range colleges {
		if college.IsCenter {
			centers = append(centers, college)
		}
	}

	// Prepare to store assignments and remaining capacities for each center
	var assignments []CenterAssignment
	remainingCapacities := make(map[string]int)

	// Initialize remaining capacities of each center
	for _, center := range centers {
		remainingCapacities[center.CollegeCode] = center.Capacity
	}

	rand.Seed(time.Now().UnixNano())

	// Assign students from each college to available centers
	for _, college := range colleges {
		remainingStudents := college.StudentsCount
		totalAssigned := 0 // To track the total number of students assigned for this college
		var availableCenters []models.College

		// Find available centers within 50km and exclude self-assignment
		for _, center := range centers {
			// Prevent a college that is also a center from assigning students to itself
			if college.CollegeCode == center.CollegeCode || remainingCapacities[center.CollegeCode] <= 0 {
				continue // Skip self or centers with no remaining capacity
			}
			distance := Haversine(college.Latitude, college.Longitude, center.Latitude, center.Longitude)
			if distance < 50 { // Check distance constraint
				availableCenters = append(availableCenters, center)
			}
		}

		// If no available centers, log a warning
		if len(availableCenters) == 0 {
			fmt.Printf("Warning: No available centers found for %s\n", college.Name)
			continue // No available centers; move to the next college
		}

		// Shuffle available centers to randomize assignment order
		rand.Shuffle(len(availableCenters), func(i, j int) {
			availableCenters[i], availableCenters[j] = availableCenters[j], availableCenters[i]
		})

		// Assign students to centers
		for _, center := range availableCenters {
			if remainingStudents <= 0 {
				break // Exit if all students are assigned
			}

			assignCount := min(remainingStudents, remainingCapacities[center.CollegeCode])
			if assignCount > 0 {
				assignments = append(assignments, CenterAssignment{
					CollegeName:   college.Name,
					CenterName:    center.Name,
					AssignedSeat:  assignCount,
					RemainingSeat: remainingCapacities[center.CollegeCode] - assignCount,
				})

				remainingStudents -= assignCount
				totalAssigned += assignCount
				remainingCapacities[center.CollegeCode] -= assignCount // Update capacity of the center
			}
		}

		// Log remaining students for this college
		remainingAfterAssignment := college.StudentsCount - totalAssigned
		fmt.Printf("Remaining students for %s after assignment: %d\n", college.Name, remainingAfterAssignment)

		// Additional warning if there are unassigned students
		if remainingStudents > 0 {
			fmt.Printf("Warning: %s still has unassigned students: %d\n", college.Name, remainingStudents)
		}
	}

	return assignments, nil
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
