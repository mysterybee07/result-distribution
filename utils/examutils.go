package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func ExamRoutine(batchID, programID, semesterID uint, startDate, endDate time.Time) (string, error) {
	// Check for overlapping exams
	overlapRangeStart := startDate.AddDate(0, 0, -20)
	overlapRangeEnd := endDate.AddDate(0, 0, 20)

	var overlappingExams []models.ExamRoutine
	if err := initializers.DB.Where(
		"program_id = ? AND (start_date BETWEEN ? AND ? OR end_date BETWEEN ? AND ?)",
		programID, overlapRangeStart, overlapRangeEnd, overlapRangeStart, overlapRangeEnd,
	).Find(&overlappingExams).Error; err != nil {
		return "", fmt.Errorf("database error: %w", err)
	}

	if len(overlappingExams) > 0 {
		return "", fmt.Errorf("overlapping exams detected: ensure a 20-day gap between exams for the same program")
	}

	// Fetch courses for the semester and program
	var courses []models.Course
	if err := initializers.DB.Where(
		"program_id = ? AND semester_id = ?",
		programID, semesterID,
	).Find(&courses).Error; err != nil {
		return "", fmt.Errorf("failed to fetch courses: %w", err)
	}

	if len(courses) == 0 {
		return "", fmt.Errorf("no courses found for the given program and semester")
	}

	// Validate date range
	days := int(endDate.Sub(startDate).Hours() / 24)
	if len(courses) > days {
		return "", fmt.Errorf("not enough days in the range to schedule all exams")
	}

	// Generate exam schedule
	rand.Seed(time.Now().UnixNano())
	usedDates := make(map[string]bool)
	fileContent := "Course Code,Course Name,Exam Date\n"

	for _, course := range courses {
		var examDate time.Time
		for {
			randomDay := rand.Intn(days)
			examDate = startDate.AddDate(0, 0, randomDay)
			if !usedDates[examDate.Format("2006-01-02")] {
				usedDates[examDate.Format("2006-01-02")] = true
				break
			}
		}
		fileContent += fmt.Sprintf("%s,%s,%s\n", course.CourseCode, course.Name, examDate.Format("2006-01-02"))
	}

	// Save exam routine to database
	examRoutine := models.ExamRoutine{
		StartDate:  startDate,
		EndDate:    endDate,
		BatchID:    batchID,
		ProgramID:  programID,
		SemesterID: semesterID,
	}

	if err := initializers.DB.Create(&examRoutine).Error; err != nil {
		return "", fmt.Errorf("failed to save routine details: %w", err)
	}

	// Write the exam schedule to a CSV file
	fileName := fmt.Sprintf("ExamRoutine_Batch%d_Program%d_Semester%d.csv", batchID, programID, semesterID)
	if err := os.WriteFile(fileName, []byte(fileContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return fileName, nil
}
