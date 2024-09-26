package models

import (
	"gorm.io/gorm"
)

type CandleInterval string

const (
	OneMin           CandleInterval = "1min"
	FiveMin          CandleInterval = "5min"
	FifteenMin       CandleInterval = "15min"
	ThirtyMin        CandleInterval = "30min"
	SixtyMin         CandleInterval = "60min"
	TwoFortyMin      CandleInterval = "240min"
	SevenTwentyMin   CandleInterval = "720min"
	FourteenFortyMin CandleInterval = "1440min"
)

type Candle struct {
	gorm.Model
	Symbol   string         // idx_symbol_interval_start
	Interval CandleInterval // idx_symbol_interval_start
	Start    int64          // idx_symbol_interval_start
	End      int64
	LastEnd  int64
	Op       float64 // Changed from string to float64
	Hi       float64 // Changed from string to float64
	Lo       float64 // Changed from string to float64
	Cl       float64 // Changed from string to float64
	Bv       float64 // Changed from string to float64
	Qv       float64 // Changed from string to float64
	Tbv      float64 // Changed from string to float64
	Tqv      float64 // Changed from string to float64
	Cnt      int64
}
