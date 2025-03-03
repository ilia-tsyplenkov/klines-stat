package models

type VBS struct {
	BuyBase   float64 // объём покупок в базовой валюте
	SellBase  float64 // объём продаж в базовой валюте
	BuyQuote  float64 // объём покупок в котируемой валюте
	SellQuote float64 // объём продаж в котируемой валюте
}
