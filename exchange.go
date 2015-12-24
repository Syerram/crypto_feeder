package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Parser interface implemented by concreate exchanges
type Parser interface {
	Parse(meta Exchange, row map[string]interface{}) TradeRow
}

// -- Exchange struct ---
type Exchange struct {
	Name     string
	Pair     string
	Url      string
	Parser   Parser
	Interval int
}

// Gets the data from the given URL and parse it using JSON parser
func (exchange *Exchange) Fetch() (data []map[string]interface{}, err error) {
	resp, err := EXCHANGES.HttpClient.Get(exchange.Url)
	if err != nil {
		LOGGER.Error.Println("Error occurred fetching data from exchange [", exchange.Name, "]: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	var resBody, resErr = ioutil.ReadAll(resp.Body)
	if resErr != nil {
		LOGGER.Error.Println("Error occurred reading response body [", exchange.Name, "]: ", resErr)
		return nil, resErr
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		LOGGER.Error.Println("Error occurred Parsing data->json from exchange [", exchange.Name, "]: ", err)
		return nil, err
	}
	return data, nil
}

// Convert raw data into TradeRow struct
func (exchange *Exchange) Parse(rows []map[string]interface{}) (tradeRows []TradeRow, err error) {
	// During parsing, we could encounter panic due to data mismatch or json errors.
	defer func() {
		if r := recover(); r != nil {
			LOGGER.Error.Println("Error occured in parsing on exchange [", exchange.Name, "]: ", r)
			tradeRows = []TradeRow{}
			err = nil // we ignore parsing errors and moving on
		}
	}()
	tradeRows = []TradeRow{}
	for _, row := range rows {
		tradeRow := exchange.Parser.Parse(*exchange, row)
		tradeRows = append(tradeRows, tradeRow)
	}
	return tradeRows, nil
}

// -- End of Exchange --- //

// -- Exchanges struct --- //
type Exchanges struct {
	HttpClient *http.Client
	_Exchanges []Exchange
}

// Register the http client
func (exchanges *Exchanges) RegisterHttpClient(httpClient *http.Client) {
	exchanges.HttpClient = httpClient
}

// Register an exchange
func (exchanges *Exchanges) RegisterExchange(exchange Exchange) {
	exchanges._Exchanges = append(exchanges._Exchanges, exchange)
}

func (exchanges *Exchanges) GetExchanges() []Exchange {
	return exchanges._Exchanges
}

func (exchanges *Exchanges) ClearExchanges() {
	exchanges._Exchanges = []Exchange{}
}

// Listen function launches registered exchanges as go-routines
func (exchanges *Exchanges) Listen(closeChannel chan bool, broadcastChannel chan []TradeRow) {
	for _, _exchange := range exchanges._Exchanges {
		// Each exchange gets its own go-routine
		go func(exchange Exchange) {
			LOGGER.Trace.Println("Listening for exchange: ", exchange.Name)
			// Each exchange pulls data on a specific interval
			var interval = time.Duration(exchange.Interval) * time.Second
			ticker := time.NewTicker(interval)
			// Close out the ticker timer on the exit
			// TIP: Nice way of using defer functions to close out any resources before go-routine exits
			//		We could have done in closechannel select as well but this seems like more natural
			defer func() {
				ticker.Stop()
			}()
			// Infinite loop until instructed to close
			// Fetch on each ticker and broadcast the data on the given channel
			for {
				select {
				case <-ticker.C:
					if broadcastChannel != nil {
						var data, err = exchange.Fetch()
						if err != nil {
							continue
						}
						var processedData, _ = exchange.Parse(data)
						broadcastChannel <- processedData
					}
				case <-closeChannel:
					broadcastChannel = nil
					return
				}
			}
		}(_exchange)
	}
}

var EXCHANGES = Exchanges{}

// -- End of Exchanges --- //

// -- Traderow Struct to hold Trade data --//
type TradeRow struct {
	Triplet   string
	Tid       string
	Timestamp string
	Amount    float64
	Price     float64
}

// Convert to JSON string
func (tradeRow *TradeRow) ToJson() []byte {
	var tradeRowAsJson, err = json.Marshal(tradeRow)
	if err != nil {
		//TODO: raise error
	}
	return tradeRowAsJson
}

//-- End of Traderow ---//

// Setup HTTP Client to accept all HTTPS domains
func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	EXCHANGES.RegisterHttpClient(&http.Client{Transport: tr})
}
