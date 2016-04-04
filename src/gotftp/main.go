package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pin/tftp"
)

func main() {
	addr := flag.String("s", ":69", "Server address")
	path := flag.String("p", "<path>", "Local file path")
	filename := flag.String("n", "<filename>", "Name of the file on server")
	operation := flag.String("o", "<get|put>", "What to do: download or upload file")
	mode := flag.String("m", "octet", "Transfer mode: 'octet' or 'netascii'")
	flag.Parse()
	if *path == "<path>" {
		fmt.Fprintf(os.Stderr, "missing local path!\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *filename == "<filename>" {
		fmt.Fprintf(os.Stderr, "missing filename!\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *mode != "netascii" && *mode != "octet" {
		fmt.Fprintf(os.Stderr, "invalid mode: %s\n\n", *mode)
		flag.Usage()
		os.Exit(1)
	}
	if *operation == "put" {
		err := send(*addr, *path, *filename, *mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *operation == "get" {
		err := receive(*addr, *path, *filename, *mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "missing or invalid operation!\n\n")
		flag.Usage()
		os.Exit(1)
	}
}

func send(addr string, path string, filename string, mode string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	c, err := tftp.NewClient(addr)
	if err != nil {
		return err
	}
	r, err := c.Send(filename, mode)
	if err != nil {
		return err
	}
	n, err := r.ReadFrom(file)
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

func receive(addr string, path string, filename string, mode string) error {
	c, err := tftp.NewClient(addr)
	if err != nil {
		return err
	}
	w, err := c.Receive(filename, mode)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	n, err := w.WriteTo(file)
	if err != nil {
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}
