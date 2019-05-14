package use_cases

import (
	"time"
)

type GetTime struct {}

func (g GetTime) GetNowTime() time.Time{
	return time.Now().UTC()
}