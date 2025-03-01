package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// NiftyData represents the response structure from Paytm Money API
type NiftyData struct {
	Symbol        string  `json:"symbol"`
	LastPrice     float64 `json:"lastPrice"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"pChange"`
	Volume        int64   `json:"volume"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PreviousClose float64 `json:"previousClose"`
	Timestamp     string  `json:"timestamp"`
}

// Create custom HTTP client with timeout
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func fetchNifty50Data() (*NiftyData, error) {
	// Note: This is a hypothetical URL, as Paytm Money does not provide a public API
	// You'll need to replace this with the actual API endpoint from Paytm Money's documentation
	// or use a different provider that offers this data
	url := "https://api.paytmmoney.com/markets/v1/indices/nifty50"

	// Set up request with authorization headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add required headers
	req.Header.Set("Content-Type", "application/json")
	// You would normally need to add authentication headers as well
	// For example:
	// req.Header.Set("Authorization", "Bearer YOUR_API_KEY")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Check if response was successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned error: %s, body: %s", resp.Status, string(body))
	}

	// Parse response body
	var data NiftyData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &data, nil
}

func main() {
	// Set up logging
	logFile, err := os.OpenFile("nifty50_data.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Create a multi-writer to log to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("Starting Nifty 50 data fetcher...")

	// Create ticker for 1-second intervals
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Set up channel for graceful shutdown
	done := make(chan bool)

	// Create CSV file for data storage
	csvFile, err := os.Create("nifty50_history.csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer csvFile.Close()

	// Write CSV header
	_, err = csvFile.WriteString("timestamp,symbol,lastPrice,change,percentChange,volume,high,low,open,previousClose\n")
	if err != nil {
		log.Fatalf("Error writing CSV header: %v", err)
	}

	// Set up signal handling for graceful shutdown
	go func() {
		// In a real application, you would handle OS signals here
		// For simplicity, we're just using a timeout
		time.Sleep(10 * time.Minute) // Run for 10 minutes
		done <- true
	}()

	// Main loop
	for {
		select {
		case <-done:
			log.Println("Shutting down...")
			return
		case <-ticker.C:
			// Fetch data
			data, err := fetchNifty50Data()
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}

			// Log the data
			log.Printf("Nifty 50: %s Price: %.2f Change: %.2f (%.2f%%)",
				data.Symbol, data.LastPrice, data.Change, data.PercentChange)

			// Write to CSV
			csvLine := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%d,%.2f,%.2f,%.2f,%.2f\n",
				data.Timestamp, data.Symbol, data.LastPrice, data.Change, data.PercentChange,
				data.Volume, data.High, data.Low, data.Open, data.PreviousClose)

			if _, err := csvFile.WriteString(csvLine); err != nil {
				log.Printf("Error writing to CSV: %v", err)
			}
			csvFile.Sync()
		}
	}
}
