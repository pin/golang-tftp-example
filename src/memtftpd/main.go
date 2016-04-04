package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pin/tftp"
)

func main() {
	addr := flag.String("l", ":69", "Address to listen")
	flag.Parse()
	b := &backend{}
	b.m = make(map[string][]byte)
	s := tftp.NewServer(b.handleRead, b.handleWrite)
	err := s.ListenAndServe(*addr)
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}

type backend struct {
	m  map[string][]byte
	mu sync.Mutex
}

func (b *backend) handleWrite(filename string, wt io.WriterTo) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.m[filename]
	if ok {
		fmt.Fprintf(os.Stderr, "file %s already exists\n", filename)
		return fmt.Errorf("file already exists")
	}
	buf := &bytes.Buffer{}
	n, err := wt.WriteTo(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't receive %s: %v\n", filename, err)
		return err
	}
	b.m[filename] = buf.Bytes()
	fmt.Fprintf(os.Stderr, "received %s (%d bytes)\n", filename, n)
	return nil
}

func (b *backend) handleRead(filename string, rf io.ReaderFrom) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bs, ok := b.m[filename]
	if !ok {
		fmt.Fprintf(os.Stderr, "file %s not found\n", filename)
		return fmt.Errorf("file not found")
	}
	n, err := rf.ReadFrom(bytes.NewBuffer(bs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't send %s: %v\n", filename, err)
		return err
	}
	fmt.Fprintf(os.Stderr, "sent %s (%d bytes)\n", filename, n)
	return nil
}
