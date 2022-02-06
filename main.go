package main

import (
	"encoding/csv"
	"io"
	"log"
	"mrsydar/uberinvoice/invoice"
	"os"
)

func main() {
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
		log.SetOutput(io.Writer(errFile))
	}

	dirPath := "/home/mrsydar/Desktop/warsaw 2 styczen/korekty"

	iw.Write([]string{"no", "nip", "date", "gross_percent", "net", "gross"})
	if nextInvoice, err := invoice.GetInvoices(dirPath); err != nil {
		log.Printf("error reading invoice directory: %v\n", err)
	} else {
		for inv, err := nextInvoice(); err != io.EOF; inv, err = nextInvoice() {
			if err != nil {
				log.Printf("error extracting invoice: %v\n", err)
			} else {
				fields := inv.GetAllFields()
				iw.Write(fields)
			}
		}
	}
}
