package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"portal-fullstack/internal"

	"github.com/gorilla/mux"
)

var tpl = template.Must(template.ParseFiles("../web/templates/index.html"))

func main() {
	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
	r.HandleFunc("/download", DownloadHandler).Methods("GET")

	// Start the web server
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// HomeHandler renders the upload form
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// UploadHandler handles the uploaded CSV file
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve the file
	err := r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Process the CSV file
	products, err := internal.ReadCSVFromReader(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Sort and summarize the products
	sortedProducts := internal.SortAndSummarise(products)

	// Save the processed CSV file
	err = internal.GenerateCSVFile("sorted_products.csv", sortedProducts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating CSV: %v", err), http.StatusInternalServerError)
		return
	}

	// Render the result in the template
	data := struct {
		Products []internal.Product
	}{
		Products: sortedProducts,
	}
	tpl.Execute(w, data)
}

// DownloadHandler serves the downloadable CSV file
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "sorted_products.csv"

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers and serve the file for download
	w.Header().Set("Content-Disposition", "attachment; filename=sorted_products.csv")
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, r, filePath)
}
