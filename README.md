# CSV Processor Application

This project is a web-based CSV file processor built using **Go**. It allows users to upload a CSV file, processes it on the backend, sorts and summarizes the data, and displays the result on a webpage. Additionally, users can download the processed CSV file.

## Table of Contents

- [Project Overview](#project-overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Setup and Installation](#setup-and-installation)
- [Usage](#usage)
- [How It Works](#how-it-works)
- [Middleware and Validation](#middleware-and-validation)
- [CSV Download Feature](#csv-download-feature)
- [Best Practices and Design Considerations](#best-practices-and-design-considerations)

## Project Overview

This project was created to demonstrate skills in backend development using **Go**. The goal was to build an application that can accept CSV files, process them, sort product data according to specific warehouse bay and shelf logic, and present the summarized data to the user in both a web interface and downloadable CSV file format.

The warehouse is structured into 52 bays, labelled from A to AZ, and each bay contains 10 shelves. Products need to be sorted first by bay and then by shelf, and any products located on the same bay and shelf should have their quantities summed.

## Features

- **CSV File Upload:** Users can upload a CSV file that contains product data with product codes, quantities, and pick locations.
- **Data Sorting:** The backend processes the uploaded CSV, sorting products by bay and shelf.
- **Data Summarization:** Products with identical codes, bays, and shelves are summarized (quantities are summed).
- **CSV Download:** After processing, users can download the sorted and summarized data as a new CSV file.
- **Middleware for Data Validation:** Ensures that the CSV data is properly cleaned and validated before processing, removing invalid or malformed entries.
- **Error Handling:** The application gracefully handles errors like invalid CSV format, incorrect data, and file size limitations.

## Technologies Used

- **Go:** The backend is written in Go to handle file uploads, CSV processing, and data sorting.
- **HTML & CSS:** For rendering the upload form and displaying the processed results in a table on the frontend.
- **Gorilla Mux:** A powerful routing library used for managing HTTP routes in the Go web server.
- **CSV Package:** The Go `encoding/csv` package is used to parse and write CSV files.

## Project Structure

The project is organized as follows:

```bash
.
├── cmd/
│   └── main.go                # Main entry point of the application
├── internal/
│   └── processor.go           # Handles CSV processing logic
├── web/
│   └── templates/
│       └── index.html         # HTML template for frontend
├── go.mod                     # Go module configuration
├── go.sum                     # Go module checksum
├── README.md                  # Project documentation
```

- **cmd/main.go**: Contains the main application logic, setting up HTTP routes and managing the file upload.

- **internal/processor.go**: Handles CSV file reading, data sorting, summarization, and validation.

- **web/templates/index.html**: Frontend template for user interaction, including CSV upload and displaying results.

## Setup and Installation

To set up this project locally, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/richardwaters9049/portal-fullstack.git
```

2. Navigate to the project directory:

```bash
cd protal-fullstack
```

3. Install dependencies:

```bash
go mod tidy
```

4. Build and run the project:

```bash
go build -o portal-fullstack ./cmd/main.go
./portal-fullstack

```

5. Access the application in your web browser at `http://localhost:8080`.

## Usage

# Uploading a CSV File:

Navigate to the home page.
Select a CSV file using the upload form and click "Upload and Process."

# Processing the Data:

The backend processes the file, sorting and summarizing the data based on the product bay and shelf.

# Viewing Results:

The results will be displayed in a table on the webpage.
The user can also download the processed results as a new CSV file.

## How It Works

The core functionality is powered by Go. When a user uploads a CSV file:

The CSV file is parsed using Go’s encoding/csv package.
Middleware cleans and validates the data to ensure it is well-formatted.
The products are sorted by bay and shelf.
Products with the same code, bay, and shelf have their quantities summed.
The sorted and summarized products are rendered in the frontend and can be downloaded as a CSV file.

Here is an example of the CSV input:

```mathematica
Product Code,Quantity,Pick Location
12456,10,AB 9
36389,4,AC 5
```

And an example of the processed output:

```mathematica
Product Code,Quantity,Pick Location
12456,10,AB 9
36389,4,AC 5
```

## Middleware and Validation

The application employs middleware to clean and validate the CSV records before processing. This is necessary because user-uploaded data can often be inconsistent, containing unwanted characters, malformed rows, or incorrect data.

In the processor.go file, the CleanCSVRecords function ensures that all data is correctly formatted. Invalid rows are skipped, and only well-formed rows are processed. Here’s an example of the middleware logic:

```go
func CleanCSVRecords(records [][]string) [][]string {
    var cleanedRecords [][]string
    for _, record := range records {
        var cleanedRecord []string
        for _, field := range record {
            filteredField := strings.Map(func(r rune) rune {
                if isValidCharacter(r) {
                    return r
                }
                return -1
            }, field)
            cleanedRecord = append(cleanedRecord, strings.TrimSpace(filteredField))
        }
        if len(cleanedRecord) == 3 {
            cleanedRecords = append(cleanedRecords, cleanedRecord)
        }
    }
    return cleanedRecords
}
```

## CSV Download Feature

Once the CSV data is processed and displayed on the frontend, a new feature has been added to allow users to download the processed data as a CSV file. This gives the user flexibility to store or further manipulate the data.

The backend generates a new CSV file, and a download button is provided on the frontend. The file is written dynamically based on the results of the processing.

Here’s an example of how the CSV download is implemented in Go:

```go
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
    products := // Get the processed product data
    w.Header().Set("Content-Disposition", "attachment; filename=processed_products.csv")
    w.Header().Set("Content-Type", "text/csv")
    csvWriter := csv.NewWriter(w)
    defer csvWriter.Flush()
    for _, product := range products {
        csvWriter.Write([]string{product.Code, strconv.Itoa(product.Quantity), product.Bay + " " + strconv.Itoa(product.Shelf)})
    }
}
```

## Best Practices and Design Considerations

- Separation of Concerns: The project is divided into clear sections, with routing in main.go, processing logic in processor.go, and the frontend in index.html.
- Modularity: Functions are kept small and focused, making them easier to maintain and test.
- Error Handling: The application uses proper error handling to ensure robust performance, particularly when dealing with user inputs like CSV files.
- Middleware for Data Validation: We clean and validate the uploaded CSV data before processing, ensuring that any invalid entries are filtered out, which protects the system from crashes or incorrect outputs.

## Conclusion

This project demonstrates key skills in backend web development using Go, including file handling, data processing, sorting, and summarization, while also incorporating frontend interaction with a user-friendly interface. By implementing middleware and validation, the application ensures data integrity, making it a strong example of production-ready Go code.
