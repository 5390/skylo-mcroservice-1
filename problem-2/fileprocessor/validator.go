package fileprocessor

import (
	"fmt"
	"log"
	"problem-2/db"
	"problem-2/models"
)

// ValidateFile checks for duplicates within the file and against the database.
func ValidateFile(data models.FileData, dbHandler *db.DBHandler) (models.FileData, error) {
	seenIMSI := make(map[string]bool)
	if data.Header.Quantity != len(data.Records) {
		return models.FileData{}, fmt.Errorf("quantity mismatch between records: %v and header quantity : %v", len(data.Records), data.Header.Quantity)
	}
	// Iterate over the records
	for i := 0; i < len(data.Records); i++ {
		record := data.Records[i]

		// Check for duplicate IMSI within the file
		if seenIMSI[record.IMSI] {
			// Print the duplicate IMSI value
			log.Printf("Duplicate IMSI found:", record.IMSI)
			data.ValidationResults.DuplicateImsiInFile++

			// Remove the duplicate IMSI from the slice
			data.Records = append(data.Records[:i], data.Records[i+1:]...)

			// Decrement the index to account for the shifted elements
			i--
			continue
		}
		seenIMSI[(record.IMSI)] = true

		// Check for existing IMSI in the database
		exists, err := dbHandler.IMSIExists(record.IMSI)
		if err != nil {
			return models.FileData{}, err
		}
		if exists {
			log.Println("IMSI already exists in database :", record.IMSI)
			data.ValidationResults.AlreadyExistingImsi++
			continue
		}
	}
	return data, nil
}
