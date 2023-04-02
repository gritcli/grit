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

// Discard is a log that discards all messages.
func Discard(Message) {}

// Write writes a log message.
func (log Log) Write(format string, args ...any) {
	if log != nil {
		log(Message{fmt.Sprintf(format, args...), false})
	}
}

// WriteVerbose writes a log message that is only shown if verbose logging is
// enabled.
func (log Log) WriteVerbose(format string, args ...any) {
	if log != nil {
		log(Message{fmt.Sprintf(format, args...), true})
	}
}

// WithPrefix returns a new log that prefixes all messages with a string.
func (log Log) WithPrefix(format string, args ...any) Log {
	prefix := fmt.Sprintf(format, args...)
	if log == nil || prefix == "" {
		return log
	}

	return func(m Message) {
		m.Text = prefix + m.Text
		log(m)
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
func (b *Buffer) Log(m Message) {
	*b = append(*b, m)
}
