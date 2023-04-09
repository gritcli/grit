package logs_test

import (
	. "github.com/gritcli/grit/daemon/internal/logs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Log", func() {
	Describe("func WithPrefix()", func() {
		It("adds the prefix to the message text", func() {
			var buffer Buffer

			buffer.
				Log().
				WithPrefix("outer: ").
				WithPrefix("inner: ").
				Write("message")

			Expect(buffer).To(ConsistOf(
				Message{
					Text: "outer: inner: message",
				},
			))
		})
	})
})
