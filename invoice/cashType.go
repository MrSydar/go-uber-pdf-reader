package invoice

import (
	"fmt"
	"mrsydar/uberinvoice/util"
	"regexp"
	"strings"
)

const (
	cashInvoiceTypeEvidence = "Płatność gotówką"

	cashInvoiceNoRegexStr         = `Numer Faktury:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	cashInvoiceCustomerRegexStr   = `Płatność gotówką\s*(.+)\s*Dane do faktury są pobierane z`
	cashInvoiceNipRegexStr        = `NIP:\s*(.+)\s*Dane do faktury są pobierane z`
	cashInvoiceDateRegexStr       = `Data faktury:\s+(\d+\s*[a-z]+\s*\d{4})`
	cashInvoiceVatPercentRegexStr = `(\d+)\s*%`
	cashInvoiceNetRegexStr        = `Razem po potrąceniu VAT\s*(\d+,\d+)\s*PLN`
	cashInvoiceVatRegexStr        = `Wartość całkowita VAT\s*\d+\s*%\s*(\d+\,\d+)\s*PLN`
)

var cashInvoiceNoRegex = regexp.MustCompile(cashInvoiceNoRegexStr)
var cashInvoiceCustomerRegex = regexp.MustCompile(cashInvoiceCustomerRegexStr)
var cashInvoiceDateRegex = regexp.MustCompile(cashInvoiceDateRegexStr)
var cashInvoiceNetRegex = regexp.MustCompile(cashInvoiceNetRegexStr)
var cashInvoiceVatRegex = regexp.MustCompile(cashInvoiceVatRegexStr)
var cashInvoiceNipRegex = regexp.MustCompile(cashInvoiceNipRegexStr)
var cashInvoiceVatPercentRegex = regexp.MustCompile(cashInvoiceVatPercentRegexStr)

var cashMonthNumber = map[string]string{
	"stycznia":     "01",
	"lutego":       "02",
	"marca":        "03",
	"kwietnia":     "04",
	"maja":         "05",
	"czerwca":      "06",
	"lipca":        "07",
	"sierpnia":     "08",
	"września":     "09",
	"października": "10",
	"lisopada":     "11",
	"grudnia":      "12",
}

type cashInvoice struct {
	content string
}

func (invoice *cashInvoice) getNo() (string, error) {
	fmt.Println(invoice.content)
	return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceNoRegex)
}

func (invoice *cashInvoice) getCustomer() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceCustomerRegex)
}

func (invoice *cashInvoice) getNip() (string, error) {
	if strings.Count(invoice.content, defaultInvoiceNipEvidence) >= 2 {
		return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceNipRegex)
	} else {
		return "", nil
	}
}

func (invoice *cashInvoice) getFormattedDate() (string, error) {
	extractedDate, err := util.GetFirstSubgroupMatch(invoice.content, cashInvoiceDateRegex)
	if err != nil {
		return "", fmt.Errorf("can't find date")
	}

	day := dayInDateRegex.FindString(extractedDate)
	monthName := monthNameInDateRegex.FindString(extractedDate)
	year := yearInDateRegex.FindString(extractedDate)

	if day == "" {
		return "", fmt.Errorf("can't find day in the extracted date: %v", extractedDate)
	}

	if monthName == "" {
		return "", fmt.Errorf("can't find month name in the extracted date: %v", extractedDate)
	}

	if year == "" {
		return "", fmt.Errorf("can't find year in the extracted date: %v", extractedDate)
	}

	if len(day) == 1 {
		day = "0" + day
	}

	if cashMonthNumber[monthName] == "" {
		return "", fmt.Errorf("unexpected month name: %v", monthName)
	}

	return fmt.Sprintf("%s.%s.%s", day, cashMonthNumber[monthName], year), nil
}

func (invoice *cashInvoice) getVatPercent() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceVatPercentRegex)
}

func (invoice *cashInvoice) getNet() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceNetRegex)
}

func (invoice *cashInvoice) getVat() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, cashInvoiceVatRegex)
}

func (invoice *cashInvoice) getPaymentType() string {
	return cash
}
