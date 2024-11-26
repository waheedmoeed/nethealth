package storage

import (
	"fmt"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func StoreTransactionDetailsToPDF(fileName string, transactionDetails []*model.TransactionDetail) error {
	cfg := config.NewBuilder().
		WithDebug(true).
		WithLeftMargin(5).
		WithTopMargin(15).
		WithRightMargin(5).
		Build()

	m := maroto.New(cfg)

	m.AddRows(text.NewRow(10, "Transaction Breakdown", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: &props.WhiteColor,
	}).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()}))

	for _, detail := range transactionDetails {

		m.AddRows(text.NewRow(10, detail.GetDetailedName(), props.Text{
			Top:   3,
			Style: fontstyle.Bold,
			Align: align.Center,
		}).WithStyle(&props.Cell{BackgroundColor: getGrayColor()}))

		renderBreakdownLagger(m, detail.TransactionBreakdown)
		m.AddRow(5)
	}

	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to save PDF : %w", err)
	}

	err = document.Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save PDF : %w", err)
	}
	return nil
}

func renderBreakdownLagger(m core.Maroto, breakdowns []*model.TransactionBreakdown) {
	headers := []string{"Date", "Reason Code", "Description", "Amount", "Ref No", "Payer", "Batch No", "PDF Link"}
	widths := []float64{1, 2, 2, 1, 2, 2, 1, 1}

	headerCols := make([]core.Col, 0) //core.Col
	for i, header := range headers {
		headerCols = append(headerCols, text.NewCol(int(widths[i]), header, props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold, Bottom: 3, Top: 3}))
	}
	m.AddAutoRow(headerCols...)

	for _, breakdown := range breakdowns {
		recordCols := make([]core.Col, 0)
		values := []string{breakdown.Date, breakdown.ResonCode, breakdown.Description, breakdown.Amount, breakdown.Reference, breakdown.Payer, breakdown.Batch, breakdown.PDFLink}
		for i, value := range values {
			if i == len(values)-1 && value != "" {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), "Document", props.Text{Size: 7.5, Align: align.Center, Bottom: 3, Top: 3, Hyperlink: &value}))
			} else {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), value, props.Text{Size: 7.5, Align: align.Center, Bottom: 3, Top: 3}))
			}
		}
		m.AddAutoRow(recordCols...)
	}
}
