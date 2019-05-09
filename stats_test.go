package accessmon

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"strconv"
	"testing"
)

func TestNewStatsCount(t *testing.T) {
	r := &Request{
		SourceIP:    net.ParseIP("127.0.0.1"),
		User:        "user",
		Section:     "/api",
		HTTPVersion: "HTTP/1.0",
		Code:        200,
	}

	requests := make([]*Request, 1000)
	for i := range requests {
		requests[i] = r
	}

	stats := NewStats(requests, 1)
	require.Equal(t, 1000, stats.Count)

	stats = NewStats(requests[:100], 1)
	require.Equal(t, 100, stats.Count)

	stats = NewStats(requests, 1)
	require.Equal(t, 1000, stats.Count)
}

func TestNewStatsPercentages(t *testing.T) {

	requests := make([]*Request, 1000)
	for i := range requests {
		r := &Request{
			SourceIP:    net.ParseIP("127.0.0.1"),
			User:        "user",
			Section:     "/api",
			HTTPVersion: "HTTP/1.0",
			Code:        200,
		}
		requests[i] = r
	}

	stats := NewStats(requests, 1)
	assert.Equal(t, float64(0), stats.Ipv6)
	assert.Equal(t, float64(0), stats.HTTP2)
	assert.Equal(t, float64(0), stats.ServerError)

	ipv6 := net.ParseIP("::1")
	for _, r := range requests[:100] {
		r.SourceIP = ipv6
		r.HTTPVersion = "HTTP/2.0"
		r.Code = 500
	}

	stats = NewStats(requests, 1)
	assert.Equal(t, float64(10), stats.Ipv6)
	assert.Equal(t, float64(10), stats.HTTP2)
	assert.Equal(t, float64(10), stats.ServerError)

}

func TestNewStatsTop(t *testing.T) {

	requests := make([]*Request, 1000)
	for i := range requests {
		r := &Request{
			User:        "user",
			Section:     "/api",
			Code:        200,
			HTTPVersion: "HTTP/1.0",
			SourceIP:    net.ParseIP("127.0.0.1"),
		}
		requests[i] = r
	}

	stats := NewStats(requests, -1)
	assert.Len(t, stats.TopSources, 0)
	assert.Len(t, stats.TopUsers, 0)
	assert.Len(t, stats.TopSection, 0)

	stats = NewStats(requests, 0)
	assert.Len(t, stats.TopSources, 0)
	assert.Len(t, stats.TopUsers, 0)
	assert.Len(t, stats.TopSection, 0)

	stats = NewStats(requests, 1)
	assert.Len(t, stats.TopSources, 1)
	assert.Len(t, stats.TopUsers, 1)
	assert.Len(t, stats.TopSection, 1)

	stats = NewStats(requests, 2)
	assert.Len(t, stats.TopSources, 1)
	assert.Len(t, stats.TopUsers, 1)
	assert.Len(t, stats.TopSection, 1)

	alter := func(i int, from int, to int) {
		index := strconv.Itoa(i)
		ip := net.ParseIP("127.0.0." + index)
		user := "user_" + index
		section := "section_" + index
		for _, r := range requests[from:to] {
			r.SourceIP = ip
			r.User = user
			r.Section = section
		}
	}

	alter(0, 0, 10)     // 10
	alter(1, 10, 100)   // 90
	alter(2, 100, 200)  // 100
	alter(3, 200, 400)  // 200
	alter(4, 400, 1000) // 600

	counts := []int{10, 90, 100, 200, 600}
	check := func(stats *Stats, i int) {
		index := strconv.Itoa(i)
		section := "section_" + index
		user := "user_" + index
		ip := net.ParseIP("127.0.0." + index)

		require.NotNil(t, stats.TopSources[i])
		require.Equal(t, ip.String(), stats.TopSources[i].Key)
		require.Equal(t, counts[i], stats.TopSources[i].Count)

		require.NotNil(t, stats.TopUsers[i])
		require.Equal(t, user, stats.TopUsers[i].Key)
		require.Equal(t, counts[i], stats.TopUsers[i].Count)

		require.NotNil(t, stats.TopSection[i])
		require.Equal(t, section, stats.TopSection[i].Key)
		require.Equal(t, counts[i], stats.TopSection[i].Count)
	}

	stats = NewStats(requests, 5)
	assert.Len(t, stats.TopSources, 5)
	assert.Len(t, stats.TopUsers, 5)
	assert.Len(t, stats.TopSection, 5)

	for i := range counts {
		check(stats, i)
	}
}
