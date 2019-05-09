package accessmon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStore_AddRequest(t *testing.T) {
	p := &Store{}

	now := time.Now()
	err := p.AddRequest(&Request{Time: now})
	require.NoError(t, err)

	err = p.AddRequest(&Request{Time: now})
	require.NoError(t, err)

	err = p.AddRequest(&Request{Time: now.Add(1 * time.Second)})
	require.NoError(t, err)

	err = p.AddRequest(&Request{Time: now.Add(-1 * time.Second)})
	require.Error(t, err)

	require.Len(t, p.requests, 3)
}

func TestStore_GetRequests(t *testing.T) {
	var err error

	p := &Store{}
	now := time.Date(2019, time.May, 3, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 60; i++ {
		err = p.AddRequest(&Request{Time: now})
		require.NoError(t, err)
		err = p.AddRequest(&Request{Time: now})
		require.NoError(t, err)
		now = now.Add(time.Second)
	}

	now = now.Add(-1 * time.Second)

	reqs := p.Since(Deadline(now, -1*time.Second))
	require.Len(t, reqs, 0)

	reqs = p.Since(Deadline(now, 0*time.Second))
	require.Len(t, reqs, 0)

	reqs = p.Since(Deadline(now, 24*time.Hour))
	require.Len(t, reqs, 120)

	reqs = p.Since(Deadline(now, 1*time.Second))
	require.Len(t, reqs, 2)

	reqs = p.Since(Deadline(now, 10*time.Second))
	require.Len(t, reqs, 20)
	require.Equal(t, reqs[len(reqs)-1].Time, reqs[len(reqs)-2].Time)
	require.Equal(t, reqs[len(reqs)-1].Time, now)

	reqs = p.Since(Deadline(now, 60*time.Second))
	require.Len(t, reqs, 120)
}

func TestStore_Clean(t *testing.T) {
	var err error

	p := &Store{}
	now := time.Date(2019, time.May, 3, 0, 0, 0, 0, time.UTC)

	p.Clean(Deadline(now, time.Minute))

	for i := 0; i < 60; i++ {
		err = p.AddRequest(&Request{Time: now})
		require.NoError(t, err)
		err = p.AddRequest(&Request{Time: now})
		require.NoError(t, err)
		now = now.Add(time.Second)
	}

	now = now.Add(-1 * time.Second)

	require.Equal(t, 120, len(p.Since(Deadline(now, time.Minute))))

	p.Clean(Deadline(now, 10*time.Second))

	require.Equal(t, 20, len(p.Since(Deadline(now, time.Minute))))
}
