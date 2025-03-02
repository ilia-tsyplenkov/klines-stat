package models

// Структура RT:
type RecentTrade struct {
	Tid       string `json:"i"` // id транзакции
	Pair      string `json:"s"` // название валютной пары (как у нас)
	Price     string `json:"p"` // цена транзакции
	Amount    string `json:"v"` // объём транзакции в базовой валюте
	Side      string `json:"S"` // как биржа засчитала эту сделку (как buy или как sell)
	Timestamp int64  `json:"T"` // время UTC UnixNano
}
