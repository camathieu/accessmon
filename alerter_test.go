package accessmon

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var start = time.Date(2019, time.May, 3, 0, 0, 0, 0, time.UTC)

func (a *Alerter) playFixedInterval(now time.Time, interval time.Duration, values []float64) {
	var times []time.Time
	for i := 0; i < len(values); i++ {
		times = append(times, now.Add(time.Duration(i)*interval))
	}
	a.play(now, times, values)
}

func (a *Alerter) play(now time.Time, times []time.Time, values []float64) {
	if len(times) != len(values) {
		panic("times and values length mismatch")
	}
	for i, now := range times {
		a.Check(now, values[i])
	}
}

func TestNewAlerter_CheckNoStart(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0  1  2   3  4   5   6  7   8
	values := []float64{0, 0, 0, 50, 0, 50, 50, 0, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 0)
}

func TestNewAlerter_CheckStart(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0, 1  2   3   4   5   6   7   8
	values := []float64{0, 0, 0, 50, 50, 50, 50, 50, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(5*time.Second), a.alerts[0].Start)
	require.True(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckStart2(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1   2   3   4   5   6   7
	values := []float64{0, 50, 50, 50, 50, 50, 50, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(3*time.Second), a.alerts[0].Start)
	require.True(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckStartNoInit(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1   2   3   4   5   6   7   8
	values := []float64{0, 50, 50, 50, 50, 50, 50, 50, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(3*time.Second), a.alerts[0].Start)
	require.True(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckStart3(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1  2   3   4   5   6   7
	values := []float64{0, 50, 0, 50, 50, 50, 50, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(5*time.Second), a.alerts[0].Start)
	require.True(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckContinue(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1   2   3   4  5   6  7  8   9
	values := []float64{0, 50, 50, 50, 50, 0, 50, 0, 0, 50}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(3*time.Second), a.alerts[0].Start)
	require.True(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckEnd(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1   2   3   4  5  6  7  8  9
	values := []float64{0, 50, 50, 50, 50, 0, 0, 0, 0, 0}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 1)
	require.Equal(t, start.Add(3*time.Second), a.alerts[0].Start)
	require.Equal(t, start.Add(7*time.Second), a.alerts[0].End)
	require.False(t, a.alerts[0].IsOngoing())
}

func TestNewAlerter_CheckNew(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	//                  0   1   2   3   4  5  6  7  8   9  10  11  12 13 14 15 16
	values := []float64{0, 50, 50, 50, 50, 0, 0, 0, 0, 50, 50, 50, 50, 0, 0, 0, 0}
	a.playFixedInterval(start, time.Second, values)

	require.Len(t, a.Alerts(), 2)
	require.Equal(t, start.Add(3*time.Second), a.alerts[0].Start)
	require.Equal(t, start.Add(7*time.Second), a.alerts[0].End)
	require.False(t, a.alerts[0].IsOngoing())
	require.Equal(t, start.Add(11*time.Second), a.alerts[1].Start)
	require.Equal(t, start.Add(15*time.Second), a.alerts[1].End)
	require.False(t, a.alerts[1].IsOngoing())
}

func TestNewAlerter_CheckIrregularWindow(t *testing.T) {
	a := NewAlerter(3*time.Second, float64(10))

	times := []time.Time{start, start.Add(20 * time.Second), start.Add(45 * time.Second), start.Add(46 * time.Second), start.Add(55 * time.Second)}
	values := []float64{0, 50, 0, 0, 50}
	a.play(start, times, values)

	require.Len(t, a.Alerts(), 2)
	require.Equal(t, start.Add(20*time.Second), a.alerts[0].Start)
	require.Equal(t, start.Add(45*time.Second), a.alerts[0].End)
	require.False(t, a.alerts[0].IsOngoing())
	require.Equal(t, start.Add(55*time.Second), a.alerts[1].Start)
	require.True(t, a.alerts[1].IsOngoing())
}
