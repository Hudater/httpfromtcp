package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// open file and defer close
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("could not open file with error: %v", err)
	}
	defer file.Close()

	// 8 bytes at a time buffer
	buffer := make([]byte, 8)

	// each line is one unit now
	currLine := ""
	splitBuffer := []string{}

	for {
		// read the open file
		bytesRead, err := file.Read(buffer)

		// End of file means print whatever is there and exit, we done
		if err == io.EOF {
			if len(currLine) > 0 {
				fmt.Printf("read: %s\n", currLine)
			}
			break
		}
		// handle the err, we are not savages
		if err != nil {
			log.Fatalf("Error reading buffer: %v", err)
		}

		// fmt.Printf("OG: '%s' of type: '%T'\n", string(buffer[:bytesRead]), buffer)
		//
		splitBuffer = strings.Split(string(buffer[:bytesRead]), "\n")
		// fmt.Printf("%q\n", splitBuffer)
		for i, part := range splitBuffer {
			if i == 0 {
				currLine += part
			} else {
				fmt.Printf("read: %s\n", currLine)
				currLine = ""
				currLine += part
			}
		}
	}
}
