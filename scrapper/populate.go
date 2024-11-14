package scrapper

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/chromedp"
)

// populateUsersBasedOnAgency populates users based on the given agency name.
// It scrapes all the agency names and stores them in a CSV file.
// Then it scrapes all the users for each agency and stores them in a CSV file.

func populateUsersAndAgencies(ctx context.Context) error {
	agencies, err := scrapAgencyNameFromHtml(ctx)
	if err != nil {
		return err
	}
	err = storeAgenciesToCSV(agencies)
	if err != nil {
		return err
	}
	for index, agency := range agencies {
		if index < 95 {
			continue
		}
		users, err := scrapUsersFromHtml(ctx, index)
		if err != nil {
			return err
		}
		err = storeUsersToCSV(agency, users)
		if err != nil {
			return err
		}
		chromedp.Run(ctx,
			chromedp.Click(`#s2id_facility_search`, chromedp.ByID),
		)
	}

	return nil
}

func scrapAgencyNameFromHtml(ctx context.Context) ([]string, error) {
	var agencies []string

	err := chromedp.Run(ctx,
		chromedp.Click(`#s2id_facility_search`, chromedp.ByID),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('#select2-results-1 > li > div > div > div')).map(el => el.textContent)`, &agencies),
	)
	if err != nil {
		return nil, err
	}
	return agencies, nil
}

func storeAgenciesToCSV(agencies []string) error {
	file, err := os.Create("agencies.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, agency := range agencies {
		record := []string{
			agency,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

const csvDir = "userscvs"

func storeUsersToCSV(fileName string, users []*model.User) error {
	dir := path.Join(csvDir)
	file, err := os.Create(path.Join(dir, fileName+"_users.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, user := range users {
		record := []string{
			strconv.FormatInt(int64(user.ID), 10),
			user.GetID(),
			user.FirstName,
			user.LastName,
			user.MI,
			strconv.FormatInt(user.AccountNumber, 10),
			strconv.FormatBool(user.IsMigrated),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func scrapUsersFromHtml(ctx context.Context, index int) ([]*model.User, error) {
	var users []*model.User

	err := chromedp.Run(ctx,
		chromedp.Click(fmt.Sprintf(`#select2-results-1 > li:nth-child(%d)`, index+1), chromedp.ByQuery),
		chromedp.Click(`#btnSearchPatients`, chromedp.ByID),
		chromedp.Sleep(10*time.Second), // Adjust this time as needed
		chromedp.WaitVisible(`#patientSearch_tbl > tbody > tr:nth-child(1)`, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// pagination handle
	pageNumber := 1
	for {
		var pageUsers []*model.User
		err = chromedp.Run(ctx,
			chromedp.Sleep(2*time.Second), // Adjust this time as needed
			chromedp.Evaluate(`Array.from(document.querySelectorAll('#patientSearch_tbl > tbody > tr')).map((row, index) => ({
				id: index,
				uniqueIdentifier: row.cells[1].textContent.trim() + '_' + row.cells[0].textContent.trim() + '_' + row.cells[3].textContent.trim(),
				lastName: row.cells[0].textContent.trim(),
				firstName: row.cells[1].textContent.trim(),
				mi: row.cells[2].textContent.trim(),
				accountNumber: parseInt(row.cells[3].textContent.trim(), 10),
				enity: row.cells[4].textContent.trim(),
				isMigrated: false
			}))`, &pageUsers),
		)
		if err != nil {
			return nil, err
		}
		users = append(users, pageUsers...)

		var nextClick string
		var found bool
		err = chromedp.Run(ctx,
			chromedp.OuterHTML(`#patientSearch_tbl_next`, &nextClick, chromedp.ByID),
			chromedp.AttributeValue(`#patientSearch_tbl_next`, "class", &nextClick, &found, chromedp.ByID),
		)
		if err != nil {
			return nil, err
		}
		if nextClick != "paginate_button next disabled" {
			err = chromedp.Run(ctx,
				chromedp.Click(`#patientSearch_tbl_next`, chromedp.ByID),
			)
			if err != nil {
				fmt.Printf("Agency index:%d Pages:%d , Current page number %v", index, pageNumber, pageNumber)
				break
				//return nil, err
			}
		} else {
			break
		}
		pageNumber++
	}

	return users, nil
}

// chromedp.Evaluate(`Array.from(document.querySelectorAll('#patientSearch_tbl > tbody > tr')).map((row, index) => ({
// 	id: index,
// 	uniqueIdentifier: row.cells[1].textContent.trim() + '_' + row.cells[0].textContent.trim() + '_' + row.cells[3].textContent.trim(),
// 	lastName: row.cells[0].textContent.trim(),
// 	firstName: row.cells[1].textContent.trim(),
// 	mi: row.cells[2].textContent.trim(),
// 	accountNumber: parseInt(row.cells[3].textContent.trim(), 10),
// 	enity: row.cells[4].textContent.trim(),
// 	isMigrated: false
// }))`, &users),
