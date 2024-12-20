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

var Token = ""

func StartPDFDownloader(ctx context.Context, config model.Config) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("DownloaderScrapper panicked: %v\n", r)
		}
	}()

	err := loginAndSetAuthKey(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to login and set auth key: %w", err)
	}

	go autoUpdateAuthToken(ctx, config)
	startTime := time.Now()
	executeJob(ctx, config)
	elapsedTime := time.Since(startTime)
	fmt.Printf("DownloaderScrapper finished in %s\n", elapsedTime)
	return nil
}

func loginAndSetAuthKey(basectx context.Context, config model.Config) error {
	ctx, cancel := chromedp.NewExecAllocator(basectx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", true))...)
	defer cancel()

	scrapperContext, cancel := chromedp.NewContext(ctx)
	defer cancel()

	err := Login(scrapperContext, config.DownloaderUser, config.DownloaderPassword)
	if err != nil {
		return err
	}
	fmt.Println("Login Success for downloading pdf")

	authKey, err := getCookiesKeys(scrapperContext)
	if err != nil {
		return errors.Wrap(err, "failed to get cookies keys")
	}
	Token = authKey
	return nil
}

func executeJob(ctx context.Context, config model.Config) {
	for {
		jobs, err := leveldb.GetJobs()
		if err != nil {
			fmt.Println("Error getting jobs: ", err)
		}
		var wg sync.WaitGroup
		if len(jobs) == 0 {
			fmt.Println("No downloading jobs found")
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
				err := downloadAndSavePDF(ctx, config, jobs[i].PDFLink, jobs[i].FilePath)
				if err != nil {
					fmt.Println("Error downloading file: ", jobs[0].FilePath, " Error: ", err)
					return
				}
				err = leveldb.DeleteJob(jobs[i].FileName)
				if err != nil {
					fmt.Println("Error deleting job: ", jobs[0].FilePath, " Error: ", err)
				}
			}()
		}
		wg.Wait()
	}
}

func downloadAndSavePDF(ctx context.Context, config model.Config, url string, filePath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	client := &http.Client{}
	if err != nil {
		fmt.Println(err)
	}
	key := ".ASPXAUTH=" + Token
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
		if res.StatusCode == http.StatusUnauthorized {
			fmt.Println("Unauthorized")
		}
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

func autoUpdateAuthToken(ctx context.Context, config model.Config) {
	time.Sleep(time.Minute * 15)
	for {
		err := loginAndSetAuthKey(ctx, config)
		if err != nil {
			fmt.Printf("failed to relogin for downloads job and set auth key: %w", err)
		}
		time.Sleep(time.Second * 5)
	}
}
