package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	exam "github.com/mysterybee07/result-distribution-system/models/exam"
)

const (
	PREF_DISTANCE_THRESHOLD = 2    // Preferred threshold distance in km
	ABS_DISTANCE_THRESHOLD  = 7    // Absolute threshold distance in km
	MIN_STUDENT_IN_CENTER   = 10   // Min. no of students from a school to be assigned to a center in normal circumstances
	STRETCH_CAPACITY_FACTOR = 0.02 // How much can center capacity be stretched if need arises
	PREF_CUTOFF             = -4   // Do not allocate students with pref score less than cutoff
)

// ParseColleges reads a TSV file and returns a slice of College structs
func ParseColleges(filePath string) ([]exam.College, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Specify tab delimiter for TSV

	var colleges []exam.College
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from file %s: %w", filePath, err)
	}

	// Iterate over records starting from index 1 to skip header row
	for _, record := range records[1:] {
		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		if err != nil {
			log.Printf("Error parsing latitude for college %s: %v", record[2], err)
			continue
		}
		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		if err != nil {
			log.Printf("Error parsing longitude for college %s: %v", record[2], err)
			continue
		}
		students, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			log.Printf("Error parsing student count for college %s: %v", record[2], err)
			continue
		}

		college := exam.College{
			CollegeCode:   strings.TrimSpace(record[0]),
			Name:          strings.TrimSpace(record[2]),
			Address:       strings.TrimSpace(record[2]),
			Latitude:      latitude,
			Longitude:     longitude,
			StudentsCount: students,
		}
		colleges = append(colleges, college)
	}
	return colleges, nil
}

// ParseCenters reads a TSV file and returns a slice of Center structs
func ParseCenters(filePath string) ([]exam.Center, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Set the delimiter to tab

	var centers []exam.Center
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from file %s: %w", filePath, err)
	}

	for _, record := range records[1:] {
		capacity, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			log.Printf("Error parsing capacity for center %s: %v", record[2], err)
			continue
		}
		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
		if err != nil {
			log.Printf("Error parsing latitude for center %s: %v", record[2], err)
			continue
		}
		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[7]), 64)
		if err != nil {
			log.Printf("Error parsing longitude for center %s: %v", record[2], err)
			continue
		}

		center := exam.Center{
			CenterCollegeCode: strings.TrimSpace(record[0]),
			Name:              strings.TrimSpace(record[2]),
			Address:           strings.TrimSpace(record[3]),
			Latitude:          latitude,
			Longitude:         longitude,
			Capacity:          capacity,
		}
		centers = append(centers, center)
	}
	return centers, nil
}

// ParsePreferences reads a TSV file and returns a slice of Preference structs
func ParsePreferences(filePath string) ([]exam.Preference, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Set the delimiter to tab

	var preferences []exam.Preference
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from file %s: %w", filePath, err)
	}

	// Skip the header and iterate over the records
	for _, record := range records[1:] {
		score, err := strconv.Atoi(strings.TrimSpace(record[2])) // Parsing preference score
		if err != nil {
			log.Printf("Error parsing preference score for college %s and center %s: %v", record[0], record[1], err)
			continue
		}

		// Append a new preference to the list
		preferences = append(preferences, exam.Preference{
			CollegeCode:       strings.TrimSpace(record[0]), // College code
			CenterCollegeCode: strings.TrimSpace(record[1]), // Center college code
			Preference:        score,                        // Preference score
		})
	}

	return preferences, nil
}

