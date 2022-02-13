package invoice

import (
	"fmt"
	"mrsydar/uberinvoice/util"
	"regexp"
	"strings"
)

const (
	grossObligatoryCorrectionInvoiceTypeEvidence = "Datapowstaniaobowiązkupodatkowego"

	gocInvoiceNoRegexStr         = `Numer faktury korygującej:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	gocInvoiceCustomerStr        = `Faktura VAT\s*(.*)\s*Faktura wystawiona przez Uber Poland`
	gocInvoiceNipRegexStr        = `NIP:\s(\d{10})\sFaktura wystawiona przez`
	gocInvoiceDateRegexStr       = `\d{2} [a-z]{3} \d{4}`
	gocInvoiceVatPercentRegexStr = `(\d+)%`
	gocInvoiceNetRegexStr        = `Wartość całkowita netto\s+(\d+\,\d+)\s+zł`
	gocInvoiceVatRegexStr        = `Wartość całkowita VAT\s+(\d+\,\d+)\s+zł`
)

var gocInvoiceNoRegex = regexp.MustCompile(gocInvoiceNoRegexStr)
var gocInvoiceCustomerRegex = regexp.MustCompile(gocInvoiceCustomerStr)
var gocInvoiceNipRegex = regexp.MustCompile(gocInvoiceNipRegexStr)
var gocInvoiceDateRegex = regexp.MustCompile(gocInvoiceDateRegexStr)
var gocInvoiceNetRegex = regexp.MustCompile(gocInvoiceNetRegexStr)
var gocInvoiceVatRegex = regexp.MustCompile(gocInvoiceVatRegexStr)
var gocInvoiceVatPercentRegex = regexp.MustCompile(gocInvoiceVatPercentRegexStr)

var gocMonthNumber = map[string]string{
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

type gocInvoice struct {
	content string
}

func (invoice *gocInvoice) getNo() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceNoRegex)
}

func (invoice *gocInvoice) getCustomer() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceCustomerRegex)
}

func (invoice *gocInvoice) getNip() (string, error) {
	if strings.Count(invoice.content, defaultInvoiceNipEvidence) >= 2 {
		return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceNipRegex)
	} else {
		return "", nil
	}
}

func (invoice *gocInvoice) getFormattedDate() (string, error) {
	extractedDate := gocInvoiceDateRegex.FindString(invoice.content)
	monthName := monthNameInDateRegex.FindString(extractedDate)
	if gocMonthNumber[monthName] == "" {
		return "", fmt.Errorf("unexpected month name: %v", monthName)
	}
	date := strings.Replace(extractedDate, monthName, gocMonthNumber[monthName], 1)
	return date[6:10] + date[3:5] + date[0:2] + "000000", nil
}

func (invoice *gocInvoice) getVatPercent() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceVatPercentRegex)
}

func (invoice *gocInvoice) getNet() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceNetRegex)
}

func (invoice *gocInvoice) getVat() (string, error) {
	return util.GetFirstSubgroupMatch(invoice.content, gocInvoiceVatRegex)
}

func (invoice *gocInvoice) getPaymentType() string {
	return cash
}
