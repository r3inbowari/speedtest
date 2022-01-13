package test

import (
	"fmt"
	"speedtest/internal/api"
	"strconv"
	"testing"
)

func BenchmarkStrAppend0(b *testing.B) {
	delta := api.InitialDelta
	for i := 0; i < b.N; i++ {
		_ = append(api.MessagePing, []byte(strconv.FormatInt(delta, 10))...)
		delta++
	}
}

func BenchmarkStrAppend1(b *testing.B) {
	delta := api.InitialDelta
	const P = "PING %d"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf(P, delta)
		delta++
	}
}

func BenchmarkStrAppend2(b *testing.B) {
	delta := api.InitialDelta
	const P = "PING "
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf(P + strconv.FormatInt(delta, 10))
		delta++
	}
}

func BenchmarkStrAppend3(b *testing.B) {
	delta := api.InitialDelta
	for i := 0; i < b.N; i++ {
		_ = strconv.AppendInt(api.MessagePing, delta, 10)
		delta++
	}
}
