package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Product represents a product in the warehouse
type Product struct {
	Code     string
	Quantity int
	Bay      string
	Shelf    int
}

// isValidCharacter checks if a rune is a valid character for our data
func isValidCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == ','
}

// CleanCSVRecords removes unwanted characters and validates rows
func CleanCSVRecords(records [][]string) [][]string {
	var cleanedRecords [][]string
	for _, record := range records {
		var cleanedRecord []string
		for _, field := range record {
			// Remove any unwanted characters
			filteredField := strings.Map(func(r rune) rune {
				if isValidCharacter(r) {
					return r
				}
				return -1
			}, field)

			// Add cleaned field if not empty
			trimmedField := strings.TrimSpace(filteredField)
			if trimmedField != "" {
				cleanedRecord = append(cleanedRecord, trimmedField)
			}
		}

		// Check if the cleaned record has exactly 3 fields
		if len(cleanedRecord) == 3 {
			cleanedRecords = append(cleanedRecords, cleanedRecord)
		} else {
			fmt.Printf("Skipping malformed record: %v\n", record)
		}
	}
	return cleanedRecords
}

// ReadCSVFromReader reads the CSV file from an io.Reader
func ReadCSVFromReader(reader io.Reader) ([]Product, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Allow variable field count

	// Read all CSV records
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	// Skip the header row (assuming the first row is the header)
	if len(records) > 0 {
		records = records[1:] // Remove the first row (header)
	}

	// Clean the records and validate them
	cleanedRecords := CleanCSVRecords(records)

	var products []Product
	for i, record := range cleanedRecords {
		// Ensure each record has 3 fields (Code, Quantity, Location)
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid record on line %d: %v", i+2, record)
		}

		// Parse the quantity field
		quantity, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity on line %d: %v", i+2, err)
		}

		// Parse the location field (e.g., "A3 5")
		locationParts := strings.Fields(record[2])
		if len(locationParts) != 2 {
			return nil, fmt.Errorf("invalid location format on line %d: %v", i+2, record[2])
		}

		// Parse the shelf number
		shelf, err := strconv.Atoi(locationParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid shelf on line %d: %v", i+2, err)
		}

		// Create a product instance
		product := Product{
			Code:     record[0],
			Quantity: quantity,
			Bay:      locationParts[0],
			Shelf:    shelf,
		}

		products = append(products, product)
	}

	return products, nil
}

// SortAndSummarise sorts and summarizes the products
func SortAndSummarise(products []Product) []Product {
	// Sort products by bay and shelf
	sort.Slice(products, func(i, j int) bool {
		if products[i].Bay != products[j].Bay {
			return products[i].Bay < products[j].Bay
		}
		return products[i].Shelf < products[j].Shelf
	})

	// Summarize product quantities
	summary := make(map[string]*Product)
	for _, p := range products {
		key := p.Code + "," + p.Bay + " " + strconv.Itoa(p.Shelf)
		if existing, found := summary[key]; found {
			existing.Quantity += p.Quantity
		} else {
			summary[key] = &Product{
				Code:     p.Code,
				Quantity: p.Quantity,
				Bay:      p.Bay,
				Shelf:    p.Shelf,
			}
		}
	}

	// Convert the summary map back to a slice
	var summarizedProducts []Product
	for _, p := range summary {
		summarizedProducts = append(summarizedProducts, *p)
	}

	return summarizedProducts
}
