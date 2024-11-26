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

func StartAgingSummaryScrapper(ctx context.Context, agingSummaryUrl string, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/agingsummary", userDataPath)
	file, _ := os.Stat(userDataPath + "/agingsummary.pdf")
	if  file != nil &&file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}
	err := chromedp.Run(ctx,
		chromedp.Navigate(agingSummaryUrl),
		chromedp.Sleep(20*time.Second), // Adjust this time as needed
		chromedp.Navigate(agingSummaryUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}
	agingSummary, err := scrapeAgingSummary(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreAgingSummaryToPDF(userDataPath+"/agingsummary.pdf", agingSummary)
	if err != nil {
		return err
	}
	return nil
}

func scrapeAgingSummary(ctx context.Context) ([]*model.AgingSummary, error) {
	summary := make([]*model.AgingSummary, 0)
	var foundRecord string
	err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.InnerHTML("#agingSummary_tbl > tbody > tr > td:nth-child(1)", &foundRecord, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	if foundRecord == "" {
		return summary, nil
	}

	summary, err = scrapeAgingSummaryTbody(ctx)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func scrapeAgingSummaryTbody(ctx context.Context) ([]*model.AgingSummary, error) {
	records := []*model.AgingSummary{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#agingSummary_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.AgingSummary{}
		// Run chromedp tasks
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.PayerPlan, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(2)`, i), &record.Credits, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(3)`, i), &record.Total, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(4)`, i), &record.Current, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(5)`, i), &record.ThirtyDays, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(6)`, i), &record.SixtyDays, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(7)`, i), &record.NinetyDays, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(8)`, i), &record.OneTwentyDays, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(9)`, i), &record.OneEightyDays, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#agingSummary_tbl > tbody > tr:nth-child(%d) > td:nth-child(10)`, i), &record.MoreThanEigthy, chromedp.ByQuery),
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
