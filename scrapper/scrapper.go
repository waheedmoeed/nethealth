package scrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

func StartScrapper(ctx context.Context, config model.Config) error {
	ctx, cancel := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()

	scrapperContext, cancel := chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	err := login(scrapperContext, config.Email, config.Password)
	if err != nil {
		return err
	}
	fmt.Println("Login Success")

	err = startScrapper(scrapperContext, "./data/AgeilityatBearCreek/abnerclaire_8108")
	if err != nil {
		return err
	}
	return nil
}

func startScrapper(ctx context.Context, userDataPath string) error {
	err := prepareDataRoomDir(userDataPath)
	if err != nil {
		return err
	}

	var g errgroup.Group

	startTime := time.Now()

	g.Go(func() error {
		ledgerCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartLaggerScrapper(ledgerCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/ledger", userDataPath)
		return fmt.Errorf("error while running lagger scrapper: %w", err)
	})

	g.Go(func() error {
		claimsCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartClaimsScrapper(claimsCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/claims", userDataPath)
		return fmt.Errorf("error while running claims scrapper: %w", err)
	})

	g.Go(func() error {
		transactionCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartTransactionScrapper(transactionCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/transactions", userDataPath)
		return fmt.Errorf("error while running transaction scrapper: %w", err)
	})

	g.Go(func() error {
		benefitCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartBenefitScrapper(benefitCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/28816/benefitsVerification", userDataPath)
		return fmt.Errorf("error while running benefit scrapper: %w", err)
	})

	g.Go(func() error {
		agingCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartAgingSummaryScrapper(agingCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/agingSummary", userDataPath)
		return fmt.Errorf("error while running aging summary scrapper: %w", err)
	})

	g.Go(func() error {
		transactionDetailCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartTransactionDetailScrapper(transactionDetailCTX, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/transactions", userDataPath)
		return fmt.Errorf("error while running transaction detail scrapper: %w", err)
	})

	if err := g.Wait(); err != nil {
		elapsedTime := time.Since(startTime)
		fmt.Printf("Scrapper finished  with error in %s\n", elapsedTime)
		return err
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Scrapper finished in %s\n", elapsedTime)
	return nil
}

func prepareDataRoomDir(userDataPath string) error {
	directories := []string{"transactions", "laggers", "claims", "agingsummary", "benefits", "transactionbreakdowns"}
	for _, dir := range directories {
		path := fmt.Sprintf("%s/%s", userDataPath, dir)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
