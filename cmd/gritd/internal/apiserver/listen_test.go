package apiserver_test

import (
	"net"
	"os"
	"path/filepath"
	"time"

	. "github.com/gritcli/grit/cmd/gritd/internal/apiserver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Listen()", func() {
	It("returns a listener for the given socket", func() {
		dir, err := os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
		defer os.RemoveAll(dir)

		socket := filepath.Join(dir, "test.sock")
		l, err := Listen(socket)
		Expect(err).ShouldNot(HaveOccurred())
		defer l.Close()

		go func() {
			defer GinkgoRecover()
			conn, err := net.Dial("unix", socket)
			Expect(err).ShouldNot(HaveOccurred())
			conn.Close()
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(50 * time.Millisecond)
			l.Close() // prevent hanging indefinitely
		}()

		conn, err := l.Accept()
		Expect(err).ShouldNot(HaveOccurred())
		conn.Close()
	})

	It("deletes the socket file if it already exists", func() {
		fp, err := os.CreateTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())
		defer os.Remove(fp.Name())

		socket := fp.Name()
		fp.Close()

		l, err := Listen(socket)
		Expect(err).ShouldNot(HaveOccurred())
		defer l.Close()

		go func() {
			defer GinkgoRecover()
			conn, err := net.Dial("unix", socket)
			Expect(err).ShouldNot(HaveOccurred())
			conn.Close()
		}()

		go func() {
			defer GinkgoRecover()
			time.Sleep(50 * time.Millisecond)
			l.Close() // prevent hanging indefinitely
		}()

		conn, err := l.Accept()
		Expect(err).ShouldNot(HaveOccurred())
		conn.Close()
	})
})
