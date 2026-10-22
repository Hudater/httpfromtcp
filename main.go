package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// open file, close is done by `getLinesChannel()`
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("could not open file with error: %v", err)
	}

	lineChannel := getLinesChannel(file)
	// channels don't return index
	for currLine := range lineChannel {
		fmt.Printf("read: %s\n", currLine)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChan := make(chan string)

	go func() {
		// 8 bytes at a time buffer
		buffer := make([]byte, 8)

		// TODO: read more about this, it runs on LIFO priniple
		defer f.Close()
		defer close(lineChan)

		currLine := ""
		splitBuffer := []string{}

		for {
			// read the open file
			bytesRead, err := f.Read(buffer)

			// End of file means send to channel whatever is there
			if err == io.EOF {
				if len(currLine) > 0 {
					lineChan <- currLine
				}
				break
			}
			// handle the err, we are not savages
			if err != nil {
				log.Fatalf("Error reading buffer: %v", err)
			}

			splitBuffer = strings.Split(string(buffer[:bytesRead]), "\n")
			for i, part := range splitBuffer {
				if i == 0 {
					currLine += part
				} else {
					lineChan <- currLine
					currLine = ""
					currLine += part
				}
			}
		}
	}()
	return lineChan
}
