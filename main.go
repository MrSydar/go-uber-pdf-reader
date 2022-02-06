package main

import (
	"encoding/csv"
	"io"
	"log"
	"mrsydar/uberinvoice/invoice"
	"os"
)

const argErrorMessage = `
	Not enough arguments.
	
	First argument:
	- path to directory with uber invoices to extract
`

func main() {
	if len(os.Args) != 2 {
		log.Fatalf(argErrorMessage)
	}

	var iw *csv.Writer
	invFile, err := os.Create("invoices.csv")
	if err != nil {
		log.Fatalf("error opening invoices.csv: %v\n", err)
	} else {
		defer invFile.Close()

		iw = csv.NewWriter(invFile)
		defer iw.Flush()
	}

	errFile, err := os.Create("errors.log")
	if err != nil {
		log.Fatalf("error opening errors.log: %v", err)
	} else {
		defer errFile.Close()
	}

	iw.Write([]string{"no", "nip", "date", "vat_percent", "net", "vat"})
	if nextInvoice, err := invoice.GetInvoices(os.Args[1]); err != nil {
		log.Printf("error reading invoice directory: %v\n", err)
	} else {
		log.SetOutput(io.Writer(errFile))
		for inv, err := nextInvoice(); err != io.EOF; inv, err = nextInvoice() {
			if err != nil {
				log.Printf("error extracting invoice: %v\n", err)
			} else {
				fields := []string{
					inv.No,
					inv.Nip,
					inv.FormattedDate,
					inv.VatPercent,
					inv.Net,
					inv.Vat,
				}
				iw.Write(fields)
			}
		}
	}
}
