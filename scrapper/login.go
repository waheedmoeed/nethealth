package scrapper

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func login(ctx context.Context, email string, password string) error {

	// Define URL of the login page
	loginURL := "https://p13006.therapy.nethealth.com/login"

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		// Open the login page
		chromedp.Navigate(loginURL),

		chromedp.Sleep(60*time.Second), // Adjust this time as needed
		// Wait for the page to load
		chromedp.WaitVisible(`#userName`, chromedp.ByQuery),

		// Enter email and password
		chromedp.SendKeys(`#userName`, email, chromedp.ByQuery),
		chromedp.SendKeys(`#container > div > div > div.login-screen > div.login-input-fields > div:nth-child(4) > input`, password, chromedp.ByQuery),
		///html/body/div/div/div/section/form/div/div/div/div[2]/div[1]/div[4]/input
		// Click the login button
		chromedp.Click(`#container > div > div > div.login-screen > div.login-controls > button`, chromedp.ByQuery),

		// Wait for navigation or a specific element that appears after login
		chromedp.Sleep(60*time.Second), // Adjust this time as needed
		// Wait for redirection to the target page by waiting for an element unique to that page
		chromedp.WaitVisible(`#s2id_facility_search`, chromedp.ByID),
	)

	return err
}
