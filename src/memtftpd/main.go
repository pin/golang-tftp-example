package main

import (
	"fmt"
	"github.com/pin/tftp/v3"
	"io"
	"bytes"
	"os"
	"time"
)

// Map stores "files" in memory.
var m map[string][]byte

// Handler for write (aka PUT) requests.
func writeHandler(filename string, wt io.WriterTo) error {
	_, exists := m[filename]
	if exists {
		return fmt.Errorf("File already exists: %s", filename)
	}
	buffer := &bytes.Buffer{}
	_,err := wt.WriteTo(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't receive %s: %v\n", filename, err)
		return err
	}
	m[filename] = buffer.Bytes()
	return nil
}

// Hander for read (aka GET) requests.
func readHandler(filename string, rt io.ReaderFrom) error {
	b, exists := m[filename]
	if !exists {
		return fmt.Errorf("File not found: %s", filename)
	}
	buffer := bytes.NewBuffer(b)
	_, err := rt.ReadFrom(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't send %s: %v\n", filename, err)
	}
	return nil
}


// Hook for logging on every transfer completion or failure.
type logHook struct {}
func (h *logHook) OnSuccess(stats tftp.TransferStats) {
	fmt.Printf("Transfer of %s to %s complete\n", stats.Filename, stats.RemoteAddr)
}
func (h *logHook) OnFailure(stats tftp.TransferStats, err error) {
	fmt.Printf("Transfer of %s to %s failed: %v\n", stats.Filename, stats.RemoteAddr, err)
}

func main() {
	m = make(map[string][]byte)

	// Start the server.
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetHook(&logHook{})
	go func() {
		err := s.ListenAndServe(":69")
		if err != nil {
			fmt.Fprintf(os.Stdout, "Can't start server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Do some other stuff.
	for {
		_, ok := m["secret-magic-shutdown-file"]
		if ok {
			break
		}
		time.Sleep(5 * time.Minute)
	}

	// Eventually shutdown the server.
	s.Shutdown()
}
