package steelix_test

import (
	"time"
)

type mockBackoff struct{}

func (m mockBackoff) NextInterval() time.Duration {
	return 0
}
