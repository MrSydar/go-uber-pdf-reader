package invoice

import "regexp"

const (
	cashInvoiceCustomerRegexStr = `Płatność gotówką\s*(.+)\s*Dane do faktury są pobierane z riders.uber.com`

	cashInvoiceTypeEvidence = "Płatność gotówką"
)

var cashInvoiceCustomerRegex = regexp.MustCompile(cashInvoiceCustomerRegexStr)
