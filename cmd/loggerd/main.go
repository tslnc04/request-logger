/*
Loggerd is a web server that logs all requests to stdout. All requests are logged as JSON to stdout, separated by a
newline. Any errors or other messages are printed to stderr.

Usage:

	loggerd [flags]

The flags are:

	-h, -help
		print this help message

	-p, -port string
		Port to listen on. Defaults to 8080.
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tslnc04/request-logger/internal/server"
)

//nolint:lll
const usage = `Loggerd is a web server that logs all requests to stdout. All requests are logged as JSON to stdout, separated by a
newline. Any errors or other messages are printed to stderr.

Usage:

	loggerd [flags]

The flags are:

	-h, -help
		print this help message

	-p, -port string
		Port to listen on. Defaults to 8080.
`

var (
	help bool
	port string
)

func init() {
	const (
		helpUsage = "print this help message"
		portUsage = "port to listen on"

		defaultHelp = false
		defaultPort = ":8080"
	)

	flag.BoolVar(&help, "help", defaultHelp, helpUsage)
	flag.BoolVar(&help, "h", defaultHelp, helpUsage+" (shorthand)")

	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&port, "p", defaultPort, portUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	if help {
		fmt.Fprint(os.Stderr, usage)

		return
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	http.HandleFunc("/", server.Handle)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen on port %s: %v\n", port, err)

		os.Exit(2)
	}
}
