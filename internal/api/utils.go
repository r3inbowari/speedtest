package api

import (
	"strconv"
	"time"
)

// MinString get minimum delay and convert it to string
func MinString(tds []time.Duration) string {
	return strconv.FormatFloat(float64(max(tds).Microseconds())/1000, 'f', 2, 32)
}

func max(tds []time.Duration) time.Duration {
	if len(tds) > 0 {
		m := tds[0]
		for _, v := range tds {
			if v > m {
				m = v
			}
		}
		return m
	}
	return 0
}

func min(tds []time.Duration) time.Duration {
	if len(tds) > 0 {
		m := tds[0]
		for _, v := range tds {
			if v < m {
				m = v
			}
		}
		return m
	}
	return 0
}

// RemoteAverage calculate the average latency using timestamp(len) recorded by the
// remote server. notice that the remote-side-rtt represent prev-round + next-trip
// prev-round + next-trip = timestamp(n+1) - timestamp(n) { 0 ≤ n < len }
func RemoteAverage(remoteTimestamps []string) float64 {
	intRemoteTimestamps := make([]int64, len(remoteTimestamps))
	for i := range remoteTimestamps {
		ist, err := strconv.ParseInt(remoteTimestamps[i], 10, 64)
		if err != nil {
			return 0
		}
		intRemoteTimestamps[i] = ist
	}
	var sum int64 = 0
	for i := 0; i <= len(intRemoteTimestamps)-2; i++ {
		sum += intRemoteTimestamps[i+1] - intRemoteTimestamps[i]
	}
	return float64(sum) / float64(len(intRemoteTimestamps))
}

// LocalAverage directly average by local latency sequence
func LocalAverage(localLatencies []time.Duration) float64 {
	var sum int64 = 0
	for _, latency := range localLatencies {
		sum += int64(latency)
	}
	return float64(sum) / float64(len(localLatencies)) / float64(time.Millisecond)
}

// Jitter obtained by calculating the local delay sequence
func jitter(tds []time.Duration) time.Duration {
	minLatency := min(tds)
	maxLatency := max(tds)
	if maxLatency == 0 || minLatency == 0 {
		return -1
	}
	return maxLatency - minLatency
}

func JitterString(tds []time.Duration) string {
	return strconv.FormatFloat(float64(jitter(tds).Microseconds())/1000, 'f', 2, 32)
}
