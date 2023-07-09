package kucoin_hf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/georgexdz/ccxt/go/base"
)

var symbol = "BTC/USDT"
var ex *Kucoin
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

func loadApiKey(ex *Kucoin) {
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
	testFetchOrderBook(t)
	//testFetchTicker(t)
	//testFetchOHLCV(t)
	testFetchBalance(t)
	//order := testCreateOrder(t); _ = order
	//testFetchOrder(t, "63a42a5f77056300018c3e48")
	//testFetchOpenOrders(t)
	//testCancelOrder(t, "63a42a5f77056300018c3e48")

	// 余额划转
	//result := ex.ApiFunc("privatePostV2AccountsInnerTransfer",
	//	map[string]interface{}{
	//		"clientOid": ex.Uuid(),
	//		"currency": "BTC",
	//		"from": "trade",
	//		"to": "trade_hf",
	//		"amount": "0.01",
	//	},
	//	nil, nil,
	//)
	//log.Println(result)
}

func testFetchOrderBook(t *testing.T) {
	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook(symbol, 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrderBook:", symbol, ex.JsonIndent(orderbook))
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
	balance.Info = nil
	log.Println("##### FetchBalance:", ex.JsonIndent(balance))
}

func testCreateOrder(t *testing.T) *base.Order {
	// @ CreateOrder
	order, err := ex.CreateOrder(symbol, "limit", "buy", 0.001 /*amount*/, 10000 /*price*/, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CreateOrder:", symbol, order.Id)
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
	log.Println("##### FetchOpenOrders:", ex.Json(openOrders))
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
