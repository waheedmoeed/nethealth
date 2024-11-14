package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/scrapper"
)

func main() {
	config, err := loadConfigs()
	if err != nil {
		panic("fail to get the credential files")
	}
	ctx := context.Background()
	err = scrapper.StartScrapper(ctx, config)
	fmt.Println(err)
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
