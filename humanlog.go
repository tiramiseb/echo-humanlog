// Package humanlog provides a logger output that makes Echo logs human-readable
//
// Instead of reimplementing log capabilities, this package receives output from
// Labstack's gommon logger and writes them in a human-readable form, thus
// allowing maximum compatibility.
//
// This package provides two things:
//
// * a logreader, where the Echo logger must write
// * a config for the Logger middleware
package humanlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/labstack/echo/middleware"
)

type (
	logEntry struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Prefix  string `json:"prefix"`
		File    string `json:"file"`
		Line    string `json:"line"`
		Message string `json:"message"`
	}

	// HumanReadableLogger is the object where logs must be written: it receives
	// JSON data from Echo logs and writes them as human-readable strings.
	//
	// This logger also supports the Recover middleware
	//
	// The most simple way to use it is the following:
	//
	// 	 e.Logger.SetOutput(humanlog.New(e.Logger.Output()))
	HumanReadableLogger struct {
		output   io.Writer
		toRemove []*regexp.Regexp
	}
)

var (
	// LoggerConfig is a human-readable log format for the Echo Logger middleware
	LoggerConfig = middleware.LoggerConfig{
		Format: "${time_rfc3339} [ HTTP] ${remote_ip} \"${method} ${uri}\" ${status} (in: ${bytes_in}, out: ${bytes_out}, latency: ${latency_human})\n",
	}
)

// New returns an instance of HumanReadableLogger
func New(output io.Writer) *HumanReadableLogger {
	var (
		re *regexp.Regexp
	)
	re, _ = regexp.Compile("\\\\x1b\\[[0-9]*m")
	return &HumanReadableLogger{
		output:   output,
		toRemove: []*regexp.Regexp{re},
	}
}

func (h *HumanReadableLogger) Write(p []byte) (n int, err error) {
	if p[0] == '{' {
		var (
			d   *json.Decoder
			err error
			re  *regexp.Regexp
			l   logEntry
		)
		for _, re = range h.toRemove {
			p = re.ReplaceAll(p, []byte{})
		}
		d = json.NewDecoder(bytes.NewReader(p))
		err = d.Decode(&l)
		if err != nil {
			return fmt.Fprintf(h.output, "ERROR: Could not read log entry as JSON (%s): %s", err.Error(), string(p))
		}
		// It is a JSON object...
		return fmt.Fprintf(h.output, "%s+%s [%5s] %s:%s: %s\n", strings.Split(l.Time, ".")[0], strings.Split(l.Time, "+")[1], l.Level, l.File, l.Line, l.Message)
	}
	// It is a single line
	return h.output.Write(p)
}
