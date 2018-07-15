package logging

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Based heavily on: https://gist.github.com/cespare/3985516

const (
	// ApacheFormatPattern is the Common Logging format
	ApacheFormatPattern = "%s - - [%s] \"%s\" %d %d %s %dms\n"
)

// ApacheLogRecord represents all of the data that is needed to log in the common logging formate
type ApacheLogRecord struct {
	http.ResponseWriter

	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	responseBytes         int64
	userAgent             string
	elapsedTime           time.Duration
}

// Log writes the ApacheLogRecord to the Writer in the ApacheFormatPattern
func (r *ApacheLogRecord) Log(out io.Writer) {
	timeFormatted := r.time.Format("02/Jan/2006 03:04:05 -0700")
	requestLine := fmt.Sprintf("%s %s %s", r.method, r.uri, r.protocol)
	fmt.Fprintf(out, ApacheFormatPattern, r.ip, timeFormatted, requestLine, r.status, r.responseBytes,
		r.userAgent, r.elapsedTime.Nanoseconds()/1e3)
}

// ApacheLoggingHandler is an http.Handler that wraps another http.Handler
// and logs the http requests in the common logging format
type ApacheLoggingHandler struct {
	handler http.Handler
	out     io.Writer
}

// NewApacheLoggingHandler creates a new ApacheLoggingHandler that wraps the
// input handler and writes to the given writer
func NewApacheLoggingHandler(handler http.Handler, out io.Writer) http.Handler {
	return &ApacheLoggingHandler{
		handler: handler,
		out:     out,
	}
}

// ServerHTTP tracks http request times and metadata, logging the http request
// after calling the wrapped http.Handler
func (h *ApacheLoggingHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	record := &ApacheLogRecord{
		ResponseWriter: rw,
		ip:             clientIP,
		time:           time.Time{},
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
		status:         http.StatusOK,
		userAgent:      r.UserAgent(),
		elapsedTime:    time.Duration(0),
	}

	startTime := time.Now()
	h.handler.ServeHTTP(record, r)
	finishTime := time.Now()

	record.time = finishTime.UTC()
	record.elapsedTime = finishTime.Sub(startTime)

	record.Log(h.out)
}
