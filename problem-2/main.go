package main

import (
	"log"
	"os"
	"problem-2/config"
	"problem-2/fileprocessor"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	var FileProcessedMap = make(map[string]struct{})
	log.Println("Starting SIM Data Processor...")

	// Periodic folder scanning
	for {
		files, err := fileprocessor.ScanFolder(cfg.InputFolder)
		if err != nil {
			log.Printf("Error scanning folder: %v\n", err)
			continue
		}

		for _, file := range files {
			if _, ok := FileProcessedMap[file]; ok {
				log.Printf("File Is Already Processed: %v\n", file)
				log.Printf("Deleting the file: %v\n", file)
				os.Remove(file)
				continue
			}
			err := fileprocessor.ProcessFile(file, cfg)
			if err != nil {
				log.Printf("Error processing file %s: %v\n", file, err)
			}
			FileProcessedMap[file] = struct{}{}
		}

		time.Sleep(cfg.ScanInterval)
	}
}
