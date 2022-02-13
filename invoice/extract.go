package invoice

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	cash = "karta"
	card = "gotówka"
)

type pureInvoice interface {
	getNo() (string, error)
	getCustomer() (string, error)
	getNip() (string, error)
	getFormattedDate() (string, error)
	getVatPercent() (string, error)
	getNet() (string, error)
	getVat() (string, error)
	getPaymentType() string
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

func getPureInvoice(content string) (pureInvoice, error) {
	if strings.Contains(content, correctionInvoiceTypeEvidence) {
		return nil, fmt.Errorf("unsupported correction invoice type")
	}

	if strings.Contains(content, cashInvoiceTypeEvidence) {
		return &cashInvoice{content}, nil
	} else if strings.Contains(content, grossObligatoryCorrectionInvoiceTypeEvidence) {
		return &gocInvoice{content}, nil
	} else {
		return nil, fmt.Errorf("default invoice not supported")
	}
}

func extractInvoiceData(content string) (Invoice, error) {
	cleanContent := strings.ReplaceAll(content, "\n", " ")
	cleanContent = strings.ReplaceAll(cleanContent, "\u00a0", " ")

	pi, err := getPureInvoice(cleanContent)
	if err != nil {
		return Invoice{}, fmt.Errorf("can't process invoice: %v", err)
	}

	no, err := pi.getNo()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get no: %v", err)
	}

	customer, err := pi.getCustomer()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get customer: %v", err)
	}

	nip, err := pi.getNip()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get nip: %v", err)
	}

	date, err := pi.getFormattedDate()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get date: %v", err)
	}

	vatPercent, err := pi.getVatPercent()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get vat percent: %v", err)
	}

	net, err := pi.getNet()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get net: %v", err)
	}

	vat, err := pi.getVat()
	if err != nil {
		return Invoice{}, fmt.Errorf("can't get vat: %v", err)
	}

	paymentType := pi.getPaymentType()

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
