/*
Loggerd is a web server that logs all requests to stdout and optionally to Loki. All requests are logged as JSON to
stdout, separated by a newline. Any errors or other messages are printed to stderr. When logging to Loki, errors are
also sent to Loki in addition to stderr.

Usage:

	loggerd [flags]

The flags are:

	-h, -help
	        print this help message

	-l, -loki-url string
	        Loki API base URL such as http://localhost:3100. Must be provided to enable logging to Loki.

	-p, -port string
	        Port to listen on. Defaults to 8080.
*/
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	loki "github.com/tslnc04/loki-logger/pkg/client"
	lokislog "github.com/tslnc04/loki-logger/pkg/slog"
	"github.com/tslnc04/request-logger/internal/server"
)

//nolint:lll
const usage = `Loggerd is a web server that logs all requests to stdout and optionally to Loki. All requests are logged as JSON to
stdout, separated by a newline. Any errors or other messages are printed to stderr. When logging to Loki, errors are
also sent to Loki in addition to stderr.

Usage:

        loggerd [flags]

The flags are:

        -h, -help
                print this help message

        -l, -loki-url string
                Loki API base URL such as http://localhost:3100. Must be provided to enable logging to Loki.

        -p, -port string
                Port to listen on. Defaults to 8080.
`

var (
	help    bool
	port    string
	lokiURL string
)

func init() {
	const (
		helpUsage    = "print this help message"
		portUsage    = "port to listen on"
		lokiURLUsage = "loki api base url such as http://localhost:3100"

		defaultHelp    = false
		defaultPort    = ":8080"
		defaultLokiURL = ""
	)

	flag.BoolVar(&help, "help", defaultHelp, helpUsage)
	flag.BoolVar(&help, "h", defaultHelp, helpUsage+" (shorthand)")

	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&port, "p", defaultPort, portUsage+" (shorthand)")

	flag.StringVar(&lokiURL, "loki-url", defaultLokiURL, lokiURLUsage)
	flag.StringVar(&lokiURL, "l", defaultLokiURL, lokiURLUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	if help {
		fmt.Fprint(os.Stderr, usage)

		return
	}

	var lokiClient loki.Client

	if lokiURL != "" {
		lokiClient = loki.NewLokiClient(lokiURL + loki.PushPath)
	}

	errorLogger := newErrorLogger(lokiClient)

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	handler := server.NewHandler(errorLogger, newLogger(lokiClient))
	err := http.ListenAndServe(":8080", handler)

	if err != nil {
		errorLogger.Error("Failed to listen and serve", "port", port, "error", err)

		os.Exit(2)
	}
}

func newErrorLogger(lokiClient loki.Client) *slog.Logger {
	options := slog.HandlerOptions{AddSource: true, Level: slog.LevelError}
	errorTextHandler := slog.NewTextHandler(os.Stderr, &options)

	if lokiClient == nil {
		return slog.New(errorTextHandler)
	}

	newOptions := options
	lokiHandler := lokislog.NewHandler(lokiClient, &newOptions)

	return lokislog.NewJoinedLogger(lokiHandler, errorTextHandler)
}

func newLogger(lokiClient loki.Client) *slog.Logger {
	options := slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &options)

	if lokiClient == nil {
		return slog.New(jsonHandler)
	}

	newOptions := options
	lokiHandler := lokislog.NewHandler(lokiClient, &newOptions).
		WithAttrs([]slog.Attr{slog.String("service_name", "loggerd")})

	return lokislog.NewJoinedLogger(lokiHandler, jsonHandler)
}
