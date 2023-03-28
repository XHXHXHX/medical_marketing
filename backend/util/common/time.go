package common

import "time"

func PTime(t time.Time) *time.Time {
	return &t
}

func PNow() *time.Time {
	return PTime(time.Now())
}

func PTimeUnix(s int64) *time.Time {
	if s == 0 {
		return nil
	}
	return PTime(time.Unix(s, 0))
}

func TimeToUnix(t *time.Time) int64 {
	if t == nil {
		return 0
	}

	return t.Unix()
}