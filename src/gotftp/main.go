package main;

import (
	"github.com/pin/tftp"
	"net"
	"fmt"
	"os"
	"log"
	"io"
	"bufio"
	"flag"
)

func main() {
	addrStr := flag.String("s", "localhost:69", "Server address")
	pathStr := flag.String("p", "<path>", "Local file path")
	filenameStr := flag.String("n", "<filename>", "Name of the file on server")
	operation := flag.String("o", "<get|put>", "What to do: download or upload file")
	mode := flag.String("m", "octet", "Transfer mode: 'octet' or 'netascii'")
	flag.Parse()
	if *pathStr == "<path>" {
		fmt.Fprintf(os.Stderr, "missing local path!\n\n");
		flag.Usage()
		os.Exit(1)
	}
	if *filenameStr == "<filename>" {
		fmt.Fprintf(os.Stderr, "missing filename!\n\n");
		flag.Usage()
		os.Exit(1)
	}
	if *mode != "netascii" && *mode != "octet" {
		fmt.Fprintf(os.Stderr, "invalid mode: %s\n\n", *mode);
		flag.Usage()
		os.Exit(1)
	}
	if *operation == "put" {
		putFile(*addrStr, *pathStr, *filenameStr, *mode, *pathStr)
	} else if *operation == "get" {
		getFile(*addrStr, *pathStr, *filenameStr, *mode, *pathStr)
	} else {
		fmt.Fprintf(os.Stderr, "missing or invalid operation!\n\n");
		flag.Usage()
		os.Exit(1)
	}
}

func putFile(addrStr string, pathStr string, filename string, mode string, path string) {
	addr, e := net.ResolveUDPAddr("udp", addrStr)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		return
	}
	file, e := os.Open(pathStr)
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(file)
	log := log.New(os.Stderr, "", log.Ldate | log.Ltime)
	c := tftp.Client{addr, log}
	c.Put(filename, mode, func(writer *io.PipeWriter) {
		n, writeError := r.WriteTo(writer)
		if writeError != nil {
			fmt.Fprintf(os.Stderr, "Can't put %s: %v\n", filename, writeError);
		} else {
			fmt.Fprintf(os.Stderr, "Put %s (%d bytes)\n", filename, n);
		}
		writer.Close()
	})
}

func getFile(addrStr string, pathStr string, filename string, mode string, path string) {
	addr, e := net.ResolveUDPAddr("udp", addrStr)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		return
	}
	file, e := os.Create(pathStr)
	if e != nil {
		panic(e)
	}
	w := bufio.NewWriter(file)
	log := log.New(os.Stderr, "", log.Ldate | log.Ltime)
	c := tftp.Client{addr, log}
	c.Get(filename, mode, func(reader *io.PipeReader) {
		n, readError := w.ReadFrom(reader)
		if readError != nil {
			fmt.Fprintf(os.Stderr, "Can't get %s: %v\n", filename, readError);
		} else {
			fmt.Fprintf(os.Stderr, "Got %s (%d bytes)\n", filename, n);
		}
		w.Flush()
		file.Close()
	})
}
