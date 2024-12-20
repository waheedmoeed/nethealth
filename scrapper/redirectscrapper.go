package scrapper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

func StartScrapper(scrapperContext context.Context, config model.Config) error {
	users, err := model.ReadUsersFromCSVFile(context.Background(), "./userscvs/current.csv", config.Entity)
	if err != nil {
		return err
	}

	err = startScrapperForUsers(scrapperContext, users)
	if err != nil {
		return err
	}
	return nil
}

func startScrapperForUsers(ctx context.Context, users []*model.User) error {
	userChan := make(chan *model.User)
	var g errgroup.Group
	g.Go(func() error {
		latestState, _ := leveldb.GetAgencyState(users[0].Enity)
		if latestState != "" {
			for index, user := range users {
				if user.GetID() == latestState {
					users = users[index:]
				}
			}
		}
		fmt.Printf("Latest State: %s\n", latestState)

		for _, user := range users {
			userChan <- user
		}
		close(userChan)
		return nil
	})

	for i := 0; i < 2; i++ {
		g.Go(func() error {
			for user := range userChan {
				for {
					err := startScrapperPerUser(ctx, user, user.GetUserDataRoomPath())
					if err == nil {
						break
					}
					leveldb.PutAgencyState(user.Enity, user.GetID())
					var userValidationError *UserValidationError
					if errors.As(err, &userValidationError) {
						fmt.Printf(" while running scrapper for user %v. Error: %v\n", user, err)
						err = leveldb.PutFailedUser(user)
						if err == nil {
							break
						}
					}
					time.Sleep(5 * time.Second)
				}
			}
			return nil
		})
	}

	return g.Wait()
}

func startScrapperPerUser(ctx context.Context, user *model.User, userDataPath string) error {
	err := prepareDataRoomDir(userDataPath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	userCTX, cancel := chromedp.NewContext(ctx)
	defer cancel()

	hasTransactions, err := StartLaggerScrapper(userCTX, user, user.GetLedgerPageURL(), userDataPath)
	if err != nil {
		return fmt.Errorf("error while running lagger scrapper for user %s: %w", user.GetID(), err)
	}
	fmt.Printf("Lagger scrapper finished for user %s\n", user.GetID())

	if hasTransactions {

		err = StartClaimsScrapper(userCTX, user, user.GetClaimsPageURL(), userDataPath)
		if err != nil {
			return fmt.Errorf("error while running claims scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Claims scrapper finished for user %s\n", user.GetID())

		err = StartTransactionScrapper(userCTX, user, user.GetTransactionsPageURL(), userDataPath)
		if err != nil {
			return fmt.Errorf("error while running transaction scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Transaction scrapper finished for user %s\n", user.GetID())

		err = StartAgingSummaryScrapper(userCTX, user, user.GetAgingSummaryPageURL(), userDataPath)
		if err != nil {
			return fmt.Errorf("error while running aging summary scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Aging summary scrapper finished for user %s\n", user.GetID())

		err = StartTransactionDetailScrapper(userCTX, user, user.GetTransactionsPageURL(), userDataPath)
		if err != nil {
			return fmt.Errorf("error while running transaction detail scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Transaction detail scrapper finished for user %s\n", user.GetID())
	} else {
		fmt.Printf("No transactions found for user %s\n", user.GetID())
		err = handleNoTransactions(userDataPath)
		if err != nil {
			return fmt.Errorf("error while handling no transactions for user %s: %w", user.GetID(), err)
		}
	}

	err = StartBenefitScrapper(userCTX, user, user.GetBenefitsPageURL(), userDataPath)
	if err != nil {
		return fmt.Errorf("error while running benefit scrapper for user %s: %w", user.GetID(), err)
	}
	fmt.Printf("Benefit scrapper finished for user %s\n", user.GetID())

	elapsedTime := time.Since(startTime)
	fmt.Printf("Scrapper finished in %s for user %s\n", elapsedTime, user.GetID())
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

// var g errgroup.Group
// var mu sync.Mutex

// startTime := time.Now()

// g.Go(func() error {
// 	ledgerCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartLaggerScrapper(ledgerCTX, user, &mu, user.GetLedgerPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running lagger scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Lagger scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// g.Go(func() error {
// 	claimsCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartClaimsScrapper(claimsCTX, user, &mu, user.GetClaimsPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running claims scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Claims scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// g.Go(func() error {
// 	transactionCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartTransactionScrapper(transactionCTX, user, user.GetTransactionsPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running transaction scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Transaction scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// g.Go(func() error {
// 	benefitCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartBenefitScrapper(benefitCTX, user, user.GetBenefitsPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running benefit scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Benefit scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// g.Go(func() error {
// 	agingCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartAgingSummaryScrapper(agingCTX, user, user.GetAgingSummaryPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running aging summary scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Aging summary scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// g.Go(func() error {
// 	transactionDetailCTX, cancel := chromedp.NewContext(ctx)
// 	defer cancel()
// 	err := StartTransactionDetailScrapper(transactionDetailCTX, user, &mu, user.GetTransactionsPageURL(), userDataPath)
// 	if err != nil {
// 		return fmt.Errorf("error while running transaction detail scrapper for user %s: %w", user.GetID(), err)
// 	}
// 	fmt.Printf("Transaction detail scrapper finished for user %s\n", user.GetID())
// 	return nil
// })

// if err := g.Wait(); err != nil {
// 	elapsedTime := time.Since(startTime)
// 	fmt.Printf("Scrapper finished with error in %s for user %s\n", elapsedTime, user.GetID())
// 	return err
// }
