/*
handler_test.go contains functions for testing Receipt handlers.
*/
package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestReceiptHandlers verifies the behavior of ReceiptHandler functions.
// It checks that the receipt handlers work correctly in various scenarios.
func TestReceiptHandlers(t *testing.T) {
	assert.Equal(t, 28, GetReceiptPoints(t, `{
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
	}`))
	assert.Equal(t, 109, GetReceiptPoints(t, `{
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
	}`))
	assert.Equal(t, 15, GetReceiptPoints(t, `{
		"retailer": "Walgreens",
		"purchaseDate": "2022-01-02",
		"purchaseTime": "08:13",
		"total": "2.65",
		"items": [
			{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
			{"shortDescription": "Dasani", "price": "1.40"}
		]
	}`))
	assert.Equal(t, 31, GetReceiptPoints(t, `{
		"retailer": "Target",
		"purchaseDate": "2022-01-02",
		"purchaseTime": "13:13",
		"total": "1.25",
		"items": [
			{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
		]
	}`))
}

// GetReceiptPoints posts a Receipt json and fetches the points earned.
// receipt: a json representation of Receipt
// Returns: the points earned for the given receipt
func GetReceiptPoints(t *testing.T, receipt string) (points int) {

	handler := NewReceiptHandler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/process", strings.NewReader(receipt))
	handler.PostReceiptsProcess(recorder, request)
	body, _ := io.ReadAll(recorder.Result().Body)
	receiptId := &PostReceiptsProcessResponse{}
	json.Unmarshal(body, &receiptId)
	assert.NotEmpty(t, receiptId.Id)

	recorder = httptest.NewRecorder()
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", receiptId.Id)
	request = httptest.NewRequest(http.MethodGet, "/{id}/points", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))

	handler.GetReceiptsIdPoints(recorder, request)
	body, _ = io.ReadAll(recorder.Result().Body)
	receiptPoints := &GetReceiptsIdPointsResponse{}
	json.Unmarshal(body, &receiptPoints)

	return receiptPoints.Points

}
