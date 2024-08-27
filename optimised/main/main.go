package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/google/uuid"
)

type InNetworkFile struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

func main() {
	// Connect to ClickHouse
	fmt.Println("Connecting to ClickHouse...")
	db, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000")
	if err != nil {
		log.Fatalf("Error connecting to ClickHouse: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing db connection: %v", err)
		}
	}(db)
	fmt.Println("Connected to ClickHouse")

	// Open JSON file
	fmt.Println("Opening JSON file...")
	file, err := os.Open("../../2024-08-01_anthem_index.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(file)
	fmt.Println("Opened JSON file")

	decoder := json.NewDecoder(file)
	var token json.Token

	// Navigate to the start of the reporting_structure array
	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			log.Fatalf("Error reading token: %v", err)
		}
		if delim, ok := token.(json.Delim); ok && delim.String() == "[" {
			break
		}
	}

	// Batch processing
	fmt.Println("Started processing batches...")
	start := time.Now()
	batchSize := 10_000
	var batch []InNetworkFile

	// Read through each reporting_structure item
	for decoder.More() {
		var reportingStructure struct {
			InNetworkFiles []InNetworkFile `json:"in_network_files"`
		}
		if err := decoder.Decode(&reportingStructure); err != nil {
			log.Fatalf("Error decoding reporting_structure: %v", err)
		}

		// Process in_network_files for each reporting_structure
		for _, inNetworkFile := range reportingStructure.InNetworkFiles {
			batch = append(batch, inNetworkFile)

			if len(batch) >= batchSize {
				if err := insertBatch(db, batch); err != nil {
					log.Fatalf("Error inserting batch: %v", err)
				}
				batch = nil // Reset batch
			}
		}
	}

	// Insert the final batch if it has any data
	if len(batch) > 0 {
		if err := insertBatch(db, batch); err != nil {
			log.Fatalf("Error inserting final batch: %v", err)
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("Data successfully parsed and inserted into ClickHouse in %.2f seconds.\n", elapsed)
}

func insertBatch(db *sql.DB, batch []InNetworkFile) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO urls (id, description, location) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatalf("Error closing statement: %v", err)
		}
	}(stmt)

	for _, obj := range batch {
		if _, err := stmt.Exec(uuid.New(), obj.Description, obj.Location); err != nil {
			return err
		}
	}

	return tx.Commit()
}
