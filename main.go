package main

import (
	"log"
	"time"

	"github.com/johnfercher/maroto/v2"

	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"

	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func main() {
	// initialize PDF creation with potrait and A4 size

	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).Build() // margin, orientation and size

	m := maroto.New(cfg)

	// Build sections of the pdf
	// 1. Header
	addHeader(m)
	// 2. Invoice number
	addInvoiceDetails(m)
	// 3. Item List
	addItemList(m)
	// 4. Footer
	addFooter(m)

	// Save the Pdf File

	document, err := m.Generate()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = document.Save("output/invoice_sample.pdf")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Invoice Sample created and Saved successfully.")
}

// Adds a header to the pdf

func addHeader(m core.Maroto) {
	m.AddRow(50,
		image.NewFromFileCol(12, "assets/tat.jpeg",
			props.Rect{
				Center:  true,
				Percent: 75,
			}),
	)

	m.AddRow(20,
		text.NewCol(12, "CoolTats", props.Text{
			Top:   10,
			Style: fontstyle.Bold,
			Size:  18,
			Align: align.Center,
		}),
	)

	m.AddRow(20,
		text.NewCol(12, "Invoice", props.Text{
			Top:   10,
			Style: fontstyle.Bold,
			Size:  18,
			Align: align.Center,
		}),
	)
}

func addInvoiceDetails(m core.Maroto) {
	m.AddRow(10,
		text.NewCol(6, "Date: "+time.Now().Format("02 Jan 2006"), props.Text{
			Size:  10,
			Align: align.Left,
		}),
		text.NewCol(6, "Invoice #1001: ", props.Text{
			Size:  10,
			Align: align.Left,
		}),
	)
	m.AddRow(40, line.NewCol(12))
}

type InvoiceItem struct {
	Item          string
	Description   string
	Quantity      string
	Price         string
	DiscountPrice string
	Total         string
}

func (o InvoiceItem) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(2, "Item", props.Text{Style: fontstyle.Bold}),
		text.NewCol(3, "Description", props.Text{Style: fontstyle.Bold}),
		text.NewCol(1, "Quantity", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Price", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "DiscountedPrice", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Total", props.Text{Style: fontstyle.Bold}),
	)
}

func (o InvoiceItem) GetContent(i int) core.Row {
	r := row.New(5).Add(
		text.NewCol(2, o.Item),
		text.NewCol(3, o.Description),
		text.NewCol(1, o.Quantity),
		text.NewCol(2, o.Price),
		text.NewCol(2, o.DiscountPrice),
		text.NewCol(2, o.Total),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		})
	}

	return r
}

func getObjects() []InvoiceItem {
	var items []InvoiceItem

	contents := [][]string{
		{"Shampoo", "Regular Shampoo", "10", "$10.00", "$10.00", "$10.00"},
		{"Conditioner", "Conditioner", "2", "$15.00", "$15.00", "$30.00"},
		{"Conditioner", "Conditioner", "1", "$15.00", "$15.00", "$15.00"},
	}

	for i := 0; i < len(contents); i++ {
		items = append(items, InvoiceItem{
			Item:          contents[i][0],
			Description:   contents[i][1],
			Quantity:      contents[i][2],
			Price:         contents[i][3],
			DiscountPrice: contents[i][4],
			Total:         contents[i][5],
		})
	}
	return items
}

// Adds a list of items to the invoice
func addItemList(m core.Maroto) {
	rows, err := list.Build[InvoiceItem](getObjects())

	if err != nil {
		log.Fatal(err.Error())
	}
	m.AddRows(rows...)
}

// Adds the footer to the invoice
func addFooter(m core.Maroto) {
	m.AddRow(15,
		text.NewCol(8, "Total Amount", props.Text{
			Size:  10,
			Align: align.Right,
			Top:   5,
			Style: fontstyle.Bold,
		}),
		text.NewCol(4, "$1100", props.Text{
			Size:  10,
			Align: align.Right,
			Top:   5,
			Style: fontstyle.Bold,
		}),
	)

	m.AddRow(40,
		signature.NewCol(6, "Authorized Signatory", props.Signature{FontFamily: fontfamily.Courier}),
		code.NewQrCol(6, "https://kandy-dev.vercel.app", props.Rect{
			Percent: 75,
			Center:  true,
		}),
	)
}
