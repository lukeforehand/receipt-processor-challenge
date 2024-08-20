/*
model.go contains models for helping unmarshal json
*/
package api

// PostReceiptsProcessResponse
// Id: UUID string associated with a Receipt
type PostReceiptsProcessResponse struct {
	Id string `json:"id"`
}

// GetReceiptsIdPointsResponse
// Points: points earned for a Receipt
type GetReceiptsIdPointsResponse struct {
	Points int `json:"points"`
}
