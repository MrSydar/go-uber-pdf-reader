package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	grossObligatoryCorrectionInvoiceTypeEvidence = `Datapowstaniaobowiązkupodatkowego`
	gocInvoiceNumRegexStr                        = `Numer faktury korygującej:\s+([A-Z]+-\d{2}-\d{4}-\d{7})`
	gocInvoiceDateRegexStr                       = `\d{2} [a-z]{3} \d{4}`
	gocInvoiceNetRegexStr                        = `Wartość całkowita netto\s+(\d+\,\d+)\s+zł`
	gocInvoiceGrossRegexStr                      = `Wartość całkowita brutto\s+(\d+\,\d+)\s+zł`
	gocInvoiceNipRegexStr                        = `NIP:\s(\d{10})\sFaktura wystawiona przez`
)

type invoice struct {
	no            string
	nip           string
	formattedDate string
	net           float64
	gross         float64
}

var gocInvoiceNoRegex = regexp.MustCompile(gocInvoiceNumRegexStr)
var gocInvoiceNipRegex = regexp.MustCompile(gocInvoiceNipRegexStr)
var gocInvoiceDateRegex = regexp.MustCompile(gocInvoiceDateRegexStr)
var gocInvoiceNetRegex = regexp.MustCompile(gocInvoiceNetRegexStr)
var gocInvoiceGrossRegex = regexp.MustCompile(gocInvoiceGrossRegexStr)

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

func extractNumber(text string, re *regexp.Regexp) (float64, error) {
	valStr, err := getFirstSubgroupMatch(text, re)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(strings.Replace(valStr, ",", ".", 1), 32)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func extractNip(text string) (string, error) {
	if strings.Count(text, "NIP") >= 2 {
		return getFirstSubgroupMatch(text, gocInvoiceNipRegex)
	} else {
		return "", nil
	}
}

func extractInvoiceData(content string) (invoice, error) {
	if !strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
		return invoice{}, fmt.Errorf("unsupported invoice type for extraction")
	}

	fmt.Println(content)

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

	net, err := extractNumber(content, gocInvoiceNetRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract net: %v", err)
	}

	gross, err := extractNumber(content, gocInvoiceGrossRegex)
	if err != nil {
		return invoice{}, fmt.Errorf("can't extract gross: %v", err)
	}

	invoiceData := invoice{no, nip, date, net, gross}

	return invoiceData, nil
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

		if content, err := getFirstPageContent(filepath); err != nil {
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
