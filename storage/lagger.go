package storage

import (
	"fmt"
	"strings"

	"github.com/abdulwaheed/nethealth/model"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func StoreLaggerGroupsToPDF(fileName string, laggerGroups []*model.LaggerGroup) error {
	cfg := config.NewBuilder().
		WithDebug(true).
		WithLeftMargin(5).
		WithTopMargin(15).
		WithRightMargin(5).
		Build()

	m := maroto.New(cfg)

	for _, group := range laggerGroups {
		m.AddRows(text.NewRow(10, group.ServiceDate, props.Text{
			Top:   3,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()}))

		m.AddRows(text.NewRow(10, "Ledger", props.Text{
			Top:   3,
			Style: fontstyle.Bold,
			Align: align.Center,
		}).WithStyle(&props.Cell{BackgroundColor: getGrayColor()}))

		renderLagger(m, group.Laggers)
		m.AddRow(10)
		m.AddRows(text.NewRow(10, "Estimated Adjustments", props.Text{
			Top:   3,
			Style: fontstyle.Bold,
			Align: align.Center,
		}).WithStyle(&props.Cell{BackgroundColor: getGrayColor()}))

		renderLagger(m, group.EstimatedAdjustments)

		m.AddRow(20)
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

func renderLagger(m core.Maroto, lagger []*model.Lagger) {
	headers := []string{"Tx Date", "Type", "Control #", "Description", "Seq", "Service Date", "Category", "Debit", "Credit", "Balance", "Claim PDF"}
	widths := []float64{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1}

	headerCols := make([]core.Col, 0) //core.Col
	for i, header := range headers {
		headerCols = append(headerCols, text.NewCol(int(widths[i]), header, props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold, Bottom: 3, Top: 3}))
	}
	m.AddAutoRow(headerCols...)

	for _, lagger := range lagger {
		recordCols := make([]core.Col, 0)
		values := []string{lagger.TxDate, lagger.Type, lagger.ControlNumber, lagger.Description, lagger.Seq, lagger.ServiceDate, lagger.Category, lagger.DBAmount, lagger.CRAmount, lagger.Balance, lagger.UploadedLink}
		for i, value := range values {
			if i == len(values)-1 && value != "" {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), "Claim", props.Text{Size: 9, Align: align.Center, Bottom: 3, Top: 3, Hyperlink: &value}))
			} else if i == 2 {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), sterializeAccountNumber(value), props.Text{Size: 9, Align: align.Center, Bottom: 3, Top: 3}))
			} else {
				recordCols = append(recordCols, text.NewCol(int(widths[i]), value, props.Text{Size: 9, Align: align.Center, Bottom: 3, Top: 3}))
			}
		}
		m.AddAutoRow(recordCols...)
	}
}

func sterializeAccountNumber(accountNumber string) string {
	return strings.Replace(accountNumber, ",", ", ", -1)
}
