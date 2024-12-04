package scrapper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

func StartTransactionScrapper(ctx context.Context, user *model.User, transactionsUrl string, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/transactions", userDataPath)
	file, _ := os.Stat(userDataPath + "/transaction.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}
	err := chromedp.Run(ctx,
		chromedp.Navigate(transactionsUrl),
		chromedp.Sleep(10*time.Second), // Adjust this time as needed
		chromedp.Navigate(transactionsUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}
	transactions, err := scrapeTransactions(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreTransactionsToPDF(userDataPath+"/transaction.pdf", transactions)
	if err != nil {
		return err
	}
	return nil
}

func scrapeTransactions(ctx context.Context) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0)

	var nextClick string
	var found bool
	err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.AttributeValue("#transaction_tbl > tbody > tr > td", "class", &nextClick, &found, chromedp.ByID),
	)
	if err != nil {
		return nil, err
	}

	if nextClick == "dataTables_empty" {
		return transactions, nil
	}

	for {
		transaction, err := scrapeTransactionTbody(ctx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction...)
		hasNextPage, err := hasNextPage(ctx, "#transaction_tbl_next")
		if err != nil {
			return nil, err
		}
		if !hasNextPage {
			break
		}
		err = chromedp.Run(ctx,
			chromedp.Click(`#transaction_tbl_next`, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
	}
	return transactions, nil
}

func scrapeTransactionTbody(ctx context.Context) ([]*model.Transaction, error) {
	records := []*model.Transaction{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#transaction_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.Transaction{}
		// Run chromedp tasks
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.ServiceDate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(2)`, i), &record.ServiceCode, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(3)`, i), &record.Description, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(4)`, i), &record.ClaimType, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(5)`, i), &record.Units, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(6)`, i), &record.Rate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(7)`, i), &record.Charge, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(8)`, i), &record.Payer, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(9) > a`, i), &record.Batch, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(10)`, i), &record.Balance, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transaction_tbl > tbody > tr:nth-child(%d) > td:nth-child(11)`, i), &record.Entity, chromedp.ByQuery),
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
