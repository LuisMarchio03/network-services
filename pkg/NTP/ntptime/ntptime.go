// Package ntptime converte entre time.Time e timestamps NTP (era 1900).
package ntptime

import "time"

// ntpEpochOffset é a diferença em segundos entre 1900-01-01 e 1970-01-01.
const ntpEpochOffset = 2208988800

// ToNTP converte um time.Time para (segundos, fração) no formato NTP.
func ToNTP(t time.Time) (uint32, uint32) {
	secs := uint64(t.Unix()) + ntpEpochOffset
	frac := uint64(t.Nanosecond()) * (1 << 32) / 1e9
	return uint32(secs), uint32(frac)
}

// FromNTP converte (segundos, fração) NTP de volta para time.Time (UTC).
func FromNTP(secs, frac uint32) time.Time {
	sec := int64(secs) - ntpEpochOffset
	nsec := int64(uint64(frac) * 1e9 / (1 << 32))
	return time.Unix(sec, nsec).UTC()
}
