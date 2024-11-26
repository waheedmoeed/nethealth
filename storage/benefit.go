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

func StoreBenefitsToPDF(fileName string, benefits []*model.Benefit) error {
	cfg := config.NewBuilder().
		WithDebug(true).
		WithLeftMargin(5).
		WithTopMargin(15).
		WithRightMargin(5).
		Build()

	m := maroto.New(cfg)
	m.AddRows(text.NewRow(10, "Benefits", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: &props.WhiteColor,
	}).WithStyle(&props.Cell{BackgroundColor: getDarkGrayColor()}))

	headers := []string{"Text"}
	widths := []float64{12}

	headerCols := make([]core.Col, 0)
	for i, header := range headers {
		headerCols = append(headerCols, text.NewCol(int(widths[i]), header, props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold, Bottom: 3, Top: 3}))
	}
	m.AddAutoRow(headerCols...).WithStyle(&props.Cell{BackgroundColor: getGrayColor()})

	for _, benefit := range benefits {
		recordCols := make([]core.Col, 0)
		values := []string{benefit.Text}
		for i, value := range values {
			recordCols = append(recordCols, text.NewCol(int(widths[i]), value, props.Text{Size: 9, Align: align.Center, Bottom: 3, Top: 3}))
		}
		m.AddAutoRow(recordCols...)
	}

	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate benefits PDF: %w", err)
	}

	err = document.Save(fileName)
	if err != nil {
		return fmt.Errorf("failed to save benefits PDF: %w", err)
	}
	return nil
}
