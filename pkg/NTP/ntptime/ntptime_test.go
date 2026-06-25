package ntptime

import (
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	original := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	secs, frac := ToNTP(original)
	got := FromNTP(secs, frac).UTC()
	if !got.Equal(original) {
		t.Fatalf("round-trip falhou: esperado %v, obtido %v", original, got)
	}
}

func TestEpochOffset(t *testing.T) {
	// 1970-01-01 deve corresponder a 2208988800 segundos NTP (era 1900).
	unixEpoch := time.Unix(0, 0).UTC()
	secs, _ := ToNTP(unixEpoch)
	if secs != 2208988800 {
		t.Fatalf("esperado 2208988800, obtido %d", secs)
	}
}
