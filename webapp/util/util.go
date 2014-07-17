package util

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

// for test
var Now = time.Now

func FixedTime(str string, fn func()) {
	backup := Now
	defer func() {
		Now = backup
	}()
	Now = func() time.Time {
		ret, _ := time.Parse("2006-01-02 15:04:05 -0700", str)
		return ret
	}
	fn()
}

func GetTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05 -0700", str)
}

var GenUUID = uuid.New
