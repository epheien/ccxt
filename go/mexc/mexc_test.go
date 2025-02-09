package mexc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/epheien/ccxt/go/base"
)

var symbol = "BTC/USDT"
var ex *Mexc
var err error

func setup() {
	var err error
	ex, err = New(nil)
	if err != nil {
		log.Fatal(err)
	}
	ex.Verbose = true
	ex.SetProxy("socks5://127.0.0.1:1080")
	loadApiKey(ex)
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	rc := m.Run()
	teardown()
	os.Exit(rc)
}

func loadApiKey(ex *Mexc) {
	plan, err := ioutil.ReadFile("api.json")
	if err != nil {
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(plan, &data)
	if err != nil {
		return
	}

	ex.ApiKey = data["apiKey"].(string)
	ex.Secret = data["secret"].(string)
	if data["password"] != nil {
		ex.Password = data["password"].(string)
	}
	if data["uid"] != nil {
		ex.Uid = data["uid"].(string)
	}
}

func TestAll(t *testing.T) {
	testFetchMarkets(t)
	//testFetchOrderBook(t)
	//testFetchTrades(t)
	//testFetchTicker(t)
	//testFetchOHLCV(t)
	//testFetchBalance(t)
	//order := testCreateOrder(t); _ = order
	//testFetchOrder(t, "11555864984")
	//testFetchOpenOrders(t)
	//testCancelOrder(t, "11555864984")
}

func testFetchMarkets(t *testing.T) {
	// @ FetchMarkets
	markets, err := ex.FetchMarkets(nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchMarkets:", ex.JsonIndent(markets[1]))
	count := 0
	for _, market := range markets {
		if !market.Active {
			//log.Println(market.Symbol, "is not active!")
			continue
		}
		if market.QuoteId != "USDT" {
			continue
		}
		count += 1
	}
	log.Printf("*/USDT count %d / %d", count, len(markets))
}

func testFetchOrderBook(t *testing.T) {
	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook(symbol, 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrderBook:", symbol, ex.JsonIndent(orderbook))
}

func testFetchTrades(t *testing.T) {
	// @ FetchTrades
	trades, err := ex.FetchTrades(symbol, 0, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	length := len(trades)
	if length >= 3 {
		log.Println("##### FetchTrades:", symbol, ex.JsonIndent(trades[length-4:length-1]))
	} else {
		log.Println("##### FetchTrades:", symbol, ex.JsonIndent(trades))
	}
	if length > 0 {
		length := len(trades)
		log.Println(symbol, "Trade Frequency:", float64(length)*1000/float64(trades[length-1].Timestamp-trades[0].Timestamp))
	}
}

func testFetchTicker(t *testing.T) {
	ticker, err := ex.FetchTicker(symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchTicker:", symbol, ex.JsonIndent(ticker))
}

func testFetchOHLCV(t *testing.T) {
	klines, err := ex.FetchOHLCV(symbol, "1h", 0, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOHLCV:", symbol, ex.JsonIndent(klines))
	log.Println("count:", len(klines))
}

func testFetchBalance(t *testing.T) {
	// @ FetchBalance
	balance, err := ex.FetchBalance(nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchBalance:", ex.JsonIndent(balance))
}

func testCreateOrder(t *testing.T) *base.Order {
	// @ CreateOrder
	order, err := ex.CreateOrder(symbol, "limit", "buy", 0.001 /*amount*/, 10000 /*price*/, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CreateOrder:", symbol, ex.JsonIndent(order))
	return order
}

func testFetchOrder(t *testing.T, orderId string) {
	// @ FetchOrder
	o, err := ex.FetchOrder(orderId, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrder:", ex.JsonIndent(o))
}

func testFetchOpenOrders(t *testing.T) []*base.Order {
	// @ FetchOpenOrders
	openOrders, err := ex.FetchOpenOrders(symbol, 0, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOpenOrders:", ex.JsonIndent(openOrders))
	return openOrders
}

func testCancelOrder(t *testing.T, orderId string) {
	// @ CancelOrder
	resp, err := ex.CancelOrder(orderId, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CancelOrder:", resp)
}
