package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/scrapper"
	"github.com/chromedp/chromedp"
)

func main() {
	config, err := loadConfigs()
	if err != nil {
		panic("fail to get the credential files")
	}

	// err = loadUsersToLevelDB(config)
	// if err != nil {
	// 	panic("failed to load users to leveldb")
	// }

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		err = scrapper.StartPDFDownloader(context.Background(), config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		broswerOpts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		if config.Headless {
			broswerOpts = append(broswerOpts, chromedp.Flag("headless", true))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), broswerOpts...)
		defer cancel()

		opt := []chromedp.ContextOption{}
		if config.Debug {
			opt = append(opt, chromedp.WithDebugf(log.Printf))
		}

		scrapperContext, cancel := chromedp.NewContext(ctx, opt...)
		defer cancel()
		defer wg.Done()

		err = scrapper.Login(scrapperContext, "Bottwo", "TestNetHealth@1234567890")
		if err != nil {
			panic(err)
		}
		fmt.Println("Login Success for manual scrapper")

		err = scrapper.StartManualScrapper(scrapperContext, config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		broswerOpts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		if config.Headless {
			broswerOpts = append(broswerOpts, chromedp.Flag("headless", true))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), broswerOpts...)
		defer cancel()

		opt := []chromedp.ContextOption{}
		if config.Debug {
			opt = append(opt, chromedp.WithDebugf(log.Printf))
		}

		scrapperContext, cancel := chromedp.NewContext(ctx, opt...)
		defer cancel()
		defer wg.Done()

		err = scrapper.Login(scrapperContext, "Botthree", "TestNetHealth@1234567890")
		if err != nil {
			panic(err)
		}
		fmt.Println("Login Success for manual scrapper")

		err = scrapper.StartManualScrapper(scrapperContext, config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		broswerOpts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		if config.Headless {
			broswerOpts = append(broswerOpts, chromedp.Flag("headless", true))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), broswerOpts...)
		defer cancel()

		opt := []chromedp.ContextOption{}
		if config.Debug {
			opt = append(opt, chromedp.WithDebugf(log.Printf))
		}

		scrapperContext, cancel := chromedp.NewContext(ctx, opt...)
		defer cancel()
		defer wg.Done()

		err = scrapper.Login(scrapperContext, "Botfour", "TestNetHealth@1234567890")
		if err != nil {
			panic(err)
		}
		fmt.Println("Login Success for manual scrapper")

		err = scrapper.StartManualScrapper(scrapperContext, config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		broswerOpts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		if config.Headless {
			broswerOpts = append(broswerOpts, chromedp.Flag("headless", true))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), broswerOpts...)
		defer cancel()

		opt := []chromedp.ContextOption{}
		if config.Debug {
			opt = append(opt, chromedp.WithDebugf(log.Printf))
		}

		scrapperContext, cancel := chromedp.NewContext(ctx, opt...)
		defer cancel()
		defer wg.Done()

		err = scrapper.Login(scrapperContext, "Botfive", "TestNetHealth@1234567890")
		if err != nil {
			panic(err)
		}
		fmt.Println("Login Success for manual scrapper")

		err = scrapper.StartManualScrapper(scrapperContext, config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		broswerOpts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		if config.Headless {
			broswerOpts = append(broswerOpts, chromedp.Flag("headless", true))
		}
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), broswerOpts...)
		defer cancel()

		opt := []chromedp.ContextOption{}
		if config.Debug {
			opt = append(opt, chromedp.WithDebugf(log.Printf))
		}

		scrapperContext, cancel := chromedp.NewContext(ctx, opt...)
		defer cancel()
		defer wg.Done()

		err = scrapper.Login(scrapperContext, config.Email, config.Password)
		if err != nil {
			panic(err)
		}
		fmt.Println("Login Success for manual scrapper")

		err = scrapper.StartManualScrapper(scrapperContext, config)
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

func loadUsersToLevelDB(config model.Config) error {
	users, err := model.ReadUsersFromCSVFile(context.Background(), "./userscvs/current.csv", config.Entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = leveldb.PutFailedUsers(users)
	if err != nil {
		fmt.Println(err)
	}
	return err
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
