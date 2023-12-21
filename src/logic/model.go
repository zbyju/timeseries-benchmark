package logic

import "time"

type Snapshot struct {
	StationID          int
	CreatedAt          time.Time
	OutsideTemperature float64
	Voltage            float64
	HeatingTemperature float64
	CoolingTemperature float64
}
