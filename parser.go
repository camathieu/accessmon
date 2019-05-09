package accessmon

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Parser is an interface that parses a line of log to a Request object
type Parser interface {
	Parse(line string) (*Request, error)
}

// W3CParser is a Parser implementation for W3C Common Logfile Format
// see : https://www.w3.org/Daemon/User/Config/Logging.html
// 127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12
type W3CParser struct{}

const w3cDateLayout = "02/Jan/2006:15:04:05 -0700"

// Consume the line until the separator is found and returns
// - the text until the separator
// - the rest of the line after the separator
func parseNext(line string, sep string) (newline string, part string, err error) {
	if len(line) < len(sep)+1 {
		return "", "", newParsingError("invalid line too short")
	}
	i := strings.Index(line, sep)
	if i <= 0 {
		return "", "", newParsingError("invalid line separator \"" + sep + "\" not found")
	}

	part, newline = line[:i], line[i+len(sep):]

	return newline, part, nil
}

// Parse a W3C Common Log Format line into a Request object
func (parser *W3CParser) Parse(line string) (req *Request, err error) {

	// Trim new lines

	line = strings.TrimRight(line, "\r\n")

	// Split

	line, ipStr, err := parseNext(line, " - ")
	if err != nil {
		return nil, err
	}

	line, username, err := parseNext(line, " [")
	if err != nil {
		return nil, err
	}

	line, dateStr, err := parseNext(line, "] \"")
	if err != nil {
		return nil, err
	}

	line, method, err := parseNext(line, " ")
	if err != nil {
		return nil, err
	}

	line, path, err := parseNext(line, " ")
	if err != nil {
		return nil, err
	}

	line, httpVersion, err := parseNext(line, "\" ")
	if err != nil {
		return nil, err
	}

	line, codeStr, err := parseNext(line, " ")
	if err != nil {
		return nil, err
	}

	sizeStr := line

	// Parse

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, newParsingError("invalid ip")
	}

	date, err := time.Parse(w3cDateLayout, dateStr)
	if err != nil {
		return nil, newParsingError("invalid date")
	}

	section := parseSection(path)

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return nil, newParsingError("invalid code")
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, newParsingError("invalid size")
	}

	// Build

	req = &Request{
		SourceIP:    ip,
		User:        username,
		Time:        date,
		Method:      method,
		Path:        path,
		Section:     section,
		HTTPVersion: httpVersion,
		Code:        code,
		Size:        size,
	}

	return req, nil
}

func parseSection(path string) (section string) {
	path = strings.TrimPrefix(path, "/")
	i := strings.Index(path, "/")
	if i <= 0 {
		return "/" + path
	}
	return "/" + path[:i]
}

func newParsingError(detail interface{}) error {
	return fmt.Errorf("parsing Error : %s", detail)
}
