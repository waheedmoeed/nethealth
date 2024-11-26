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

func StartTransactionDetailScrapper(ctx context.Context, transactionDetailsUrl string, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/transactionbreakdowns", userDataPath)
	file, _ := os.Stat(userDataPath + "/transactionbreakdown.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(transactionDetailsUrl),
		chromedp.Sleep(20*time.Second), // Adjust this time as needed
		chromedp.Navigate(transactionDetailsUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
		chromedp.DoubleClick("#transaction_tbl > tbody > tr:nth-child(1)", chromedp.ByQuery),
		chromedp.Sleep(8*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}
	transactions, err := scrapeTransactionDetails(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreTransactionDetailsToPDF(userDataPath+"/transactionbreakdown.pdf", transactions)
	if err != nil {
		return err
	}
	return nil
}

func scrapeTransactionDetails(ctx context.Context) ([]*model.TransactionDetail, error) {
	transactions := make([]*model.TransactionDetail, 0)

	for {
		transaction, err := scrapeTransactionDetailsTbody(ctx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction...)
		hasNextPage, err := hasNextPage(ctx, "#trans_tbl_next")
		if err != nil {
			return nil, err
		}
		if !hasNextPage {
			break
		}
		err = chromedp.Run(ctx,
			chromedp.Click(`#trans_tbl_next`, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
	}
	return transactions, nil
}

func scrapeTransactionDetailsTbody(ctx context.Context) ([]*model.TransactionDetail, error) {
	records := []*model.TransactionDetail{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#trans_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.TransactionDetail{}
		// Run chromedp tasks
		//#trans_tbl > tbody > tr:nth-child(2) > td:nth-child(2)
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.ServiceDate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d) > td:nth-child(2)`, i), &record.ServiceCode, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d) > td:nth-child(3)`, i), &record.Description, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d) > td:nth-child(4)`, i), &record.Charge, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d) > td:nth-child(5)`, i), &record.Balance, chromedp.ByQuery),
			chromedp.Click(fmt.Sprintf(`#trans_tbl > tbody > tr:nth-child(%d)`, i), chromedp.ByQuery),
			chromedp.Sleep(2*time.Second), // Adjust this time as needed
		)
		if err != nil {
			return nil, err
		}
		//check if breakdown exists
		var hasBreakdown string
		var found bool
		err := chromedp.Run(ctx,
			chromedp.AttributeValue("#transactionDetail_tbl > tbody > tr > td", "class", &hasBreakdown, &found, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}

		if hasBreakdown != "dataTables_empty" && !found {
			breakdowns, err := scrapeTransactionDetailsBreakdown(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to parse transaction breakdown at %v, %w", record.ServiceDate, err)
			}
			record.TransactionBreakdown = breakdowns
		}
		records = append(records, record)
	}

	return records, nil
}

func scrapeTransactionDetailsBreakdown(ctx context.Context) ([]*model.TransactionBreakdown, error) {
	transactions := make([]*model.TransactionBreakdown, 0)

	for {
		transaction, err := scrapeTransactionDetailsBreakdownTbody(ctx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction...)
		hasNextPage, err := hasNextPage(ctx, "#transactionDetail_tbl_next")
		if err != nil {
			return nil, err
		}
		if !hasNextPage {
			break
		}
		err = chromedp.Run(ctx,
			chromedp.Click(`#transactionDetail_tbl_next`, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
	}
	return transactions, nil
}

func scrapeTransactionDetailsBreakdownTbody(ctx context.Context) ([]*model.TransactionBreakdown, error) {
	records := []*model.TransactionBreakdown{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#transactionDetail_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.TransactionBreakdown{}
		// Run chromedp tasks
		//#transactionDetail_tbl > tbody > tr:nth-child(6)
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.Date, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(2)`, i), &record.ResonCode, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(3)`, i), &record.Description, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(4)`, i), &record.Amount, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(5)`, i), &record.Reference, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(6)`, i), &record.Payer, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#transactionDetail_tbl > tbody > tr:nth-child(%d) > td:nth-child(7) > a`, i), &record.Batch, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`document.querySelector('#transactionDetail_tbl > tbody > tr:nth-child(%d)').getElementsByClassName("view-pdf-action")[0]?.href||""`, i), &record.PDFLink),
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
