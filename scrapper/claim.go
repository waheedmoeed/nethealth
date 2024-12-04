package scrapper

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

func StartClaimsScrapper(ctx context.Context, user *model.User, mu *sync.Mutex, claimsUrl string, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/claims", userDataPath)
	file, _ := os.Stat(userDataPath + "/claim.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}
	err := chromedp.Run(ctx,
		chromedp.Navigate(claimsUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}
	claims, err := scrapeClaims(ctx, user)
	if err != nil {
		return err
	}

	err = addClaimPDFDownloadJobs(user, mu, claims)
	if err != nil {
		return fmt.Errorf("failed to add claim PDF download jobs: %w", err)
	}

	err = storage.StoreClaimsToPDF(userDataPath+"/claim.pdf", claims)
	if err != nil {
		return err
	}
	return nil
}

func scrapeClaims(ctx context.Context, user *model.User) ([]*model.Claim, error) {
	claims := make([]*model.Claim, 0)

	var nextClick string
	var found bool
	err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.AttributeValue("#claims_tbl > tbody > tr > td", "class", &nextClick, &found, chromedp.ByID),
	)
	if err != nil {
		return nil, err
	}

	if nextClick == "dataTables_empty" {
		return claims, nil
	}

	for {
		claim, err := scrapeClaimTbody(ctx, user)
		if err != nil {
			return nil, err
		}
		claims = append(claims, claim...)
		hasNextPage, err := hasNextPage(ctx, "#claims_tbl_next")
		if err != nil {
			return nil, err
		}
		if !hasNextPage {
			break
		}
		err = chromedp.Run(ctx,
			chromedp.Click(`#claims_tbl_next`, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
	}
	return claims, nil
}

func scrapeClaimTbody(ctx context.Context, user *model.User) ([]*model.Claim, error) {
	records := []*model.Claim{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#claims_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.Claim{}
		// Run chromedp tasks
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.CreationDate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(2)`, i), &record.ServicesFrom, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(3)`, i), &record.ServicesThrough, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(4)`, i), &record.ClaimNumber, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(5)`, i), &record.ClaimType, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(6) > a`, i), &record.BatchNumber, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(7)`, i), &record.Entity, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(8)`, i), &record.PayingAgency, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(9)`, i), &record.PayerPlan, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(10)`, i), &record.PayerSequence, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`#claims_tbl > tbody > tr:nth-child(%d) > td:nth-child(11)`, i), &record.ClaimAmount, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`document.querySelector('#claims_tbl > tbody > tr:nth-child(%d)').getElementsByClassName("view-pdf-action")[0]?.href||""`, i), &record.PDFLink),
		)
		if err != nil {
			return nil, err
		}
		//update the Upload PDF LINK
		if record.PDFLink != "" {
			record.UploadedLink = record.GetUploadedLink(user)
		}
		records = append(records, record)
	}

	return records, nil
}

func addClaimPDFDownloadJobs(user *model.User, mu *sync.Mutex, records []*model.Claim) error {
	jobs := []*model.Job{}
	for _, record := range records {
		if record.PDFLink != "" {
			jobs = append(jobs, &model.Job{
				FileName: record.GetFileName(),
				FilePath: record.GetFilePath(user),
				Download: false,
				PDFLink:  record.PDFLink,
			})
		}
	}
	return leveldb.PutJobs(jobs)
}