// HaversineDistance calculates the great circle distance between two points on the earth
// specified in decimal degrees.
func HaversineDistance(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	latitude1Rad, longitude1Rad := degreesToRadians(latitude1), degreesToRadians(longitude1)
	latitude2Rad, longitude2Rad := degreesToRadians(latitude2), degreesToRadians(longitude2)

	dlongitude := longitude2Rad - longitude1Rad
	dlatitude := latitude2Rad - latitude1Rad
	a := math.Sin(dlatitude/2)*math.Sin(dlatitude/2) +
		math.Cos(latitude1Rad)*math.Cos(latitude2Rad)*
			math.Sin(dlongitude/2)*math.Sin(dlongitude/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	radiusEarth := 6371.0 // Average radius of Earth in kilometers
	distance := radiusEarth * c
	return distance
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Function to handle assignment and create the required files when the API is hit
func AssignCenters(colleges []exam.College, centers []exam.Center, prefs []exam.Preference) ([]exam.AllocatedCenter, []exam.College) {
	var assignedCenters []exam.AllocatedCenter
	var unassignedColleges []exam.College

	for _, college := range colleges {
		assigned := false

		// Preference-based assignment logic
		for _, pref := range prefs {
			if pref.CollegeCode == college.CollegeCode && pref.Preference >= PREF_CUTOFF {
				for _, center := range centers {
					if center.CenterCollegeCode == pref.CenterCollegeCode {
						stretchedCapacity := int(float64(center.Capacity) * (1 + STRETCH_CAPACITY_FACTOR))

						if stretchedCapacity >= college.StudentsCount {
							distance := HaversineDistance(college.Latitude, college.Longitude, center.Latitude, center.Longitude)
							if distance <= PREF_DISTANCE_THRESHOLD {
								center.Capacity -= college.StudentsCount
								assignedCenters = append(assignedCenters, exam.AllocatedCenter{
									CollegeCode:       college.CollegeCode,
									CenterCollegeCode: center.CenterCollegeCode,
									Count:             college.StudentsCount,
								})
								assigned = true
								break
							}
						}
					}
				}
				if assigned {
					break
				}
			}
		}

		// Non-preferred center assignment
		if !assigned {
			for _, center := range centers {
				stretchedCapacity := int(float64(center.Capacity) * (1 + STRETCH_CAPACITY_FACTOR))

				if stretchedCapacity >= college.StudentsCount && college.StudentsCount >= MIN_STUDENT_IN_CENTER {
					distance := HaversineDistance(college.Latitude, college.Longitude, center.Latitude, center.Longitude)
					if distance <= ABS_DISTANCE_THRESHOLD {
						center.Capacity -= college.StudentsCount
						assignedCenters = append(assignedCenters, exam.AllocatedCenter{
							CollegeCode:       college.CollegeCode,
							CenterCollegeCode: center.CenterCollegeCode,
							Count:             college.StudentsCount,
						})
						assigned = true
						break
					}
				}
			}
		}

		// Track unassigned colleges
		if !assigned {
			unassignedColleges = append(unassignedColleges, college)
			log.Printf("School %s not assigned to any center\n", college.Name)
		}
	}

	return assignedCenters, unassignedColleges
}

// Function to write the assigned students data to a file
func WriteAssignedCentersToFile(assignedCenters []exam.AllocatedCenter, colleges []exam.College, centers []exam.Center) error {
	// Create a timestamped file for assigned centers
	fileName := fmt.Sprintf("assigned_centers_%d.csv", time.Now().Unix())
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Center Name", "College Name", "Assigned Students Count"})

	// Write assigned centers data
	for _, assigned := range assignedCenters {
		centerName := GetCenterNameByCode(assigned.CenterCollegeCode, centers)
		collegeName := GetCollegeNameByCode(assigned.CollegeCode, colleges)

		writer.Write([]string{
			centerName,
			collegeName,
			fmt.Sprintf("%d", assigned.Count),
		})
	}

	log.Printf("Assigned centers written to %s", fileName)
	return nil
}

// Function to write unassigned students data to a file
func WriteUnassignedCollegesToFile(unassignedColleges []exam.College) error {
	// Create a timestamped file for unassigned colleges
	fileName := fmt.Sprintf("unassigned_colleges_%d.csv", time.Now().Unix())
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"College Name", "Students Count"})

	// Write unassigned colleges data
	for _, college := range unassignedColleges {
		writer.Write([]string{
			college.Name,
			fmt.Sprintf("%d", college.StudentsCount),
		})
	}

	log.Printf("Unassigned colleges written to %s", fileName)
	return nil
}

func GetCenterNameByCode(centerCode string, centers []exam.Center) string {
	for _, center := range centers {
		if center.CenterCollegeCode == centerCode {
			return center.Name
		}
	}
	return "Unknown Center" // Return a default value if not found
}

// GetCollegeNameByCode fetches the college name by its code
func GetCollegeNameByCode(collegeCode string, colleges []exam.College) string {
	for _, college := range colleges {
		if college.CollegeCode == collegeCode {
			return college.Name
		}
	}
	return "Unknown College" // Return a default value if not found
}
