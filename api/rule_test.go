/*
rule_test.go contains functions for testing RuleProcessor and Rule.
*/
package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestReceiptHandlers verifies the behavior of ReceiptHandler functions.
// It checks that the receipt handlers work correctly in various scenarios.
func TestRuleProcessor(t *testing.T) {

	tests := []struct {
		points  int
		receipt string
	}{
		{
			points: 28,
			receipt: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
					{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
					{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
					{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
					{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"}
				],
				"total": "35.35"
			}`,
		},
		{
			points: 109,
			receipt: `{
				"retailer": "M&M Corner Market",
				"purchaseDate": "2022-03-20",
				"purchaseTime": "14:33",
				"items": [
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"}
				],
				"total": "9.00"
			}`,
		},
		{
			points: 15,
			receipt: `{
				"retailer": "Walgreens",
				"purchaseDate": "2022-01-02",
				"purchaseTime": "08:13",
				"total": "2.65",
				"items": [
					{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
					{"shortDescription": "Dasani", "price": "1.40"}
				]
			}`,
		},
		{
			points: 31,
			receipt: `{
               "retailer": "Target",
               "purchaseDate": "2022-01-02",
               "purchaseTime": "13:13",
               "total": "1.25",
               "items": [
                       {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
               ]
       		}`,
		},
	}

	p := NewRuleProcessor()
	for idx, tt := range tests {
		var receipt Receipt
		_ = json.Unmarshal([]byte(tt.receipt), &receipt)
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			points := p.Points(receipt)
			if points != tt.points {
				t.Errorf("Expected %d points, got %d", tt.points, points)
			}
		})
	}
}
