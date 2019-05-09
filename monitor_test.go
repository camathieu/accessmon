package accessmon

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMonitor_AddLineNoAlerter(t *testing.T) {
	mon := NewMonitor(&Config{})
	require.NotNil(t, mon)

	alert, err := mon.AddLine("127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Nil(t, alert)
}

func TestMonitor_AddLine(t *testing.T) {
	mon := NewMonitor(&Config{AlertWindow: time.Second, AlertThreshold: 10})
	require.NotNil(t, mon)
	require.Equal(t, mon.config.AlertWindow, mon.config.StoreWindow)

	alert, err := mon.AddLine("127.0.0.1 - mary [09/May/2018:16:00:42 +0000] \"POST /api/user HTTP/1.0\" 503 12")
	require.NoError(t, err)
	require.Nil(t, alert)

	date, _ := time.Parse(w3cDateLayout, "09/May/2018:16:00:42 +0000")
	require.Equal(t, date, mon.Last())
	require.Len(t, mon.Alerts(), 0)
}
