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

// TestReceiptValidation
func TestReceiptValidation(t *testing.T) {

	router := GetRouter()

	// missing content type
	request := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader("{}"))
	recorder := ProcessRequest(router, request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "header Content-Type has unexpected value")

	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name: "missing retailer",
			requestBody: `{
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"retailer\" is missing",
		},
		{
			name: "bad retailer",
			requestBody: `{
				"retailer": "^&()"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "string doesn't match the regular expression",
		},
		{
			name: "missing purchaseDate",
			requestBody: `{
				"retailer": "Target"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"purchaseDate\" is missing",
		},
		{
			name: "bad purchaseDate",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "0"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "string doesn't match the format \"date\"",
		}, {
			name: "missing purchaseTime",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"purchaseTime\" is missing",
		}, {
			name: "missing items",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"items\" is missing",
		}, {
			name: "empty items",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": []
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "minimum number of items is 1",
		}, {
			name: "missing item shortDescription",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{}
				]
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"shortDescription\" is missing",
		}, {
			name: "missing item price",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK"}
				]
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"price\" is missing",
		}, {
			name: "bad item price",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "abc"}
				]
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "string doesn't match the regular expression",
		}, {
			name: "missing total",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
				]
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "property \"total\" is missing",
		}, {
			name: "bad total",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "0",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
				],
				"total": "abc"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "string doesn't match the regular expression",
		}, {
			name: "bad purchaseTime",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "0",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
				],
				"total": "35.35"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid purchaseTime",
		}, {
			name: "OK",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
				],
				"total": "35.35"
			}`,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder = ProcessRequest(router, BuildRequest(tt.requestBody))
			if status := recorder.Code; status != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, status)
			}
			if body := recorder.Body.String(); !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}

}

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

// BuildRequest is a helper for wrapping a json body in a request with appropriate header
func BuildRequest(receipt string) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(receipt))
	request.Header.Set("Content-Type", "application/json")
	return request
}

// ProcessRequest is a helper for processing the request and returning a ResponseRecorder
func ProcessRequest(router chi.Router, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
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
