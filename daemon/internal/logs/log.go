package logs

import (
	"fmt"
)

// Message is a log message.
type Message struct {
	Text      string
	IsVerbose bool
}

// Log is a function that logs a message.
type Log func(Message)

// Terse is a log that prints non-verbose messages to STDOUT.
var Terse Log = func(m Message) {
	if !m.IsVerbose {
		fmt.Println(m.Text)
	}
}

// Verbose is a log that prints all messages to STDOUT.
var Verbose Log = func(m Message) {
	fmt.Println(m.Text)
}

// Discard is a log that discards all messages.
var Discard Log

// Write writes a log message.
func (log Log) Write(format string, args ...any) {
	if log != nil {
		log(
			Message{
				Text: fmt.Sprintf(format, args...),
			},
		)
	}
}

// WriteVerbose writes a log message that is only shown if verbose logging is
// enabled.
func (log Log) WriteVerbose(format string, args ...any) {
	if log != nil {
		log(
			Message{
				Text:      fmt.Sprintf(format, args...),
				IsVerbose: true,
			},
		)
	}
}

// WithPrefix returns a log that prefixes all messages with the given string.
func (log Log) WithPrefix(format string, args ...any) Log {
	if log == nil {
		return nil
	}

	prefix := fmt.Sprintf(format, args...)

	return func(m Message) {
		log(
			Message{
				prefix + m.Text,
				m.IsVerbose,
			},
		)
	}
}

// Tee returns a new log that writes to multiple other logs.
func Tee(logs ...Log) Log {
	return func(m Message) {
		for _, log := range logs {
			if log != nil {
				log(m)
			}
		}
	}
}

// Buffer is an in-memory buffer of log messages.
type Buffer []Message

// Log is a Log implementation that appends the message to the buffer.
func (b *Buffer) Log() Log {
	return func(m Message) {
		*b = append(*b, m)
	}
}
