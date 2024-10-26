// Package server provides the handler for the loggerd web server. It is responsible for
// marshaling the request and printing it to stdout.
package server

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// Handler is an http.Handler that logs requests to the specified loggers.
type Handler struct {
	errorLogger *slog.Logger
	logger      *slog.Logger
}

// NewHandler creates a new Handler that logs requests to the provided loggers. Neither logger should be nil.
//
// Requests themselves will be logged to the logger and any errors will be logged to the errorLogger. The errorLogger
// will only be used with level slog.LevelError.
func NewHandler(errorLogger *slog.Logger, logger *slog.Logger) *Handler {
	return &Handler{
		errorLogger: errorLogger,
		logger:      logger,
	}
}

// ServeHTTP is the handler for the loggerd web server.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	attr, err := requestToAttr(req)
	if err != nil {
		h.errorLogger.Error("Failed to create request attribute", "error", err)

		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.logger.Info("Loggerd request received", attr)

	resp.WriteHeader(http.StatusNoContent)
}

// requestToAttr creates a slog.Attr from an http.Request. It will not close the request body as it assumes the caller
// will do that. It will return an error if the request body cannot be read.
func requestToAttr(req *http.Request) (slog.Attr, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return slog.Attr{}, err
	}

	return slog.Group("request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.String("proto", req.Proto),
		slog.Int("proto_major", req.ProtoMajor),
		slog.Int("proto_minor", req.ProtoMinor),
		headerToAttr(req.Header),
		slog.String("body", string(body)),
		slog.Int64("content_length", req.ContentLength),
		slog.String("host", req.Host),
		slog.String("remote_addr", req.RemoteAddr),
		slog.String("request_uri", req.RequestURI),
		slog.String("pattern", req.URL.Path),
	), nil
}

func headerToAttr(header http.Header) slog.Attr {
	attrs := make([]any, 0, len(header))

	for key, values := range header {
		joinedValues := strings.Join(values, ", ")
		attrs = append(attrs, slog.String(key, joinedValues))
	}

	return slog.Group("header", attrs...)
}
