package dto

type CandleDto struct {
	Symbol   string  `json:"symbol"`
	Interval string  `json:"interval"`
	Start    int64   `json:"start"`
	End      int64   `json:"end"`
	Op       float64 `json:"op"`
	Hi       float64 `json:"hi"`
	Lo       float64 `json:"lo"`
	Cl       float64 `json:"cl"`
	Bv       float64 `json:"bv"`
	Qv       float64 `json:"qv"`
	Cnt      int64   `json:"cnt"`
	Tbv      float64 `json:"tbv"`
	Tqv      float64 `json:"tqv"`
	LastEnd  int64   `json:"last_end"`
}
