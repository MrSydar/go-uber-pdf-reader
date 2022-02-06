package invoice

import "regexp"

const (
	grossObligatoryCorrectionInvoiceTypeEvidence = "Datapowstaniaobowiązkupodatkowego"
	gocInvoiceNumRegexStr                        = `Numer faktury korygującej:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	gocInvoiceDateRegexStr                       = `\d{2} [a-z]{3} \d{4}`
	gocInvoiceNetRegexStr                        = `Wartość całkowita netto\s+(\d+\,\d+)\s+zł`
	gocInvoiceVatRegexStr                        = `Wartość całkowita VAT\s+(\d+\,\d+)\s+zł`
	gocInvoiceNipRegexStr                        = `NIP:\s(\d{10})\sFaktura wystawiona przez`
	gocInvoiceVatPercentRegexStr                 = `(\d+)%`
)

var gocInvoiceNoRegex = regexp.MustCompile(gocInvoiceNumRegexStr)
var gocInvoiceNipRegex = regexp.MustCompile(gocInvoiceNipRegexStr)
var gocInvoiceDateRegex = regexp.MustCompile(gocInvoiceDateRegexStr)
var gocInvoiceNetRegex = regexp.MustCompile(gocInvoiceNetRegexStr)
var gocInvoiceVatRegex = regexp.MustCompile(gocInvoiceVatRegexStr)
var gocInvoiceVatPercentRegex = regexp.MustCompile(gocInvoiceVatPercentRegexStr)
