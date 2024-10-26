package utils

import (
	"log"

	"github.com/mysterybee07/result-distribution-system/initializers"
	"github.com/mysterybee07/result-distribution-system/models"
)

func GetPassStatusBySemester(semesterID string) (map[uint]string, error) {
	var marks []models.Mark
	if err := initializers.DB.Where("semester_id = ?", semesterID).Find(&marks).Error; err != nil {
		log.Printf("Failed to fetch marks: %v\n", err)
		return nil, err
	}

	// Track pass status for each student
	passStatus := make(map[uint]string)

	for _, mark := range marks {
		if mark.Status != "pass" {
			passStatus[mark.StudentID] = "fail"
		} else if _, ok := passStatus[mark.StudentID]; !ok {
			passStatus[mark.StudentID] = "pass"
		}
	}

	return passStatus, nil
}
