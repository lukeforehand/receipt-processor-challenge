package main

import (
	"os"
	"receiptprocessor/api"
	"strconv"
)

func main() {
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	api.Serve(int(port))
}
