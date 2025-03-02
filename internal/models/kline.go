package models

type KLineResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Category string     `json:"category"`
		Symbol   string     `json:"symbol"`
		List     [][]string `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time   int64 `json:"time"`
	IsLast bool  `json:"-"`
	NextTS int64 `json:"-"`
}

// Структура KL:
type Kline struct {
	Pair      string  // название пары в Bitsgap
	TimeFrame string  // период формирования свечи (1m, 15m, 1h, 1d)
	O         float64 // open - цена открытия
	H         float64 // high - максимальная цена
	L         float64 // low - минимальная цена
	C         float64 // close - цена закрытия
	UtcBegin  int64   // время unix начала формирования свечки
	UtcEnd    int64   // время unix окончания формирования свечки
	// VolumeBS  VBS
}
