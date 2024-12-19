package scrapper

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

var mu = sync.Mutex{}

func StartManualScrapper(scrapperContext context.Context, config model.Config) error {
	err := startScrapperForFailedUsers(scrapperContext)
	if err != nil {
		return err
	}
	return nil
}

func startScrapperForFailedUsers(ctx context.Context) error {
	userChan := make(chan *model.User)
	var g errgroup.Group
	g.Go(func() error {
		for {
			mu.Lock()
			users, err := leveldb.GetFailedUsers()
			if err != nil {
				fmt.Println("Info while getting failed users: ", err)
			}
			if len(users) == 0 {
				fmt.Println("No more failed users to pull for manual processing")
				time.Sleep(time.Second * 2)
			} else {
				fmt.Printf("Pulled user %d \n", len(users))
			}
			mu.Unlock()
			for _, user := range users {
				userChan <- user
			}
		}
	})

	for i := 0; i < 1; i++ {
		g.Go(func() error {
			for user := range userChan {
				for {
					userCTX, cancel := chromedp.NewContext(ctx)
					err := chromedp.Run(userCTX,
						chromedp.Navigate("https://p13006.therapy.nethealth.com/Financials#patient/search"),
						chromedp.Sleep(10*time.Second),
						chromedp.WaitVisible(`#collapse-menu > form > div:nth-child(2) > div:nth-child(1) > select`, chromedp.ByQueryAll),
						chromedp.EvaluateAsDevTools(`document.querySelectorAll('#collapse-menu > form > div:nth-child(2) > div:nth-child(1) > select')[0].value = 'customerNo';`, nil),
						chromedp.EvaluateAsDevTools(`document.querySelectorAll('#collapse-menu > form > div:nth-child(2) > div:nth-child(1) > select')[0].dispatchEvent(new Event('change'))`, nil),
						chromedp.Sleep(time.Second),
						chromedp.SendKeys(`#collapse-menu > form > div:nth-child(2) > div:nth-child(2) > input`, strconv.FormatInt(user.AccountNumber, 10), chromedp.ByQuery),
						chromedp.DoubleClick(`#btnSearchPatients`, chromedp.ByID),
						chromedp.WaitVisible("#patientSearch_tbl > tbody > tr > td:nth-child(6) > div > a > i", chromedp.ByQuery),
						chromedp.Click("#patientSearch_tbl > tbody > tr > td:nth-child(6) > div > a > i", chromedp.ByQuery),
						chromedp.Sleep(5*time.Second),
					)
					if err != nil {
						fmt.Printf("failed in searching the user: %v ", err)
						time.Sleep(2 * time.Second)
						continue
					}

					err = startManualScrapperPerUser(userCTX, user, user.GetUserDataRoomPath())
					if err == nil {
						cancel()
						err = leveldb.DeleteFailedUser(user)
						if err != nil {
							fmt.Println("Error while putting failed user: ", err)
							continue
						}
						break
					}
					fmt.Printf(" while running scrapper for failed user %v. Error: %v\n", user, err)
					time.Sleep(5 * time.Second)
				}
			}
			return nil
		})
	}

	return g.Wait()
}

func startManualScrapperPerUser(userCTX context.Context, user *model.User, userDataPath string) error {
	err := prepareDataRoomDir(userDataPath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	hasTransactions, err := StartLaggerManualScrapper(userCTX, user, userDataPath)
	if err != nil {
		return fmt.Errorf("error while running lagger scrapper for user %s: %w", user.GetID(), err)
	}
	fmt.Printf("Lagger manual scrapper finished for user %s\n", user.GetID())

	if hasTransactions {
		err = StartManualClaimsScrapper(userCTX, user, userDataPath)
		if err != nil {
			return fmt.Errorf("error while running claims scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Claims manual scrapper finished for user %s\n", user.GetID())

		err = StartManualTransactionScrapper(userCTX, user, userDataPath)
		if err != nil {
			return fmt.Errorf("error while running transaction scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Transaction manual scrapper finished for user %s\n", user.GetID())

		err = StartManualAgingSummaryScrapper(userCTX, user, userDataPath)
		if err != nil {
			return fmt.Errorf("error while running aging summary scrapper for user %s: %w", user.GetID(), err)
		}
		fmt.Printf("Aging summary manual scrapper finished for user %s\n", user.GetID())
		// err = StartManualTransactionDetailScrapper(userCTX, user, userDataPath)
		// if err != nil {
		// 	return fmt.Errorf("error while running transaction detail scrapper for user %s: %w", user.GetID(), err)
		// }
		// fmt.Printf("Transaction manual detail scrapper finished for user %s\n", user.GetID())
	} else {
		err = handleNoTransactions(userDataPath)
		if err != nil {
			return fmt.Errorf("error while handling no transactions from manual scrapper for user %s: %w", user.GetID(), err)
		}
	}

	err = StartManualBenefitScrapper(userCTX, user, userDataPath)
	if err != nil {
		return fmt.Errorf("error while running benefit scrapper for user %s: %w", user.GetID(), err)
	}
	fmt.Printf("Benefit  manual scrapper finished for user %s\n", user.GetID())

	elapsedTime := time.Since(startTime)
	fmt.Printf("Manual Scrapper finished in %s for user %s\n", elapsedTime, user.GetID())
	return nil
}

//#main-nav-back-btn > div > i
