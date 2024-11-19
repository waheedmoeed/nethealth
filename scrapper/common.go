package scrapper

import (
	"context"

	"github.com/chromedp/chromedp"
)

func hasNextPage(ctx context.Context) (bool, error) {
	var nextClick string
	var found bool
	err := chromedp.Run(ctx,
		chromedp.OuterHTML(`#patientSearch_tbl_next`, &nextClick, chromedp.ByID),
		chromedp.AttributeValue(`#patientSearch_tbl_next`, "class", &nextClick, &found, chromedp.ByID),
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
