package scrapper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

func StartBenefitScrapper(ctx context.Context, user *model.User, benefitsUrl string, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/benefits", userDataPath)
	file, _ := os.Stat(userDataPath + "/benefit.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}
	err := chromedp.Run(ctx,
		chromedp.Navigate(benefitsUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}

	isValid := validateUser(ctx, user)
	if !isValid {
		return &UserValidationError{Message: "user validation failed", Err: fmt.Errorf("user validation failed at transaction for user %s, agency: %s", user.GetID(), user.Enity)}
	}

	benefits, err := scrapeBenefitTbody(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreBenefitsToPDF(userDataPath+"/benefit.pdf", benefits)
	if err != nil {
		return err
	}
	return nil
}

func StartManualBenefitScrapper(ctx context.Context, user *model.User, userDataPath string) error {
	userDataPath = fmt.Sprintf("%s/benefits", userDataPath)
	file, _ := os.Stat(userDataPath + "/benefit.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Claim file found for LaggerDataPath: %s", userDataPath)
		return nil
	}
	var benefitsVerificationUrl string
	err := chromedp.Run(ctx, chromedp.Location(&benefitsVerificationUrl))
	if err != nil {
		return err
	}
	ledgerUrl := fmt.Sprintf("%s/ledger", benefitsVerificationUrl[:strings.LastIndex(benefitsVerificationUrl, "/")])
	benefitsVerificationUrl = fmt.Sprintf("%s/benefitsVerification", benefitsVerificationUrl[:strings.LastIndex(benefitsVerificationUrl, "/")])
	err = chromedp.Run(ctx,
		chromedp.Navigate(ledgerUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
		chromedp.Navigate(benefitsVerificationUrl),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
	)
	if err != nil {
		return err
	}

	isValid := validateUser(ctx, user)
	if !isValid {
		return &UserValidationError{Message: "user validation failed", Err: fmt.Errorf("user validation failed at transaction for user %s, agency: %s", user.GetID(), user.Enity)}
	}

	benefits, err := scrapeBenefitTbody(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreBenefitsToPDF(userDataPath+"/benefit.pdf", benefits)
	if err != nil {
		return err
	}
	return nil
}

func scrapeBenefitTbody(ctx context.Context) ([]*model.Benefit, error) {
	records := []*model.Benefit{}
	// Run chromedp tasks

	numberOfRecords := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate("document.querySelectorAll('#benefitsVerification_tbl > tbody > tr').length", &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.Benefit{}
		// Run chromedp tasks
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`#benefitsVerification_tbl > tbody > tr:nth-child(%d) > td:nth-child(1)`, i), &record.Text, chromedp.ByQuery),
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
