package scrapper

import (
	"context"

	"github.com/chromedp/chromedp"
)

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
