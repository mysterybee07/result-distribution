package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandLetter(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SanitizeFileName(fileName string) string {
	// Remove any characters that are not alphanumeric, dot, or underscore
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return strings.ToLower(re.ReplaceAllString(fileName, "_"))
}

// ConvertIDs converts batch ID and program ID from string to pointers of uint.
// Returns an error if any conversion fails.
func ConvertIDs(batchIDStr, programIDStr string) (batchID, programID *uint, err error) {
	if batchIDStr != "" {
		batchUint, err := strconv.ParseUint(batchIDStr, 10, 32)
		if err != nil {
			return nil, nil, err // Return the error if conversion fails
		}
		batchIDVal := uint(batchUint)
		batchID = &batchIDVal // Return a pointer to the converted uint
	}

	if programIDStr != "" {
		programUint, err := strconv.ParseUint(programIDStr, 10, 32)
		if err != nil {
			return nil, nil, err // Return the error if conversion fails
		}
		programIDVal := uint(programUint)
		programID = &programIDVal // Return a pointer to the converted uint
	}

	return batchID, programID, nil // Return both pointers and nil error
}
