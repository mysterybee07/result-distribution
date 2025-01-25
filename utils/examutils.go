package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func ExamRoutine(batchID, programID, semesterID uint, startDate, endDate time.Time) (string, interface{}, error) {
	// Check for overlapping exams within a 20-day range
	overlapRangeStart := startDate.AddDate(0, 0, -20)
	overlapRangeEnd := endDate.AddDate(0, 0, 20)

	var overlappingExams []models.ExamRoutine
	if err := initializers.DB.Where(
		"program_id = ? AND (start_date BETWEEN ? AND ? OR end_date BETWEEN ? AND ?)",
		programID, overlapRangeStart, overlapRangeEnd, overlapRangeStart, overlapRangeEnd,
	).Find(&overlappingExams).Error; err != nil {
		return "", nil, fmt.Errorf("database error: %w", err)
	}

	if len(overlappingExams) > 0 {
		return "", nil, fmt.Errorf("overlapping exams detected: Ensure a 20-day gap between exams for the same program")
	}

	// Fetch courses
	var courses []models.Course
	if err := initializers.DB.Where("program_id = ? AND semester_id = ?", programID, semesterID).Find(&courses).Error; err != nil {
		return "", nil, fmt.Errorf("failed to fetch courses: %w", err)
	}

	if len(courses) == 0 {
		return "", nil, fmt.Errorf("no courses found for the given program and semester")
	}

	// Separate courses into compulsory and non-compulsory
	var compulsoryCourses, nonCompulsoryCourses []models.Course
	for _, course := range courses {
		if course.IsCompulsory {
			compulsoryCourses = append(compulsoryCourses, course)
		} else {
			nonCompulsoryCourses = append(nonCompulsoryCourses, course)
		}
	}

	// Shuffle compulsory courses for randomness
	compulsoryCourses = ShuffleCourses(compulsoryCourses)

	// Create the exam routine record
	examRoutine := models.ExamRoutine{
		StartDate:  startDate,
		EndDate:    endDate,
		BatchID:    batchID,
		ProgramID:  programID,
		SemesterID: semesterID,
		Status:     "NotPublished",
	}

	if err := initializers.DB.Create(&examRoutine).Error; err != nil {
		return "", nil, fmt.Errorf("failed to save exam routine: %w", err)
	}

	fileContent := "Course Code,Course Name,Exam Date\n"
	examSchedules := make([]models.ExamSchedules, 0)

	// Schedule compulsory courses with gap logic
	currentDate := startDate
	gap := calculateGap(startDate, endDate, len(compulsoryCourses))

	for _, course := range compulsoryCourses {
		for isWeekend(currentDate) {
			currentDate = currentDate.AddDate(0, 0, 1) // Skip weekends
		}

		fileContent += fmt.Sprintf("%s,%s,%s\n", course.CourseCode, course.Name, currentDate.Format("2006-01-02"))

		examSchedule := models.ExamSchedules{
			CourseID:      course.ID,
			ExamRoutineID: examRoutine.ID,
			ExamDate:      currentDate,
		}

		if err := initializers.DB.Create(&examSchedule).Error; err != nil {
			return "", nil, fmt.Errorf("failed to save exam schedule: %w", err)
		}

		examSchedules = append(examSchedules, examSchedule)
		currentDate = currentDate.AddDate(0, 0, gap) // Move to the next exam date
	}

	// Schedule non-compulsory courses on the last day (endDate)
	for _, course := range nonCompulsoryCourses {
		for isWeekend(endDate) {
			endDate = endDate.AddDate(0, 0, -1) // Avoid weekends
		}

		fileContent += fmt.Sprintf("%s,%s,%s\n", course.CourseCode, course.Name, endDate.Format("2006-01-02"))

		examSchedule := models.ExamSchedules{
			CourseID:      course.ID,
			ExamRoutineID: examRoutine.ID,
			ExamDate:      endDate,
		}

		if err := initializers.DB.Create(&examSchedule).Error; err != nil {
			return "", nil, fmt.Errorf("failed to save non-compulsory exam schedule: %w", err)
		}

		examSchedules = append(examSchedules, examSchedule)
	}

	// Save the exam routine to a CSV file
	fileName := fmt.Sprintf("ExamRoutine_Batch%d_Program%d_Semester%d.csv", batchID, programID, semesterID)

	if err := os.WriteFile(fileName, []byte(fileContent), 0644); err != nil {
		return "", nil, fmt.Errorf("failed to write to file: %w", err)
	}

	return fileName, examSchedules, nil
}

// Helper Functions
func calculateGap(startDate, endDate time.Time, numCourses int) int {
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	if numCourses == 0 {
		return 0
	}
	return totalDays / numCourses
}

// Helper function to check if a date falls on a weekend
func isWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}

// shuffleCourses randomizes the course order
func ShuffleCourses(courses []models.Course) []models.Course {
	shuffled := make([]models.Course, len(courses))
	copy(shuffled, courses)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

// calculateFitness evaluates the quality of each schedule
func CalculateFitness(population [][]models.Course, startDate time.Time) []int {
	fitnessScores := make([]int, len(population))
	for i, schedule := range population {
		score := 0
		currentDate := startDate
		for _, course := range schedule {
			if course.IsCompulsory {
				if currentDate.Weekday() == time.Saturday {
					currentDate = currentDate.AddDate(0, 0, 1) // Skip Saturday
				}
				score++
			}
			currentDate = currentDate.AddDate(0, 0, rand.Intn(3)+1) // Add random gap
		}
		fitnessScores[i] = score
	}
	return fitnessScores
}

// selectBest selects the top schedules based on fitness
func SelectBest(population [][]models.Course, fitnessScores []int) [][]models.Course {
	// Sort by fitness scores and select the best half
	bestIndexes := make([]int, len(fitnessScores))
	for i := range fitnessScores {
		bestIndexes[i] = i
	}

	// Sort indices based on scores
	for i := 0; i < len(fitnessScores)-1; i++ {
		for j := i + 1; j < len(fitnessScores); j++ {
			if fitnessScores[i] < fitnessScores[j] {
				bestIndexes[i], bestIndexes[j] = bestIndexes[j], bestIndexes[i]
			}
		}
	}

	selected := make([][]models.Course, len(population)/2)
	for i := 0; i < len(population)/2; i++ {
		selected[i] = population[bestIndexes[i]]
	}
	return selected
}

// crossover combines two schedules to produce offspring
func Crossover(parent1, parent2 []models.Course) ([]models.Course, []models.Course) {
	crossoverPoint := rand.Intn(len(parent1))
	child1 := append([]models.Course{}, parent1[:crossoverPoint]...)
	child1 = append(child1, parent2[crossoverPoint:]...)

	child2 := append([]models.Course{}, parent2[:crossoverPoint]...)
	child2 = append(child2, parent1[crossoverPoint:]...)

	return child1, child2
}

// mutate makes small random changes to a schedule
func Mutate(schedule []models.Course) {
	if len(schedule) > 1 {
		i, j := rand.Intn(len(schedule)), rand.Intn(len(schedule))
		schedule[i], schedule[j] = schedule[j], schedule[i]
	}
}

// getBestSchedule returns the most optimized schedule from the population
func GetBestSchedule(population [][]models.Course) []models.Course {
	return population[0] // Assuming population is sorted by fitness
}
