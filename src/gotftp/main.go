package main

import (
	"flag"
	"fmt"
	"github.com/pin/tftp/v3"
	"os"
)

func putFile(addr string, localPath string, remoteFilename string, mode string) error {
	c, err := tftp.NewClient(addr)
	if err != nil {
		return fmt.Errorf("connecting %s: %v", addr, err)
	}
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("can't open %s: %v", localPath, err)
	}
	rf, err := c.Send(remoteFilename, mode)
	if err != nil {
		return fmt.Errorf("starting transfer: %v", err)
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("transferring %s: %v", remoteFilename, err)
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

func getFile(addr string, localPath string, localFilename string, mode string) error {
	c, err := tftp.NewClient(addr)
	if err != nil {
		return fmt.Errorf("connecting %s: %v", addr, err)
	}
	wt, err := c.Receive(localFilename, mode)
	if err != nil {
		return fmt.Errorf("requesting transfer: %v", err)
	}
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("creating %s: %v", localPath, err)
	}
	// Optionally obtain transfer size before actual data.
	if n, ok := wt.(tftp.IncomingTransfer).Size(); ok {
		fmt.Printf("Transfer size: %d\n", n)
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Errorf("receiving %s: %v", localFilename, err)
	}
	fmt.Printf("%d bytes received.\n", n)
	return nil
}

func main() {
	addr := flag.String("a", "localhost:69", "Server address")
	localPath := flag.String("l", "", "Local file path")
	remoteFilename := flag.String("r", "", "Remote filename")
	putOperation := flag.Bool("p", false, "Upload (aka PUT) transfer")
	getOperation := flag.Bool("g", false, "Download (aka GET) transfer")
	mode := flag.String("m", "octet", "Mode of transfer: 'octet' or 'netascii'")
	flag.Parse()

	if *localPath == "" {
		fmt.Fprintf(os.Stderr, "Error: Local file path is missing.\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *remoteFilename == "" {
		fmt.Fprintf(os.Stderr, "Error: Remote filename is missing.\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *mode != "netascii" && *mode != "octet" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s\n\n", *mode)
		flag.Usage()
		os.Exit(1)
	}
	if *getOperation && *putOperation {
		fmt.Fprintf(os.Stderr, "Error: Upload and download at once is confusing. Choose one.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *putOperation {
		err := putFile(*addr, *localPath, *remoteFilename, *mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Uploading: %v\n", err)
			os.Exit(3)
		}
	} else if *getOperation {
		err := getFile(*addr, *localPath, *remoteFilename, *mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Downloading: %v\n", err)
			os.Exit(2)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: Choose upload or download.\n\n")
		flag.Usage()
		os.Exit(1)
	}
}
