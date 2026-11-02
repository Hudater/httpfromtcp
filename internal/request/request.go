package request

import (
	"bytes"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

type StateTracker int

const bufferSize = 8

const (
	Initialized StateTracker = iota
	Done
)

type Request struct {
	RequestLine  RequestLine
	StateTracker StateTracker
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.StateTracker {
	case Initialized:
		n, line, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = line
		r.StateTracker = Done
		return n, nil
	case Done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	parsedRequest := Request{
		StateTracker: Initialized,
	}
	for parsedRequest.StateTracker != Done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		readToIndex += n
		// reached EOF so state is done then write back and reset for next read op
		if err == io.EOF {
			parsedRequest.StateTracker = Done
			break
		}
		if err != nil {
			return nil, err
		}
		consumed, err := parsedRequest.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf[consumed:readToIndex])
		buf = newBuf
		readToIndex -= consumed
	}
	return &parsedRequest, nil
}

func parseRequestLine(rs []byte) (int, RequestLine, error) {
	bytesRead := 0
	byteIndex := bytes.Index(rs, []byte("\r\n"))

	if byteIndex == -1 {
		log.Printf("empty request line")
		return bytesRead, RequestLine{}, nil
	}

	requestLine := strings.Split(string(rs[:byteIndex]), " ")
	if len(requestLine) != 3 {
		log.Printf("request line invalid, must contain exactly 3 items")
		return bytesRead, RequestLine{}, errors.New("invalid request line")
	}
	if !isOnlyUppercaseAlpha(requestLine[0]) {
		log.Printf("invalid request method, must only have uppercase alphabets")
		return bytesRead, RequestLine{}, errors.New("invalid method")
	}

	requestScheme, httpVer, found := strings.Cut(requestLine[2], "/")
	if !found {
		log.Printf("request line could not be parsed")
		return bytesRead, RequestLine{}, errors.New("request line could not be parsed")
	}

	if httpVer != "1.1" {
		log.Printf("invalid http version, must be HTTP/1.1")
		return bytesRead, RequestLine{}, errors.New("invalid http version")
	}

	if requestScheme != "HTTP" {
		log.Printf("invalid scheme, must be HTTP")
		return bytesRead, RequestLine{}, errors.New("invalid scheme")
	}

	// finally creating the struct to return

	reqLine := RequestLine{
		HttpVersion:   httpVer,
		RequestTarget: requestLine[1],
		Method:        requestLine[0],
	}
	bytesRead = byteIndex + 2
	return bytesRead, reqLine, nil
}

// forgive me roadmap, I copied this helpler straight out of google ai overview
func isOnlyUppercaseAlpha(s string) bool {
	// ^[A-Z]+$ matches 1 or more uppercase letters from start to end
	// Change to ^[A-Z]*$ if you want to consider an empty string valid
	matched, _ := regexp.MatchString("^[A-Z]+$", s)
	return matched
}
