package net

import "encoding/json"

// APNSRequest holds a request data
type APNSRequest struct {
	APS json.RawMessage `json:"aps"`
}

// APNSResponse holds a APNS response message
type APNSResponse struct {
	Reason string `json:"reason"`
}

// ResponseMessage holds a successful response body content
type ResponseMessage struct {
	Message string `json:"message"`
}

// ErrorMessage holds an unsuccessful response body content
type ErrorMessage struct {
	Error string `json:"error"`
}
