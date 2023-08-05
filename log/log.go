package log

import (
	"fmt"
	"strings"
	"time"
)

const (
	// V is Level Zero
	V = 0
	// VV is Level 1
	VV = 1
	// VVV is Level 2
	VVV = 2
)

var (
	// Verbosity is the configured level of the package, any message marked below
	// this value will not be logged.
	Verbosity = 0

	closed   chan struct{}
	messages chan message

	open = false
)

type message struct {
	timestamp time.Time
	level     int
	message   string
	options   []interface{}
}

// Start opens the background logged go routine
func Start() {
	closed = make(chan struct{})
	messages = make(chan message, 20)

	started := make(chan struct{})

	go func() {

		close(started)
		for msg := range messages {
			if Verbosity >= msg.level {
				text := fmt.Sprintf(strings.ToLower(msg.message), msg.options...)
				fmt.Printf("[%v] --- "+text+"\n", msg.timestamp.Format(time.ANSIC))
			}
		}
		close(closed)
	}()
	<-started
	open = true
	Msg(V, "Started log system.")
}

// Close gracefully shutsdown the background logging go routine
func Close() {
	Msg(V, "shutting down log system.")
	open = false
	close(messages)
	<-closed
}

// Msg logs a message
func Msg(level int, text string) {
	Msgf(level, text)
}

// Msgf logs a message with some optional formatting. Uses same syntax as Printf
func Msgf(level int, text string, options ...interface{}) {
	if open {
		messages <- message{
			timestamp: time.Now(),
			level:     level,
			message:   text,
			options:   options,
		}
	}
}
