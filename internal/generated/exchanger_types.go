// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package generated

// CodeAndMessage defines model for CodeAndMessage.
type CodeAndMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ExchangerResult defines model for ExchangerResult.
type ExchangerResult struct {
	Err    *CodeAndMessage `json:"err,omitempty"`
	Result float32         `json:"result"`
}