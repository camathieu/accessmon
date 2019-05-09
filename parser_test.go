package accessmon

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// 127.0.0.1 - mary [09/May/2018:16:00:42 +0000] "POST /api/user HTTP/1.0" 503 12

func TestParseNext(t *testing.T) {
	_, _, err := parseNext("", " - ")
	require.Error(t, err)

	_, _, err = parseNext("this line is not valid", " - ")
	require.Error(t, err)

	line, part, err := parseNext("part - line", " - ")
	require.NoError(t, err)
	require.Equal(t, "part", part)
	require.Equal(t, "line", line)

	line, part, err = parseNext("世界 - 你好好好", " - ")
	require.NoError(t, err)
	require.Equal(t, "世界", part)
	require.Equal(t, "你好好好", line)

	line, part, err = parseNext("::1 - 世界", " - ")
	require.NoError(t, err)
	require.Equal(t, "::1", part)
	require.Equal(t, "世界", line)
}

func TestW3CParser_Parse(t *testing.T) {
	parser := &W3CParser{}

	_, err := parser.Parse("")
	require.Errorf(t, err, "missing error")

	req, err := parser.Parse("127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", req.SourceIP.String())
	require.Equal(t, "mary", req.User)
	require.Equal(t, "09/May/2018:16:00:42 +0000", req.Time.Format(w3cDateLayout))
	require.Equal(t, "POST", req.Method)
	require.Equal(t, "/api/user", req.Path)
	require.Equal(t, "/api", req.Section)
	require.Equal(t, "HTTP/1.0", req.HTTPVersion)
	require.Equal(t, 503, req.Code)
	require.Equal(t, 12, req.Size)

	req, err = parser.Parse("127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST / HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Equal(t, "/", req.Section)
}

func TestW3CParser_ParseIPv6(t *testing.T) {
	parser := &W3CParser{}

	req, err := parser.Parse("::1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Equal(t, "::1", req.SourceIP.String())
}

func TestW3CParser_ParseUTF8(t *testing.T) {
	parser := &W3CParser{}

	req, err := parser.Parse("127.0.0.1 - 世界 [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Equal(t, "世界", req.User)
}

func TestW3CParser_ParseError(t *testing.T) {
	parser := &W3CParser{}

	inputs := []string{
		"invalid - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12",
		"127.0.0.1 - mary [invalid] \"POST /api/user HTTP/1.0\" 503 12",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" invalid 12",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 invalid",
		"127.0.0.1",
		"127.0.0.1 - mary",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000]",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\"",
		"127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503",
	}

	for _, input := range inputs {
		_, err := parser.Parse(input)
		require.Error(t, err, input)
	}
}

func TestParseSection(t *testing.T) {
	section := parseSection("")
	require.Equal(t, "/", section)

	section = parseSection("/")
	require.Equal(t, "/", section)

	section = parseSection("/api")
	require.Equal(t, "/api", section)

	section = parseSection("/api/path")
	require.Equal(t, "/api", section)

	section = parseSection("/世界/你好好好")
	require.Equal(t, "/世界", section)
}
