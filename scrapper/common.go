package scrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/abdulwaheed/nethealth/model"
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
	name, accountNumberStr, entity := "", "", ""
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('#content > section > div > div > span.section-text').textContent`, &name),
		chromedp.Evaluate(`document.querySelector('#content > section > div > div > span:nth-child(3)').textContent`, &entity),
		chromedp.Evaluate(`document.querySelector('#content > section > div > span > span:nth-child(1)').textContent`, &accountNumberStr),
	)
	if err != nil {
		return false
	}

	accountNumber := sterializeAccountNumber(accountNumberStr)
	entity = sterializeEntity(entity)
	fullName := user.LastName + ", " + user.FirstName
	if name != fullName || user.AccountNumber != accountNumber || user.Enity != entity {
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