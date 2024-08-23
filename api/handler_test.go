/*
handler_test.go contains functions for testing Receipt handlers.
*/
package api

import (
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
