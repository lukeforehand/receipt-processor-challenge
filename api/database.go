package api

import (
	"fmt"
	"sync"
)

type Database struct {
	receipts sync.Map
}

func (d *Database) GetReceipt(id string) (Receipt, error) {
	value, ok := d.receipts.Load(id)
	if !ok {
		return Receipt{}, fmt.Errorf("receipt not found")
	}
	receipt, _ := value.(Receipt)
	return receipt, nil
}

func (d *Database) PutReceipt(id string, receipt Receipt) {
	d.receipts.Store(id, receipt)
}
