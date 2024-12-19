package scrapper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

const basePath = "#DataTables_Table_1 > tbody > tr"

func StartLaggerScrapper(ctx context.Context, user *model.User, laggerURL string, userDataPath string) (hasTransactions bool, err error) {
	userDataPath = fmt.Sprintf("%s/laggers", userDataPath)
	file, _ := os.Stat(userDataPath + "/lagger.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Lagger file found for LaggerDataPath: %s", laggerURL)
		return true, nil
	}

	err = chromedp.Run(ctx,
		chromedp.Navigate(laggerURL),
		chromedp.Sleep(20*time.Second),
		chromedp.Navigate(laggerURL),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`#cust-btn-show-all-children`, chromedp.ByID),
		chromedp.Click(`#cust-btn-show-all-children`, chromedp.ByID),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		return true, fmt.Errorf("failed to navigate and interact with lagger URL: %w", err)
	}

	isValid := validateUser(ctx, user)
	if !isValid {
		return true, &UserValidationError{Message: "user validation failed", Err: fmt.Errorf("user validation failed for user %s, agency: %s", user.GetID(), user.Enity)}
	}

	laggerGroup, err := scrapeLagger(ctx, user)
	if err != nil {
		return true, fmt.Errorf("failed to scrape lagger groups: %w", err)
	}

	err = addLaggerPDFDownloadJobs(user, laggerGroup)
	if err != nil {
		return true, fmt.Errorf("failed to add lagger PDF download jobs: %w", err)
	}

	fileName := fmt.Sprintf("%s/lagger.pdf", userDataPath)
	err = storage.StoreLaggerGroupsToPDF(fileName, laggerGroup)
	if err != nil {
		return true, fmt.Errorf("failed to store lagger groups to PDF: %w", err)
	}
	if len(laggerGroup) == 0 {
		return false, nil
	}
	return true, nil
}

func StartLaggerManualScrapper(ctx context.Context, user *model.User, userDataPath string) (hasTransactions bool, err error) {
	userDataPath = fmt.Sprintf("%s/laggers", userDataPath)
	file, _ := os.Stat(userDataPath + "/lagger.pdf")
	if file != nil && file.Name() != "" {
		fmt.Printf("Lagger file found for LaggerDataPath: %s", userDataPath)
		return true, nil
	}

	//#patientSearch_tbl > tbody > tr
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#cust-btn-show-all-children`, chromedp.ByID),
		chromedp.Click(`#cust-btn-show-all-children`, chromedp.ByID),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		return true, fmt.Errorf("failed to navigate and interact with lagger URL: %w", err)
	}

	isValid := validateUser(ctx, user)
	if !isValid {
		return true, &UserValidationError{Message: "user validation failed", Err: fmt.Errorf("user validation failed for user %s, agency: %s", user.GetID(), user.Enity)}
	}

	laggerGroup, err := scrapeLagger(ctx, user)
	if err != nil {
		return true, fmt.Errorf("failed to scrape lagger groups: %w", err)
	}

	err = addLaggerPDFDownloadJobs(user, laggerGroup)
	if err != nil {
		return true, fmt.Errorf("failed to add lagger PDF download jobs: %w", err)
	}

	fileName := fmt.Sprintf("%s/lagger.pdf", userDataPath)
	err = storage.StoreLaggerGroupsToPDF(fileName, laggerGroup)
	if err != nil {
		return true, fmt.Errorf("failed to store lagger groups to PDF: %w", err)
	}
	if len(laggerGroup) == 0 {
		return false, nil
	}
	return true, nil
}

func scrapeLagger(ctx context.Context, user *model.User) ([]*model.LaggerGroup, error) {
	lagger := make([]*model.LaggerGroup, 0)
	//check if there is any record at all
	var nextClick string
	var found bool
	err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.AttributeValue("#DataTables_Table_1 > tbody > tr > td", "class", &nextClick, &found, chromedp.ByID),
	)
	if err != nil {
		return nil, err
	}

	if nextClick == "dataTables_empty" {
		return lagger, nil
	}
	for {
		groups, err := scrapeLaggerGroups(ctx, user)
		if err != nil {
			return nil, err
		}
		lagger = append(lagger, groups...)
		hasNextPage, err := hasNextPage(ctx, "#DataTables_Table_1_next")
		if err != nil {
			return nil, err
		}
		if !hasNextPage {
			break
		}
		err = chromedp.Run(ctx,
			chromedp.Click(`#DataTables_Table_1_next`, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
	}
	return lagger, nil
}

func scrapeLaggerGroups(ctx context.Context, user *model.User) ([]*model.LaggerGroup, error) {
	groups := make([]*model.LaggerGroup, 0)
	numberOfGroups := 0
	// Run chromedp tasks
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.querySelectorAll('%s').length", basePath), &numberOfGroups),
	)
	if err != nil {
		return nil, err
	}
	if numberOfGroups == 0 {
		fmt.Printf("No groups found, path: %v", basePath)
		return nil, nil
	}
	//#DataTables_Table_1 > tbody > tr:nth-child(2) > td > table > tbody:nth-child(2)
	for i := 1; i <= numberOfGroups; i = i + 2 {
		group, err := scrapeLaggerGroup(ctx, user, fmt.Sprintf("%s:nth-child(%d)", basePath, i+1))
		if err != nil {
			return nil, err
		}
		serviceDate := ""
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf("%s:nth-child(%d) > td:nth-child(2) > strong", basePath, i), &serviceDate, chromedp.ByQuery),
		)
		if err != nil {
			return nil, err
		}
		group.ServiceDate = serviceDate
		groups = append(groups, group)
	}
	return groups, nil
}

