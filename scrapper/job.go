package scrapper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/abdulwaheed/nethealth/leveldb"
	"github.com/abdulwaheed/nethealth/model"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func StartPDFDownloader(ctx context.Context, config model.Config) error {
	scrapperContext, cancel := chromedp.NewContext(ctx)
	defer cancel()
	err := login(scrapperContext, config.DownloaderUser, config.DownloaderPassword)
	if err != nil {
		return err
	}
	fmt.Println("Login Success for downloading pdf")

	authKey, err := getCookiesKeys(scrapperContext)
	if err != nil {
		return errors.Wrap(err, "failed to get cookies keys")
	}

	startTime := time.Now()
	executeJob(ctx, authKey)
	if err != nil {
		return err
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Scrapper finished in %s\n", elapsedTime)
	return nil
}

func executeJob(ctx context.Context, authKey string) {
	for {
		jobs, err := leveldb.GetJobs()
		if err != nil {
			fmt.Println("Error getting jobs: ", err)
		}
		var wg sync.WaitGroup
		if len(jobs) == 0 {
			time.Sleep(time.Second * 2)
			continue
		}

		maxSize := 20

		if len(jobs) < maxSize {
			maxSize = len(jobs)
		}
		for i := 0; i < maxSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := downloadAndSavePDF(jobs[i].PDFLink, jobs[i].FilePath, authKey)
				if err != nil {
					fmt.Println("Error downloading file: ", jobs[0].FileName, " Error: ", err)
				}
				err = leveldb.DeleteJob(jobs[i].FileName)
				if err != nil {
					fmt.Println("Error deleting job: ", jobs[0].FileName, " Error: ", err)
				}
			}()
		}
		wg.Wait()
	}
}

func downloadAndSavePDF(url string, filePath string, authKey string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	client := &http.Client{}
	if err != nil {
		fmt.Println(err)
	}
	key := ".ASPXAUTH=" + authKey
	req.Header.Add("Cookie", key)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	f, err := os.Create(filePath + ".pdf")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func getCookiesKeys(ctx context.Context) (string, error) {
	var authKey string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}

			for _, v := range cookies {
				if v.Name == ".ASPXAUTH" {
					authKey = v.Value
					break
				}
			}
			return nil
		}),
	)
	if err != nil {
		return "", err
	}

	return authKey, nil
}