package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	correctionInvoiceTypeEvidence = "Faktura Korygująca"

	grossObligatoryCorrectionInvoiceTypeEvidence = "Datapowstaniaobowiązkupodatkowego"
	gocInvoiceNumRegexStr                        = `Numer faktury korygującej:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	gocInvoiceDateRegexStr                       = `\d{2} [a-z]{3} \d{4}`
	gocInvoiceNetRegexStr                        = `Wartość całkowita netto\s+(\d+\,\d+)\s+zł`
	gocInvoiceGrossRegexStr                      = `Wartość całkowita brutto\s+(\d+\,\d+)\s+zł`
	gocInvoiceNipRegexStr                        = `NIP:\s(\d{10})\sFaktura wystawiona przez`
	gocInvoiceGrossPercentRegexStr               = `(\d+)%`

	invoiceNipEvidence = "NIP"
	invoiceRowEvidence = "%"
)

type invoice struct {
	no            string
	nip           string
	formattedDate string
	grossPercent  string
	net           string
	gross         string
}

var gocInvoiceNoRegex = regexp.MustCompile(gocInvoiceNumRegexStr)
var gocInvoiceNipRegex = regexp.MustCompile(gocInvoiceNipRegexStr)
var gocInvoiceDateRegex = regexp.MustCompile(gocInvoiceDateRegexStr)
var gocInvoiceNetRegex = regexp.MustCompile(gocInvoiceNetRegexStr)
var gocInvoiceGrossRegex = regexp.MustCompile(gocInvoiceGrossRegexStr)
var gocInvoiceGrossPercentRegex = regexp.MustCompile(gocInvoiceGrossPercentRegexStr)

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

func getFirstSubgroupMatch(text string, re *regexp.Regexp) (string, error) {
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return "", fmt.Errorf("no submatch found")
	} else {
		return match[1], nil
	}
}

func extractNip(text string) (string, error) {
	if strings.Count(text, invoiceNipEvidence) >= 2 {
		return getFirstSubgroupMatch(text, gocInvoiceNipRegex)
	} else {
		return "", nil
	}
}

func extractInvoiceData(content string) (invoice, error) {
	if strings.Contains(content, correctionInvoiceTypeEvidence) {
		return invoice{}, fmt.Errorf("unsupported correction invoice type")
	}

	if !strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
		return invoice{}, fmt.Errorf("unsupported invoice type for extraction")
	}

	if strings.Count(content, invoiceRowEvidence) > 1 {
		return invoice{}, fmt.Errorf("unsupported number of rows in invoice")
	}

	no, err := getFirstSubgroupMatch(content, gocInvoiceNoRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract no: %v", err)
	}

	nip, err := extractNip(content)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract nip: %v", err)
	}

	date := gocInvoiceDateRegex.FindString(content)
	if date == "" {
		return invoice{}, fmt.Errorf("can't extract date")
	}

	net, err := getFirstSubgroupMatch(content, gocInvoiceNetRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract net: %v", err)
	}

	gross, err := getFirstSubgroupMatch(content, gocInvoiceGrossRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract gross: %v", err)
	}

	grossPercent, err := getFirstSubgroupMatch(content, gocInvoiceGrossPercentRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract gross percent: %v", err)
	}

	invoiceData := invoice{no, nip, date, grossPercent, net, gross}

	return invoiceData, nil
}

func getInvoices(dirPath string) (func() (invoice, error), error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	i := 0
	return func() (invoice, error) {
		if i == len(files) {
			return invoice{}, io.EOF
		} else {
			filepath := dirPath + "/" + files[i].Name()
			i++
			if content, err := getFirstPageContent(filepath); err != nil {
				return invoice{}, fmt.Errorf("can't read %q file contents: %v", filepath, err)
			} else {
				invoice_data, err := extractInvoiceData(content)
				if err != nil {
					return invoice{}, fmt.Errorf("can't extract invoice from %q file: %v", filepath, err)
				}
				return invoice_data, nil
			}
		}
	}, nil
}

func main() {
	dirPath := "/home/mrsydar/Desktop/warsaw 2 styczen/korekty"

	if nextInvoice, err := getInvoices(dirPath); err != nil {
		fmt.Printf("error reading invoice directory: %v\n", err)
	} else {
		for inv, err := nextInvoice(); err != io.EOF; inv, err = nextInvoice() {
			if err != nil {
				fmt.Printf("error extracting invoice: %v\n", err)
			} else {
				fmt.Println(inv)
			}
		}
	}
}
