// Package server provides the handler for the loggerd web server. It is responsible for
// marshaling the request and printing it to stdout.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// MarhalableRequest contains all of the information from the request that we want to log. As its name suggests, it can
// be used with [json.Marshal]. It is designed for use in a server rather than a client.
type MarhalableRequest struct {
	Method           string      `json:"method"`
	URL              *url.URL    `json:"url"`
	Proto            string      `json:"proto"`
	ProtoMajor       int         `json:"proto_major"`
	ProtoMinor       int         `json:"proto_minor"`
	Header           http.Header `json:"header"`
	Body             []byte      `json:"body"`
	ContentLength    int64       `json:"content_length"`
	TransferEncoding []string    `json:"transfer_encoding"`
	Host             string      `json:"host"`
	Trailer          http.Header `json:"trailer,omitempty"`
	RemoteAddr       string      `json:"remote_addr"`
	RequestURI       string      `json:"request_uri"`
	Pattern          string      `json:"pattern"`
}

// NewMarshalableRequest uses an [http.Request] to create a new MarhalableRequest. It will not close the request body as
// it assumes the caller will do that.
func NewMarshalableRequest(req *http.Request) (*MarhalableRequest, error) {
	// We do not need to close the body reader since the http server will do it for us.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return &MarhalableRequest{
		Method:           req.Method,
		URL:              req.URL,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Header:           req.Header,
		Body:             body,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Host:             req.Host,
		Trailer:          req.Trailer,
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
		Pattern:          req.URL.Path,
	}, nil
}

// Handle is the handler for the loggerd web server. Unless the marshaling fails, it will print the marshaled request to
// stdout and return a 204 No Content response. If the marshaling fails, it will return a 500 Internal Server Error
// response and print the error to stderr.
func Handle(resp http.ResponseWriter, req *http.Request) {
	marshalableReq, err := NewMarshalableRequest(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create marshalable request: %v\n", err)

		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	marshaled, err := json.Marshal(marshalableReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal request: %v\n", err)

		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	fmt.Println(string(marshaled))

	resp.WriteHeader(http.StatusNoContent)
}
