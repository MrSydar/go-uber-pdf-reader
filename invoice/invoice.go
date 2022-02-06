package invoice

import (
	"fmt"
	"io"
	"io/ioutil"
	"mrsydar/uberinvoice/util"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	invoiceNipEvidence = "NIP"
	invoiceRowEvidence = "%"

	regularInvoiceCustomerRegexStr = `^(.*)Dane do faktury są pobierane z riders.uber.com`

	cash = "karta"
	card = "gotówka"
)

var regularInvoiceCustomerRegex = regexp.MustCompile(regularInvoiceCustomerRegexStr)

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
	No            string
	Customer      string
	Nip           string
	FormattedDate string
	VatPercent    string
	Net           string
	Vat           string
	PaymentType   string
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

	var customerRegexp *regexp.Regexp

	var paymentType string
	if strings.Contains(content, cashInvoiceTypeEvidence) {
		customerRegexp = cashInvoiceCustomerRegex
		paymentType = cash
	} else {
		if strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
			customerRegexp = gocInvoiceCustomerRegex
		} else {
			customerRegexp = regularInvoiceCustomerRegex
		}
		paymentType = card
	}

	if strings.Count(content, invoiceRowEvidence) > 1 {
		return Invoice{}, fmt.Errorf("unsupported number of rows in invoice")
	}

	no, err := util.GetFirstSubgroupMatch(content, gocInvoiceNoRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract no: %v", err)
	}

	customer, err := util.GetFirstSubgroupMatch(content, customerRegexp)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract customer: %v", err)
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

	vat, err := util.GetFirstSubgroupMatch(content, gocInvoiceVatRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract gross: %v", err)
	}

	vatPercent, err := util.GetFirstSubgroupMatch(content, gocInvoiceVatPercentRegex)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't extract gross percent: %v", err)
	}

	invoiceData := Invoice{
		no,
		customer,
		nip,
		date,
		vatPercent,
		strings.Replace(net, ",", ".", 1),
		strings.Replace(vat, ",", ".", 1),
		paymentType,
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