// #DataTables_Table_1 > tbody > tr:nth-child(2) > td > table
func scrapeLaggerGroup(ctx context.Context, user *model.User, groupPath string) (*model.LaggerGroup, error) {
	groupPath = fmt.Sprintf("%s > td > table > tbody", groupPath)
	group := &model.LaggerGroup{
		Laggers:              make([]*model.Lagger, 0),
		EstimatedAdjustments: make([]*model.Lagger, 0),
	}
	numberOfRecords := 0
	// Run chromedp tasks
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.querySelectorAll('%s').length", groupPath), &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		fmt.Printf("No body found for this group, path: %v", groupPath)
		return nil, fmt.Errorf("no body found for this group although group is there")
	}
	//#DataTables_Table_1 > tbody > tr:nth-child(2) > td > table > tbody:nth-child(2)
	for i := 1; i <= numberOfRecords; i++ {
		records, err := scrapeTbody(ctx, user, fmt.Sprintf("%s:nth-child(%d)", groupPath, i+1))
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			continue
		}
		if records[0].Type == "Adjustment" {
			group.EstimatedAdjustments = append(group.EstimatedAdjustments, records...)
		} else {
			group.Laggers = append(group.Laggers, records...)
		}
	}

	return group, nil
}

// #DataTables_Table_1 > tbody > tr:nth-child(2) > td > table > tbody:nth-child(2)
func scrapeTbody(ctx context.Context, user *model.User, bodyPath string) ([]*model.Lagger, error) {
	records := []*model.Lagger{}
	var isDataRow bool
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(fmt.Sprintf("document.querySelector('%s > tr').attributes.rowId ? true : false", bodyPath), &isDataRow),
	)
	if err != nil {
		return nil, err
	}

	if !isDataRow {
		return records, nil
	}

	numberOfRecords := 0
	// Run chromedp tasks
	err = chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.querySelectorAll('%s > tr').length", bodyPath), &numberOfRecords),
	)
	if err != nil {
		return nil, err
	}
	if numberOfRecords == 0 {
		return records, nil
	}

	for i := 1; i <= numberOfRecords; i++ {
		record := &model.Lagger{}
		// Run chromedp tasks
		err = chromedp.Run(ctx,
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(1)`, bodyPath, i), &record.TxDate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(2)`, bodyPath, i), &record.Type, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`Array.from(document.querySelectorAll('%s > tr:nth-child(%d) > td:nth-child(3) a')).map(link => link.textContent.trim()).toString()`, bodyPath, i), &record.ControlNumber),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(4)`, bodyPath, i), &record.Description, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(5)`, bodyPath, i), &record.Seq, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(6)`, bodyPath, i), &record.ServiceDate, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(7)`, bodyPath, i), &record.Category, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(8)`, bodyPath, i), &record.DBAmount, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(9)`, bodyPath, i), &record.CRAmount, chromedp.ByQuery),
			chromedp.InnerHTML(fmt.Sprintf(`%s > tr:nth-child(%d) > td:nth-child(10)`, bodyPath, i), &record.Balance, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`document.querySelector('%s > tr:nth-child(%d)').getElementsByClassName("view-pdf-action")[0]?.href||""`, bodyPath, i), &record.PDFLink),
		)
		if err != nil {
			return nil, err
		}
		//update the Upload PDF LINK
		if record.PDFLink != "" {
			if record.Type == "Adjustment" {
				record.UploadedLink = record.GetAdjustmentUploadedLink(user)
			} else {
				record.UploadedLink = record.GetUploadedLink(user)
			}
		}
		records = append(records, record)
	}

	return records, nil
}

func addLaggerPDFDownloadJobs(user *model.User, records []*model.LaggerGroup) error {
	jobs := []*model.Job{}
	for _, record := range records {

		for _, lagger := range record.Laggers {
			if lagger.PDFLink != "" {
				jobs = append(jobs, &model.Job{
					FileName: lagger.GetFileName() + "_" + fmt.Sprintf("%d", user.AccountNumber),
					FilePath: lagger.GetFilePath(user),
					Download: false,
					PDFLink:  lagger.PDFLink,
				})
			}
		}

		for _, lagger := range record.EstimatedAdjustments {
			if lagger.PDFLink != "" {
				jobs = append(jobs, &model.Job{
					FileName: lagger.GetAdjustmentFileName() + "_" + fmt.Sprintf("%d", user.AccountNumber),
					FilePath: lagger.GetAdjustmentFilePath(user),
					Download: false,
					PDFLink:  lagger.PDFLink,
				})
			}
		}
	}
	return leveldb.PutJobs(jobs)
}

//#DataTables_Table_1 > tbody > tr:nth-child(2) > td > table > tbody
