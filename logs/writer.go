package logs

import (
	"bytes"
	"strings"
	"sync"
)

// Writer is an adaptor that presents a Logger as an io.WriteCloser.
//
// Each line of text written via Write() is logged as a separate message. Any
// call to write with text that does not end in line separator is buffered until
// a line separator is written or Close() is called. Blank lines are ignored.
//
// Any instance of a LF, CR or CRLF is treated as a line separator. This allows
// usage with Unix-style or Windows-style text output, as well as console output
// that uses CR to overwrite the current line.
type Writer struct {
	// Target is the log that receives the log messages.
	Target Log

	m   sync.Mutex
	buf strings.Builder
}

func (w *Writer) Write(data []byte) (int, error) {
	n := len(data)

	w.m.Lock()
	defer w.m.Unlock()

	for {
		i := bytes.IndexAny(data, "\r\n")

		// No CRs or LFs in the string, write the remaining unterminated line
		// (if any) to the buffer.
		if i == -1 {
			w.buf.Write(data)
			return n, nil
		}

		// Write the data up to the line separator to the buffer (in case it's
		// the tail end of some longer line), then flush the entire line to the
		// logger.
		w.buf.Write(data[:i])
		w.flush()

		sep := data[i]
		data = data[i+1:]

		// If we've found a Windows-style new line (CRLF) we need to consume the
		// extra byte too.
		if sep == '\r' && len(data) > 0 && data[0] == '\n' {
			data = data[1:]
		}
	}
}

// Close closes the writer, producing a log message from any remaining buffered
// text.
func (w *Writer) Close() error {
	w.m.Lock()
	defer w.m.Unlock()

	w.flush()

	return nil
}

// flush writes the contents of w.buf to the log.
//
// An empty buffer is never flushed to the log, which has the effect of removing
// blank lines from the output.
func (w *Writer) flush() {
	if w.buf.Len() > 0 {
		w.Target(
			Message{
				Text: w.buf.String(),
			},
		)
		w.buf.Reset()
	}
}
