package main

import (
	"receiptprocessor/backend"
)

func main() {
	processor := backend.NewReceiptProcessor()
	processor.Start()
}
