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

// AssignCenters assigns colleges to centers based on capacity and distance,
// preventing circular assignment between centers.
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

		// Find available centers within 50km and exclude self-assignment and circular assignment
		for _, center := range centers {
			// Prevent a college that is also a center from assigning students to itself
			if college.CollegeCode == center.CollegeCode || remainingCapacities[center.CollegeCode] <= 0 {
				continue // Skip self or centers with no remaining capacity
			}
			// Circular assignment prevention: if College A is a center assigned to Center B, Center B cannot assign students to College A
			if hasAssignment(assignments, center.Name, college.Name) {
				continue // Skip if the center already has students assigned from the current college
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
