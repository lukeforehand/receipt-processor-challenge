/*
database.go contains methods for accessing stored receipts
*/
package api

import (
	"fmt"
	"sync"
)

// Database is a simple in memory key/value store implementation
// for id/Receipt
type Database struct {
	// receipts is a shared map of id->Receipt
	receipts sync.Map
}

// GetReceipt retrieves a Receipt from the Database
// id: the uuid string associated with a Receipt
// Returns: the Receipt for the given id
func (d *Database) GetReceipt(id string) (Receipt, error) {
	value, ok := d.receipts.Load(id)
	if !ok {
		return Receipt{}, fmt.Errorf("receipt not found")
	}
	receipt, _ := value.(Receipt)
	return receipt, nil
}

// PutReceipt stores a Receipt in the Database
// id: the uuid string associated with the given Receipt
// receipt: the Receipt to store
func (d *Database) PutReceipt(id string, receipt Receipt) {
	d.receipts.Store(id, receipt)
}
