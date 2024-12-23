package fileprocessor

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"problem-2/config" // Import the Config package
	"problem-2/models"
	"strconv"
	"strings"
)

// ParseFile parses the content of a file and extracts relevant information.
func ParseFile(filePath string) (models.FileData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return models.FileData{}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var data models.FileData
	section := ""

	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		counter++
		// Identify section markers 0,1,2, 5, 6
		if counter < 13 {
			if strings.Contains(line, "*") {
				continue
			} else if counter > 0 && counter < 7 {
				section = "header"
			} else if counter >= 8 && counter < 12 {
				section = "variables"
			} else {
				section = "records"
			}
		}

		// Parse the current section
		switch section {
		case "header":
			parseHeader(line, &data)
		case "variables":
			parseVariables(line, &data)
		case "records":
			if counter == 12 {
				continue
			}
			record, err := parseRecord(line, &data)
			if err != nil {
				return models.FileData{}, fmt.Errorf("failed to parse SIM record: %w", err)
			}
			data.Records = append(data.Records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return models.FileData{}, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return data, nil
}

// parseHeader extracts header information.
func parseHeader(line string, data *models.FileData) {
	// Example header format: "Quantity: 100"
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Quantity":
		if qty, err := strconv.Atoi(value); err == nil {
			data.Header.Quantity = qty
		}
	}
}

// parseVariables extracts input variable details like the starting IMSI.
func parseVariables(line string, data *models.FileData) {
	// Example input variable format: "Var_In_List: IMSI=100000000000001"
	if strings.Contains(line, "IMSI") {
		imsi := strings.Split(line, ":")
		if len(imsi) == 2 {
			data.Header.VarIn = strings.TrimSpace(imsi[1])
		} else {
			log.Println("Starting IMSI value Not found")
		}
	}
}

// parseRecord parses a single SIM record line.
func parseRecord(line string, data *models.FileData) (models.SIMRecord, error) {
	// Example record format: "100000000000001,PIN1,PUK1,PIN2,PUK2,AAM1,KI_UMTS_ENC,ACC"
	fields := strings.Split(line, ",")
	if len(fields) < 8 {
		return models.SIMRecord{}, errors.New("invalid SIM record format")
	}
	validatImsi := ValidateIMSIRange(data.Header.VarIn, fields[1], data.Header.Quantity)
	if !validatImsi {
		data.ValidationResults.InvalidImsiRangeCount++
	}
	return models.SIMRecord{
		IMSI:      fields[1],
		PIN1:      fields[0],
		PUK1:      fields[2],
		PIN2:      fields[3],
		PUK2:      fields[4],
		AAM1:      fields[5],
		KIUMTSEnc: fields[6],
		ACC:       fields[7],
	}, nil
}

// validating the imsi range is in between
func ValidateIMSIRange(startingImsi, imis string, quantity int) bool {
	startingImsiInt := new(big.Int)
	_, ok := startingImsiInt.SetString(startingImsi, 10) // Parse the string as base 10
	if !ok {
		fmt.Println("Error: value out of range")
	}
	imisInt := new(big.Int)
	_, ok = imisInt.SetString(imis, 10) // Parse the string as base 10
	if !ok {
		fmt.Println("Error: value out of range")
	}
	bigIntQuantity := new(big.Int)
	bigIntQuantity.SetInt64(int64(quantity))
	result := new(big.Int)
	result.Add(startingImsiInt, bigIntQuantity)
	return imisInt.Cmp(result) <= 0
}

// ScanFolder scans a directory for new files.
func ScanFolder(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %s: %w", folderPath, err)
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, fmt.Sprintf("%s/%s", folderPath, file.Name()))
		}
	}

	return filePaths, nil
}

// ProcessFile processes an individual file by parsing, validating, and persisting its data.
func ProcessFile(filePath string, cfg config.Config) error {
	log.Printf("Processing file: %s\n", filePath)

	// Parse the file
	data, err := ParseFile(filePath)
	if err != nil {
		createResultFile(filePath, cfg.OutputFolder, "nok")
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}
	// Validate and process records
	if data, err = ValidateFile(data, cfg.DBHandler); err != nil {
		createResultFile(filePath, cfg.OutputFolder, "nok")
		return fmt.Errorf("validation failed for file %s: %w", filePath, err)
	}
	defer PrintFileProcessingResult(data)
	// Insert records into the database
	for _, record := range data.Records {
		if err := cfg.DBHandler.InsertSIMRecord(record); err != nil {
			createResultFile(filePath, cfg.OutputFolder, "nok")
			return fmt.Errorf("failed to insert record for file %s: %w", filePath, err)
		}
	}

	// Mark file as successfully processed
	createResultFile(filePath, cfg.OutputFolder, "ok")
	return nil
}
func PrintFileProcessingResult(data models.FileData) {
	log.Println("---------------- File Processing Report-----------------")
	log.Println("Total Records In File Quantity: ", data.Header.Quantity)
	log.Println("Total Records Processed(Inserted in DB): ", (data.Header.Quantity - data.ValidationResults.AlreadyExistingImsi - data.ValidationResults.DuplicateImsiInFile - data.ValidationResults.InvalidImsiRangeCount))
	log.Println("Total Duplicate IMSI Records In File: ", data.ValidationResults.DuplicateImsiInFile)
	log.Println("Total Already Existing IMSI Records: ", data.ValidationResults.AlreadyExistingImsi)
	log.Println("Total Invalid Existing IMSI Records: ", data.ValidationResults.InvalidImsiRangeCount)
	log.Println("---------------- File Processing Completed-----------------")
}

// createResultFile creates an output file indicating the processing result (.ok or .nok).
func createResultFile(inputFilePath, outputFolder, result string) {
	fileName := fmt.Sprintf("%s.%s", strings.TrimSuffix(filepath.Base(inputFilePath), filepath.Ext(inputFilePath)), result)
	outputFilePath := fmt.Sprintf("%s/%s", outputFolder, fileName)

	file, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create result file %s: %v\n", outputFilePath, err)
		return
	}
	defer file.Close()

	log.Printf("Created result file: %s\n", outputFilePath)
}
