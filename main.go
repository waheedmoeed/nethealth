package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/scrapper"
)

func main() {
	config, err := loadConfigs()
	if err != nil {
		panic("fail to get the credential files")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = scrapper.StartPDFDownloader(context.Background(), config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err = scrapper.StartScrapper(context.Background(), config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	// totalRecords, err := countTotalRecords()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func countTotalRecords() {
	files, err := os.ReadDir("userscvs")
	if err != nil {
		return
	}
	var totalRecords int
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".csv") {
			continue
		}
		csvFile, err := os.Open(fmt.Sprintf("userscvs/%s", f.Name()))
		if err != nil {

		}
		defer csvFile.Close()
		csvReader := csv.NewReader(csvFile)
		records, err := csvReader.ReadAll()
		if err != nil {
		}
		totalRecords += len(records)
	}
	fmt.Printf("Total Users to scrape: %d\n", totalRecords)
	minutesPerRecord := 15
	totalMinutes := totalRecords * minutesPerRecord
	totalHours := totalMinutes / 60
	remainingMinutes := totalMinutes % 60
	totalDays := totalHours / 24
	remainingHours := totalHours % 24

	fmt.Printf("It will take approximately for one bot to work %d days, %d hours, and %d minutes to complete all records.\n", totalDays, remainingHours, remainingMinutes)
	fmt.Printf("It will take approximately %.2f days and 30 bots to work in parallel to complete all records conidering 15 min per user without downloading pdf.\n", float64(totalMinutes)/(1440*30))
}

func loadConfigs() (model.Config, error) {
	var config model.Config
	// Read JSON file
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
		return config, err
	}

	// Unmarshal the JSON into the struct

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
		return config, err
	}
	return config, nil
}
