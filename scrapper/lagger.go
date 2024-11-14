package scrapper

const (
	laggerURL = "https://p13006.therapy.nethealth.com/Financials#patient/search"
)

//#DataTables_Table_3 > tbody > tr:nth-child(2) > td > table > tbody:nth-child(2)
// func ScrapLagger(ctx context.Context) error {
// 	return nil
// }
// func ScrapLagger(ctx context.Context) ([]*model.Lagger, error) {
// 	var laggerTable Table
// 	err := chromedp.Run(ctx,
// 		chromedp.Tasks{
// 			chromedp.ActionFunc(func(ctx context.Context) error {
// 				node, err := laggerTable.GetNode(ctx)
// 				if err != nil {
// 					return err
// 				}
// 				var laggerRows []model.LaggerRow
// 				err = node.Do(func(node *cdp.Node) error {
// 					err = node.Children(&laggerRows, func(node *cdp.Node) bool {
// 						return node.NodeName == "tr" && node.ChildrenNumber() > 0
// 					})
// 					if err != nil {
// 						return err
// 					}
// 					return nil
// 				})
// 				if err != nil {
// 					return err
// 				}
// 				laggers := make([]*model.Lagger, len(laggerRows)-1)
// 				for i, row := range laggerRows[1:] {
// 					laggers[i] = &model.Lagger{
// 						AccountNumber: row.Cells[0].Value,
// 						PatientName:   row.Cells[1].Value,
// 						Insurance:     row.Cells[2].Value,
// 						Total:         row.Cells[3].Value,
// 					}
// 				}
// 				return nil
// 			}),
// 		})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
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
