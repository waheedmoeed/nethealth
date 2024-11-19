package scrapper

import (
	"context"
	"fmt"
	"log"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
)

func StartScrapper(ctx context.Context, config model.Config) error {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()

	scrapperContext, cancel := chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	err := login(scrapperContext, config.Email, config.Password)
	if err != nil {
		return err
	}
	fmt.Println("Login Success")

	startScrapper(scrapperContext, config)
	return nil
}

func startScrapper(ctx context.Context, config model.Config) error {
	if config.NeedToPouplateUsers {
		err := populateUsersAndAgencies(ctx)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("fail to populate users and agencies data: %w", err)
		}
	}

	err := StartLaggerScrapper(ctx, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/ledger")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("lagger scrapper failed with error : %w", err)
	}
	return nil
}
