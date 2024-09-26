package utils

import (
	"time"

	"github.com/hoangtm1601/go-binance-crawler/internal/models"
)

func GetMinute(interval models.CandleInterval) int {
	switch interval {
	case models.OneMin:
		return 1
	case models.FiveMin:
		return 5
	case models.FifteenMin:
		return 15
	case models.ThirtyMin:
		return 30
	case models.SixtyMin:
		return 60
	case models.TwoFortyMin:
		return 240
	case models.SevenTwentyMin:
		return 720
	case models.FourteenFortyMin:
		return 1440
	default:
		return 1
	}
}

func AddInterval(date time.Time, interval models.CandleInterval) time.Time {
	return date.Add(time.Duration(GetMinute(interval)) * time.Minute)
}

func CalcIntervalEnd(date time.Time, interval models.CandleInterval) time.Time {
	return AddInterval(date, interval).Add(-time.Millisecond)
}
