package models

// Структура RT:
type RecentTrade struct {
	Tid       string // id транзакции
	Pair      string // название валютной пары (как у нас)
	Price     string // цена транзакции
	Amount    string // объём транзакции в базовой валюте
	Side      string // как биржа засчитала эту сделку (как buy или как sell)
	Timestamp int64  // время UTC UnixNano
}
