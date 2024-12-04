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

func StoreTransactionsToPDF(fileName string, transactions []*model.Transaction) error {
	cfg := config.NewBuilder().
		WithDebug(true).
		WithLeftMargin(5).
		WithTopMargin(15).
		WithRightMargin(5).
		Build()

	m := maroto.New(cfg)
	m.AddRows(text.NewRow(10, "Transactions", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: &props.WhiteColor,
	}).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()}))

	headers := []string{"Service Date", "Service Code", "Description", "Claim Type", "Units", "Rate", "Charge", "Payer", "Batch", "Balance", "Entity"}
	widths := []float64{1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1}

	headerCols := make([]core.Col, 0)
	for i, header := range headers {
		headerCols = append(headerCols, text.NewCol(int(widths[i]), header, props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold, Bottom: 3, Top: 3}))
	}
	m.AddAutoRow(headerCols...).WithStyle(&props.Cell{BackgroundColor: getGrayColor()})

	for _, transaction := range transactions {
		recordCols := make([]core.Col, 0)
		values := []string{transaction.ServiceDate, transaction.ServiceCode, transaction.Description, transaction.ClaimType, transaction.Units, transaction.Rate, transaction.Charge, transaction.Payer, transaction.Batch, transaction.Balance, transaction.Entity}
		for i, value := range values {
			recordCols = append(recordCols, text.NewCol(int(widths[i]), value, props.Text{Size: 9, Align: align.Center, Bottom: 3, Top: 3}))
		}
		m.AddAutoRow(recordCols...)
	}

	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to save transactions PDF %v : %w", fileName, err)
	}

	err = document.Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save transaction PDF Path : %s, error: %w", fileName, err)
	}
	return nil
}
