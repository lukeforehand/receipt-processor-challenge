/*
rule.go contains methods for calculating the points earned from a Receipt
*/
package api

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Rule is a function type that calculates the points earned from a given Receipt.
//
// Example usage:
//
//	// Define a function that matches the Rule type
//	var myFunc Rule = func(r Receipt) int {
//	    // Example logic to calculate points
//	    return 100
//	}
type Rule func(r Receipt) int

// RuleProcessor uses rules to determine the points earned from a receipt
type RuleProcessor struct {
	// rules contains each Rule function for determining total points
	rules []Rule
}

// Points sums the earned points from all rules for a given Receipt
func (p *RuleProcessor) Points(receipt Receipt) int {
	points := 0
	for _, rule := range p.rules {
		points += rule(receipt)
	}
	return points
}

// NewRuleProcessor initializes a RuleProcessor, defining all the rules for
// calculating total points earned for a Receipt
func NewRuleProcessor() RuleProcessor {
	return RuleProcessor{
		rules: []Rule{
			// One point for every alphanumeric character in the retailer name.
			Rule(func(r Receipt) int {
				re := regexp.MustCompile(`[a-zA-Z0-9]`)
				return len(re.FindAllString(r.Retailer, -1))
			}),
			// 50 points if the total is a round dollar amount with no cents.
			Rule(func(r Receipt) int {
				amount, _ := strconv.ParseFloat(r.Total, 64)
				cents := int(amount * 100)
				if cents%100 == 0 {
					return 50
				}
				return 0
			}),
			// 25 points if the total is a multiple of `0.25`.
			Rule(func(r Receipt) int {
				amount, _ := strconv.ParseFloat(r.Total, 64)
				cents := int(amount * 100)
				if cents%25 == 0 {
					return 25
				}
				return 0
			}),
			// 5 points for every two items on the receipt.
			Rule(func(r Receipt) int {
				return (len(r.Items) / 2) * 5
			}),
			// If the trimmed length of the item description is a multiple of 3,
			// multiply the price by `0.2` and round up to the nearest integer.
			Rule(func(r Receipt) int {
				points := 0
				for _, item := range r.Items {
					if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
						price, _ := strconv.ParseFloat(item.Price, 64)
						points += int(math.Ceil(price * 0.2))
					}
				}
				return points
			}),
			// 6 points if the day in the purchase date is odd.
			Rule(func(r Receipt) int {
				if r.PurchaseDate.Day()%2 != 0 {
					return 6
				}
				return 0
			}),
			// 10 points if the time of purchase is after 14:00 and before 16:00.
			Rule(func(r Receipt) int {
				start := time.Date(0, time.January, 1, 14, 0, 0, 0, time.UTC)
				end := time.Date(0, time.January, 1, 16, 0, 0, 0, time.UTC)
				purchaseTime, _ := time.Parse("15:04", r.PurchaseTime)
				if purchaseTime.After(start) && purchaseTime.Before(end) {
					return 10
				}
				return 0
			}),
		},
	}
}
