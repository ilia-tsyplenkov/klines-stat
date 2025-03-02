package models

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
