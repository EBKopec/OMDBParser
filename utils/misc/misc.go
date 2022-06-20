package misc

import (
	"time"
)

var seconds float64

func ParseTimeDuration(t string) (float64, error) {
	timeD, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}
	seconds = timeD.Seconds()
	return seconds, nil
}
