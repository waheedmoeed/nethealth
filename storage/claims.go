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

func StoreClaimsToPDF(fileName string, claims []*model.Claim) error {
	cfg := config.NewBuilder().
		WithDebug(true).
		WithLeftMargin(5).
		WithTopMargin(15).
		WithRightMargin(5).
		Build()

	m := maroto.New(cfg)
	m.AddRows(text.NewRow(10, "Claims", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: &props.WhiteColor,
	}).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()}))

	headers := []string{"Creation Date", "Services From", "Services Through", "Claim Number", "Claim Type", "Batch Number", "Entity", "Paying Agency", "Payer Plan", "Payer Sequence", "Claim Amount", "Claim PDF"}
	widths := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	headerCols := make([]core.Col, 0)
	for i, header := range headers {
		headerCols = append(headerCols, text.NewCol(int(widths[i]), header, props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold, Bottom: 3, Top: 3}))
	}
	m.AddAutoRow(headerCols...).WithStyle(&props.Cell{BackgroundColor: getGrayColor()})

	for _, claim := range claims {
		recordCols := make([]core.Col, 0)
		values := []string{claim.CreationDate, claim.ServicesFrom, claim.ServicesThrough, claim.ClaimNumber, claim.ClaimType, claim.BatchNumber, claim.Entity, claim.PayingAgency, claim.PayerPlan, claim.PayerSequence, claim.ClaimAmount, claim.UploadedLink}
		for i, value := range values {
			if i == len(values)-1 && value != "" {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), "Claim", props.Text{Size: 7.5, Align: align.Center, Bottom: 3, Top: 3, Hyperlink: &value}))
			} else {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), value, props.Text{Size: 7.5, Align: align.Center, Bottom: 3, Top: 3}))
			}
		}
		m.AddAutoRow(recordCols...)
	}

	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate claims PDF: %w", err)
	}

	err = document.Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}
	return nil
}
