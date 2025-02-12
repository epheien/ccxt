package bitmax2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/epheien/ccxt/go/base"
)

var symbol = "BTC/USDT"
var ex *Bitmax2
var err error

func setup() {
	var err error
	ex, err = New(nil)
	if err != nil {
		log.Fatal(err)
	}
	ex.Verbose = true
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

func loadApiKey(ex *Bitmax2) {
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
	//testFetchBalance(t)
	//order := testCreateOrder(t); _ = order
	//testFetchOrder(t, "")
	//testFetchOpenOrders(t)
	//testCancelOrder(t, order)
}

func testFetchOrderBook(t *testing.T) {
	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook(symbol, 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrderBook:", orderbook)
}

func testFetchBalance(t *testing.T) {
	// @ FetchBalance
	balance, err := ex.FetchBalance(nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchBalance:", ex.Json(balance))
}

func testCreateOrder(t *testing.T) *base.Order {
	// @ CreateOrder
	order, err := ex.CreateOrder(symbol, "limit", "buy", 0.001 /*amount*/, 10000 /*price*/, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CreateOrder:", order.Id)
	return order
}

func testFetchOrder(t *testing.T, orderId string) {
	// @ FetchOrder
	o, err := ex.FetchOrder(orderId, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrder:", ex.Json(o))
}

func testFetchOpenOrders(t *testing.T) {
	// @ FetchOpenOrders
	openOrders, err := ex.FetchOpenOrders(symbol, 0, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOpenOrders:", ex.Json(openOrders))
}

func testCancelOrder(t *testing.T, order *base.Order) {
	// @ CancelOrder
	resp, err := ex.CancelOrder(order.Id, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CancelOrder:", resp)
}
