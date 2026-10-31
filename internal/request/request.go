package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestString, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("unable to read request: %v", err)
		return nil, errors.New("could not parse request")
	}

	parsedRequestLine, err := parseRequestLine(requestString)
	if err != nil {
		log.Printf("could not parse request")
		return nil, errors.New("could not parse request")
	}
	parsedRequest := Request{
		RequestLine: parsedRequestLine,
	}
	return &parsedRequest, nil
}

func parseRequestLine(rs []byte) (RequestLine, error) {
	splitLines := strings.Split(string(rs), "\r\n")
	requestLine := strings.Split(splitLines[0], " ")
	if len(requestLine) != 3 {
		log.Printf("request line invalid, must contain exactly 3 items")
		return RequestLine{}, errors.New("invalid request line")
	}
	if !isOnlyUppercaseAlpha(requestLine[0]) {
		log.Printf("invalid request method, must only have uppercase alphabets")
		return RequestLine{}, errors.New("invalid method")
	}

	requestScheme, httpVer, found := strings.Cut(requestLine[2], "/")
	if !found {
		log.Printf("request line could not be parsed")
		return RequestLine{}, errors.New("request line could not be parsed")
	}

	if httpVer != "1.1" {
		log.Printf("invalid http version, must be HTTP/1.1")
		return RequestLine{}, errors.New("invalid http version")
	}

	if requestScheme != "HTTP" {
		log.Printf("invalid scheme, must be HTTP")
		return RequestLine{}, errors.New("invalid scheme")
	}

	// finally creating the struct to return

	reqLine := RequestLine{
		HttpVersion:   httpVer,
		RequestTarget: requestLine[1],
		Method:        requestLine[0],
	}
	return reqLine, nil
}

// forgive me roadmap, I copied this helpler straight out of google ai overview
func isOnlyUppercaseAlpha(s string) bool {
	// ^[A-Z]+$ matches 1 or more uppercase letters from start to end
	// Change to ^[A-Z]*$ if you want to consider an empty string valid
	matched, _ := regexp.MatchString("^[A-Z]+$", s)
	return matched
}
