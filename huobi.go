package main

import (
	"fmt"
	"strconv"
)

type Huobi struct{}

//Huobi
func (huobi Huobi) Parse(meta Exchange, row map[string]interface{}) TradeRow {
	return TradeRow{
		Triplet:   fmt.Sprintf("%s-%s", meta.Name, meta.Pair),
		Tid:       strconv.FormatFloat(row["tid"].(float64), 'f', 0, 64),
		Timestamp: strconv.FormatFloat(row["date"].(float64), 'f', 0, 64),
		Amount:    row["amount"].(float64),
		Price:     row["price"].(float64),
	}
}

func init() {
	EXCHANGES.RegisterExchange(Exchange{Name: "huobi", Pair: "btc_usd", Url: "http://api.huobi.com/usdmarket/trades_usd_btc_json.js", Interval: 60, Parser: Huobi{}})
}
