package invoice

import (
	"fmt"
	"io"
	"io/ioutil"
	"mrsydar/uberinvoice/util"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	invoiceNipEvidence = "NIP"
	invoiceRowEvidence = "%"
)

var monthNumber = map[string]string{
	"sty": "01",
	"lut": "02",
	"mar": "03",
	"kwi": "04",
	"maj": "05",
	"cze": "06",
	"lip": "07",
	"sie": "08",
	"wrz": "09",
	"paź": "10",
	"lis": "11",
	"gru": "12",
}

type Invoice struct {
	no            string
	nip           string
	formattedDate string
	grossPercent  string
	net           string
	gross         string
}

func (inv *Invoice) GetAllFields() []string {
	fields := []string{
		inv.no,
		inv.nip,
		inv.formattedDate,
		inv.grossPercent,
		inv.net,
		inv.gross,
	}

	return fields
}

func extractNip(text string) (string, error) {
	if strings.Count(text, invoiceNipEvidence) >= 2 {
		return util.GetFirstSubgroupMatch(text, gocInvoiceNipRegex)
	} else {
		return "", nil
	}
}

func extractAndFormatDate(text string) (string, error) {
	extractedDate := gocInvoiceDateRegex.FindString(text)
	monthName := extractedDate[3:6]
	if monthNumber[monthName] == "" {
		return "", fmt.Errorf("unexpected month name: %v", monthName)
	}
	date := strings.Replace(extractedDate, monthName, monthNumber[monthName], 1)
	return date[6:10] + date[3:5] + date[0:2] + "000000", nil
}

func extractInvoiceData(content string) (Invoice, error) {
	if strings.Contains(content, correctionInvoiceTypeEvidence) {
		return Invoice{}, fmt.Errorf("unsupported correction invoice type")
	}

	if !strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
		return Invoice{}, fmt.Errorf("unsupported invoice type for extraction")
	}

	if strings.Count(content, invoiceRowEvidence) > 1 {
		return Invoice{}, fmt.Errorf("unsupported number of rows in invoice")
	}

	no, err := util.GetFirstSubgroupMatch(content, gocInvoiceNoRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract no: %v", err)
	}

	nip, err := extractNip(content)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract nip: %v", err)
	}

	date, err := extractAndFormatDate(content)
	if date == "" {
		return Invoice{}, fmt.Errorf("can't extract date: %v", err)
	}

	net, err := util.GetFirstSubgroupMatch(content, gocInvoiceNetRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract net: %v", err)
	}

	gross, err := util.GetFirstSubgroupMatch(content, gocInvoiceGrossRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract gross: %v", err)
	}

	grossPercent, err := util.GetFirstSubgroupMatch(content, gocInvoiceGrossPercentRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract gross percent: %v", err)
	}

	invoiceData := Invoice{
		no,
		nip,
		date,
		grossPercent,
		strings.Replace(net, ",", ".", 1),
		strings.Replace(gross, ",", ".", 1),
	}

	return invoiceData, nil
}

func getFirstPageContent(pdfPath string) (string, error) {
	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	plain_page_text, err := r.Page(1).GetPlainText(nil)
	if err != nil {
		return "", err
	}

	return plain_page_text, nil
}

func GetInvoices(dirPath string) (func() (Invoice, error), error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	i := 0
	return func() (Invoice, error) {
		if i == len(files) {
			return Invoice{}, io.EOF
		} else {
			filepath := dirPath + "/" + files[i].Name()
			i++
			if content, err := getFirstPageContent(filepath); err != nil {
				return Invoice{}, fmt.Errorf("can't read %q file contents: %v", filepath, err)
			} else {
				invoice_data, err := extractInvoiceData(content)
				if err != nil {
					return Invoice{}, fmt.Errorf("can't extract invoice from %q file: %v", filepath, err)
				}
				return invoice_data, nil
			}
		}
	}, nil
}
