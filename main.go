package main

import (
	"flag"
	"io"
	"log"
)

func main() {
	flag.Parse()

	defer closeAllFiles()
	processFlags()

	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading file:", err)
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			log.Fatal("Error writing file:", err)
		}
	}
}
