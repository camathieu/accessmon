package accessmon

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Request represents a HTTP server access log entry
type Request struct {
	SourceIP    net.IP
	User        string
	Time        time.Time
	Method      string
	Path        string
	Section     string
	HTTPVersion string
	Code        int
	Size        int
}

// IsHTTP2 returns if the request is a HTTP/2.0 request
func (req *Request) IsHTTP2() bool {
	return strings.HasPrefix(req.HTTPVersion, "HTTP/2")
}

// IsIPv6 returns if the request is from an IPv6 connection
func (req *Request) IsIPv6() bool {
	if req.SourceIP == nil {
		return false
	}
	return req.SourceIP.To4() == nil
}

// IsServerError returns if the request's response status was a server error
func (req *Request) IsServerError() bool {
	return req.Code >= 500
}

// String representation in W3C Common Log Format
// 127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12
func (req *Request) String() string {
	return fmt.Sprintf(
		"%s - %s [%s] \"%s %s %s\" %d %d",
		req.SourceIP.String(),
		req.User,
		req.Time.Format(w3cDateLayout),
		req.Method,
		req.Path,
		req.HTTPVersion,
		req.Code,
		req.Size)
}
