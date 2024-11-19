package scrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

const basePath = "#DataTables_Table_1 > tbody > tr"

func StartLaggerScrapper(ctx context.Context, laggerURL string) error {
	err := chromedp.Run(ctx,
		chromedp.Navigate(laggerURL),
		chromedp.Sleep(5*time.Second), // Adjust this time as needed
		chromedp.WaitVisible(`#cust-btn-show-all-children`, chromedp.ByID),
		chromedp.Click(`#cust-btn-show-all-children`, chromedp.ByID),
	)
	if err != nil {
		return err
	}
	laggerGroup, err := scrapeLaggerGroups(ctx)
	if err != nil {
		return err
	}
	err = storage.StoreLaggerGroupsToPDF("laggers.pdf", laggerGroup)
	if err != nil {
		return err
	}
	return nil
}

func scrapeLagger(ctx context.Context, groupPath string) ([]*model.LaggerGroup, error) {
	lagger := make([]*model.LaggerGroup, 0)
	for {
		groups, err := scrapeLaggerGroups(ctx)
		if err != nil {
			return nil, err
		}
		lagger = append(lagger, groups...)
		hasNextPage, err := hasNextPage(ctx)
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

func scrapeLaggerGroups(ctx context.Context) ([]*model.LaggerGroup, error) {
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
		group, err := scrapeLaggerGroup(ctx, fmt.Sprintf("%s:nth-child(%d)", basePath, i+1))
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
func scrapeLaggerGroup(ctx context.Context, groupPath string) (*model.LaggerGroup, error) {
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
		return nil, nil
	}
	//#DataTables_Table_1 > tbody > tr:nth-child(2) > td > table > tbody:nth-child(2)
	for i := 1; i <= numberOfRecords; i++ {
		records, err := scrapeTbody(ctx, fmt.Sprintf("%s:nth-child(%d)", groupPath, i+1))
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
func scrapeTbody(ctx context.Context, bodyPath string) ([]*model.Lagger, error) {
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
		records = append(records, record)
	}

	return records, nil
}

// func DownloadLaggerPDF(ctx context.Context, url string) (string, error) {
// 	var pdf64 string
// 	err := chromedp.Run(ctx,
// 		chromedp.Tasks{
// 			chromedp.ActionFunc(func(ctx context.Context) error {
// 				tabCtx, cancel := chromedp.NewContext(ctx, chromedp.WithNewTab())
// 				defer cancel()
// 				err = chromedp.Run(tabCtx,
// 					chromedp.Tasks{
// 						chromedp.Navigate(url),
// 						chromedp.Click(`#btnDownload`, chromedp.ByID),
// 						chromedp.Sleep(10 * time.Second),
// 						chromedp.ActionFunc(func(ctx context.Context) error {
// 							var pdf []byte
// 							err = chromedp.DownloadURL(url).WithDownloadPath(".")(&pdf)
// 							if err != nil {
// 								return err
// 							}
// 							pdf64 = base64.StdEncoding.EncodeToString(pdf)
// 							return nil
// 						}),
// 					})
// 				if err != nil {
// 					return err
// 				}
// 				return nil
// 			}),
// 		})
// 	if err != nil {
// 		return "", err
// 	}
// 	return pdf64, nil
// }
//forEach(td => {td.querySelector('a').?link): ""
//#ledgerrownum-1 > td:nth-child(3)
