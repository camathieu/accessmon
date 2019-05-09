package accessmon

import (
	"errors"
	"time"
)

// Store accumulates requests
// /!\ NOT THEAD SAFE /!\
type Store struct {
	requests []*Request
}

// AddRequest adds a request to the store
func (s *Store) AddRequest(req *Request) (err error) {
	// We assume and ensure that time is continuously increasing in the log stream
	if len(s.requests) > 0 {
		last := s.requests[len(s.requests)-1]
		if req.Time.Before(last.Time) {
			return errors.New("invalid request in the past")
		}
	}

	s.requests = append(s.requests, req)

	return nil
}

// Since retrieve the requests that are newer than the provided deadline
func (s *Store) Since(deadline time.Time) (requests []*Request) {
	for i := len(s.requests) - 1; i >= 0; i-- {
		if s.requests[i].Time.Before(deadline) || s.requests[i].Time.Equal(deadline) {
			return s.requests[i+1:]
		}
	}

	return s.requests
}

// Clean remove the requests that are older than the provided deadline
func (s *Store) Clean(deadline time.Time) {
	for len(s.requests) > 0 {
		if s.requests[0].Time.Before(deadline) || s.requests[0].Time.Equal(deadline) {
			// NOTE If the type of the slice element is a pointer or a struct with pointer fields
			// which need to be garbage collected there could be a potential memory leak problem
			// where some elements with values are still referenced by the slice a
			// and thus can not be collected. The following line of code fix this problem:
			s.requests[0] = nil
			// The size of the underlying array will not go unbound as we will reuse available capacity
			s.requests = s.requests[1:]
		} else {
			break
		}
	}
}

// Deadline is a convenient helper to generate deadline parameter
func Deadline(now time.Time, window time.Duration) time.Time {
	return now.Add(-window)
}
