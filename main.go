package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	grossObligatoryCorrectionInvoiceTypeEvidence = `Datapowstaniaobowiązkupodatkowego`
	gocInvoiceNumRegexStr                        = `Numer faktury korygującej:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	gocInvoiceDateRegexStr                       = `\d{2} [a-z]{3} \d{4}`
	gocInvoiceNetStr                             = `Wartość całkowita netto\s+(\d+\,\d+)\s+zł`
)

type invoice struct {
	no             string
	formatted_date string
	net            float32
	gross          float32
}

var gocInvoiceNumRegex = regexp.MustCompile(gocInvoiceNumRegexStr)
var gocInvoiceDateRegex = regexp.MustCompile(gocInvoiceDateRegexStr)

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

func extractInvoiceData(content string) (invoice, error) {
	if !strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
		return invoice{}, fmt.Errorf("unsupported invoice type for extraction")
	}

	fmt.Println(content)

	invoice_no, err := getFirstSubgroupMatch(content, gocInvoiceNumRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("invoice number parse error: %v", err)
	}

	invoice_date := gocInvoiceDateRegex.FindString(content)
	if invoice_date == "" {
		return invoice{}, fmt.Errorf("invoice date parse error")
	}

	// invoice_net :=

	invoice_data := invoice{
		no:             invoice_no,
		formatted_date: invoice_date,
	}

	return invoice_data, nil
}

func main() {
	dirpath := "./invoices"

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("Error when listing files in directory: %v\n", err)
		return
	}

	for _, f := range files {
		filepath := dirpath + "/" + f.Name()

		if content, err := extract_first_pdf_page(filepath); err != nil {
			fmt.Printf("Error while processing %q file: %v\n", filepath, err)
		} else {
			invoice_data, err := extractInvoiceData(content)
			if err != nil {
				fmt.Printf("Error while processing invoice contents %q: %v\n", filepath, err)
				continue
			}

			fmt.Println(invoice_data)
		}
	}
}
