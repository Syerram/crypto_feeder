package main

import (
	"fmt"
	"strconv"
)

type Bitfinex struct{}

//Bitfinex
func (bitfinex Bitfinex) Parse(meta Exchange, row map[string]interface{}) TradeRow {
	amt, err := strconv.ParseFloat(row["amount"].(string), 64)
	//TODO: What should we do here? Raise error or just ignore
	if err != nil {
		amt = 0.00
	}
	price, err := strconv.ParseFloat(row["price"].(string), 64)
	if err != nil {
		price = 0.00
	}
	return TradeRow{
		Triplet:   fmt.Sprintf("%s-%s", meta.Name, meta.Pair),
		Tid:       strconv.FormatFloat(row["tid"].(float64), 'f', 0, 64),
		Timestamp: strconv.FormatFloat(row["timestamp"].(float64), 'f', 0, 64),
		Amount:    amt,
		Price:     price,
	}
}

func init() {
	EXCHANGES.RegisterExchange(Exchange{Name: "bitfinex", Pair: "btc_usd", Url: "https://api.bitfinex.com/v1/trades/btcusd", Interval: 60, Parser: Bitfinex{}})
}
