package scrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/abdulwaheed/nethealth/storage"
	"github.com/chromedp/chromedp"
)

type UserValidationError struct {
	Message string
	Err     error
}

func (e *UserValidationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *UserValidationError) Unwrap() error {
	return e.Err
}

func NewScrapperError(message string, err error) *UserValidationError {
	return &UserValidationError{
		Message: message,
		Err:     err,
	}
}

func hasNextPage(ctx context.Context, tag string) (bool, error) {
	var nextClick string
	var found bool
	err := chromedp.Run(ctx,
		chromedp.OuterHTML(tag, &nextClick, chromedp.ByID),
		chromedp.AttributeValue(tag, "class", &nextClick, &found, chromedp.ByID),
	)
	if err != nil {
		return false, err
	}
	if nextClick != "paginate_button next disabled" {
		return true, nil
	} else {
		return false, nil
	}
}

func validateUser(ctx context.Context, user *model.User) bool {
	accountNumberStr, entity := "", ""
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('#content > section > div > div > span:nth-child(3)').textContent`, &entity),
		chromedp.Evaluate(`document.querySelector('#content > section > div > span > span:nth-child(1)').textContent`, &accountNumberStr),
	)
	if err != nil {
		return false
	}

	accountNumber := sterializeAccountNumber(accountNumberStr)
	entity = sterializeEntity(entity)
	if user.AccountNumber != accountNumber || user.Enity != entity {
		return false
	}
	return true
}

func sterializeAccountNumber(accountNumber string) int64 {
	accountNumbers := strings.SplitAfter(accountNumber, ": ")
	accountNumberInt, err := strconv.Atoi(accountNumbers[1])
	if err != nil {
		return 0
	}
	return int64(accountNumberInt)
}

func sterializeEntity(entity string) string {
	entites := strings.SplitAfter(entity, ": ")
	return entites[1]
}

func sterializeName(name string) string {
	return strings.ReplaceAll(name, " ", "")
}

func handleNoTransactions(userDataPath string) error {
	err := storage.StoreClaimsToPDF(userDataPath+"/claims/claim.pdf", []*model.Claim{})
	if err != nil {
		return err
	}

	err = storage.StoreTransactionsToPDF(userDataPath+"/transactions/transaction.pdf", []*model.Transaction{})
	if err != nil {
		return err
	}

	err = storage.StoreAgingSummaryToPDF(userDataPath+"/agingsummary/agingsummary.pdf", []*model.AgingSummary{})
	if err != nil {
		return err
	}

	err = storage.StoreTransactionDetailsToPDF(userDataPath+"/transactionbreakdowns/transactionbreakdown.pdf", []*model.TransactionDetail{})
	if err != nil {
		return err
	}

	return nil
}
