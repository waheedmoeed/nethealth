package scrapper

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

func StartScrapper(ctx context.Context, config model.Config) error {
	ctx, cancel := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()

	scrapperContext, cancel := chromedp.NewContext(ctx)
	defer cancel()

	err := login(scrapperContext, config.Email, config.Password)
	if err != nil {
		return err
	}
	fmt.Println("Login Success")

	user := &model.User{
		FirstName:     "Abner",
		LastName:      "Claire",
		AccountNumber: 8108,
		Enity:         "Ageility at Bear Creek",
		IsMigrated:    false,
	}
	err = startScrapper(scrapperContext, user, user.GetUserDataRoomPath())
	if err != nil {
		return err
	}
	return nil
}

func startScrapper(ctx context.Context, user *model.User, userDataPath string) error {
	//It is data room per user
	err := prepareDataRoomDir(userDataPath)
	if err != nil {
		return err
	}

	// err = createJobFileIfNotExists(user.GetPendingJobFilePath())
	// if err != nil {
	// 	return err
	// }

	var g errgroup.Group
	var mu sync.Mutex

	startTime := time.Now()

	g.Go(func() error {
		ledgerCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartLaggerScrapper(ledgerCTX, user, &mu, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/ledger", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running lagger scrapper: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		claimsCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartClaimsScrapper(claimsCTX, user, &mu, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/claims", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running claims scrapper: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		transactionCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartTransactionScrapper(transactionCTX, user, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/transactions", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running transaction scrapper: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		benefitCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartBenefitScrapper(benefitCTX, user, "https://p13006.therapy.nethealth.com/Financials#patient/details/28816/benefitsVerification", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running benefit scrapper: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		agingCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartAgingSummaryScrapper(agingCTX, user, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/agingSummary", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running aging summary scrapper: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		transactionDetailCTX, cancel := chromedp.NewContext(ctx)
		defer cancel()
		err := StartTransactionDetailScrapper(transactionDetailCTX, user, &mu, "https://p13006.therapy.nethealth.com/Financials#patient/details/81808/transactions", userDataPath)
		if err != nil {
			return fmt.Errorf("error while running transaction detail scrapper: %w", err)
		}
		return nil
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

// // createFileIfNotExists ensures the CSV file exists
// func createJobFileIfNotExists(filePath string) error {
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		// Create the file with a header
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		writer := csv.NewWriter(file)
// 		defer writer.Flush()

// 		// Write header row
// 		header := []string{"FileName", "FilePath", "Download", "PDFLink"}
// 		if err := writer.Write(header); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
