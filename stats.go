package accessmon

import (
	"sort"
)

// Stats summarize somme statistics about a bunch of Requests
type Stats struct {
	Count       int     // total number of request
	Ipv6        float64 // percentage of IPv6 requests
	HTTP2       float64 // percentage of HTTP2 requests
	ServerError float64 // percentage of Server Error response

	TopUsers   []*CounterValue // top N users
	TopSection []*CounterValue // top N sections
	TopSources []*CounterValue // top N source IPs
}

// NewStats computes statistics about the provided requests
func NewStats(requests []*Request, top int) (s *Stats) {
	s = &Stats{}

	totalReq := 0
	totalIPv6 := 0
	totalHTTP2 := 0
	totalServerError := 0

	for _, req := range requests {
		totalReq++

		if req.IsIPv6() {
			totalIPv6++
		}

		if req.IsHTTP2() {
			totalHTTP2++
		}

		if req.IsServerError() {
			totalServerError++
		}
	}

	s.Count = totalReq
	s.Ipv6 = (float64(totalIPv6) / float64(totalReq)) * 100
	s.HTTP2 = (float64(totalHTTP2) / float64(totalReq)) * 100
	s.ServerError = (float64(totalServerError) / float64(totalReq)) * 100

	if top > 0 {
		sourceCount := newCounter()
		userCount := newCounter()
		sectionCount := newCounter()

		for _, req := range requests {
			sourceCount.incr(req.SourceIP.String())
			userCount.incr(req.User)
			sectionCount.incr(req.Section)
		}

		s.TopSources = sourceCount.top(top)
		s.TopUsers = userCount.top(top)
		s.TopSection = sectionCount.top(top)
	}

	return s
}

// counter holds the key distribution count
type counter struct {
	counts map[string]int
}

func newCounter() (c *counter) {
	c = &counter{}
	c.counts = make(map[string]int)
	return c
}

func (c *counter) incr(key string) {
	c.counts[key]++
}

// top keeps only the n most frequent values
// There is room for optimisation here using a max heap
func (c *counter) top(n int) []*CounterValue {
	var top []*CounterValue
	for key, value := range c.counts {
		top = append(top, &CounterValue{Key: key, Count: value})
	}
	sort.Sort(CounterList(top))
	if len(top) <= n {
		return top
	}
	return top[:n]
}

// CounterValue is a value returned byt the top statistics
type CounterValue struct {
	Key   string
	Count int
}

// CounterList is an interface to sort CounterValues
type CounterList []*CounterValue

func (h CounterList) Len() int           { return len(h) }
func (h CounterList) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h CounterList) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
