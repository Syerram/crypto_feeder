package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func MockHttpClient(code int, body string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	return server, &http.Client{Transport: transport}
}

type TestCryptoExchange struct{}

func (testExchange TestCryptoExchange) Parse(meta Exchange, row map[string]interface{}) TradeRow {
	return TradeRow{
		Triplet:   fmt.Sprintf("%s-%s", meta.Name, meta.Pair),
		Tid:       strconv.FormatFloat(row["tid"].(float64), 'f', 0, 64),
		Timestamp: strconv.FormatFloat(row["date"].(float64), 'f', 0, 64),
		Amount:    row["amount"].(float64),
		Price:     row["price"].(float64),
	}
}

func TestExchange(t *testing.T) {
	t.Log("Exchange Fetcher")
	{
		server, httpClient := MockHttpClient(200, `[{"tid":1,"date":11111111,"amount":1,"price":1}]`)
		defer server.Close()

		//First lets clear exchanges
		EXCHANGES.ClearExchanges()

		EXCHANGES.RegisterHttpClient(httpClient)
		EXCHANGES.RegisterExchange(Exchange{Name: "test", Pair: "btc_usd", Url: "http://testexchange.com", Interval: 1, Parser: TestCryptoExchange{}})

		var fetchChannel = make(chan []TradeRow)
		var closeChannel = make(chan bool)
		EXCHANGES.Listen(closeChannel, fetchChannel)
		var tradeRows = <-fetchChannel
		var tradeRowsExpected = [1]TradeRow{TradeRow{Triplet: "test-btc_usd",
			Tid:       "1",
			Timestamp: "11111111",
			Amount:    1,
			Price:     1}}

		if !reflect.DeepEqual(tradeRowsExpected[0], tradeRows[0]) {
			t.Fatal("Not the same trade row objects")
		}
	}

	t.Log("Exchange Fetcher handle Panic")
	{
		// send invalid float values
		server, httpClient := MockHttpClient(200, `[{"tid":1,"date":22222222,"amount":"s","price":"s"}]`)
		defer server.Close()

		//First lets clear exchanges
		EXCHANGES.ClearExchanges()

		EXCHANGES.RegisterHttpClient(httpClient)
		EXCHANGES.RegisterExchange(Exchange{Name: "test", Pair: "btc_usd", Url: "http://testexchange.com", Interval: 1, Parser: TestCryptoExchange{}})

		var fetchChannel = make(chan []TradeRow)
		var closeChannel = make(chan bool)
		EXCHANGES.Listen(closeChannel, fetchChannel)
		var tradeRows = <-fetchChannel
		var tradeRowsExpected = []TradeRow{}

		if !reflect.DeepEqual(tradeRowsExpected, tradeRows) {
			t.Fatal("Should be empty")
		}

	}

	t.Log("Handle Exchange Connectivity errors") {

	}

	t.Log("Handle invalid JSON data from Exchange") {

	}

	t.Log("Test Exchange Close Channel") {

	}

}
