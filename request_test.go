package accessmon

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestRequest_IsHTTP2(t *testing.T) {
	req := &Request{}
	require.False(t, req.IsHTTP2())

	req = &Request{HTTPVersion: "HTTP/1.0"}
	require.False(t, req.IsHTTP2())

	req = &Request{HTTPVersion: "HTTP/1.1"}
	require.False(t, req.IsHTTP2())

	req = &Request{HTTPVersion: ""}
	require.False(t, req.IsHTTP2())

	req = &Request{HTTPVersion: "blah"}
	require.False(t, req.IsHTTP2())

	req = &Request{HTTPVersion: "HTTP/2.0"}
	require.True(t, req.IsHTTP2())
}

func TestRequest_IsIPv6(t *testing.T) {
	req := &Request{}
	require.False(t, req.IsIPv6())

	req = &Request{SourceIP: net.ParseIP("127.0.0.1")}
	require.False(t, req.IsIPv6())

	req = &Request{SourceIP: net.ParseIP("::1")}
	require.True(t, req.IsIPv6())
}

func TestRequest_IsServerError(t *testing.T) {
	req := &Request{}
	require.False(t, req.IsServerError())

	req = &Request{Code: 200}
	require.False(t, req.IsServerError())

	req = &Request{Code: 0}
	require.False(t, req.IsServerError())

	req = &Request{Code: -1}
	require.False(t, req.IsServerError())

	req = &Request{Code: 500}
	require.True(t, req.IsServerError())

	req = &Request{Code: 599}
	require.True(t, req.IsServerError())
}

func TestRequest_String(t *testing.T) {
	req := &Request{}
	require.Equal(t, "<nil> -  [01/Jan/0001:00:00:00 +0000] \"  \" 0 0", req.String())

	line := "127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12"

	parser := &W3CParser{}
	req, err := parser.Parse(line)
	require.NoError(t, err)
	require.Equal(t, line, req.String())
}
