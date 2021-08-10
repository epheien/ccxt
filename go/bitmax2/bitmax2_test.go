package bitmax2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

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

func TestFetchOrderBook(t *testing.T) {
	symbol := "BTC/USDT"
	ex, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ex.Verbose = true
	loadApiKey(ex)

	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook(symbol, 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrderBook:", orderbook)

	// @ FetchBalance
	balance, err := ex.FetchBalance(nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchBalance:", ex.Json(balance))

	// @ CreateOrder
	order, err := ex.CreateOrder(symbol, "limit", "buy", 0.001 /*amount*/, 10000 /*price*/, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CreateOrder:", order.Id)

	// @ FetchOrder
	o, err := ex.FetchOrder(order.Id, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrder:", ex.Json(o))

	// @ FetchOpenOrders
	openOrders, err := ex.FetchOpenOrders(symbol, 0, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOpenOrders:", ex.Json(openOrders))

	// @ CancelOrder
	resp, err := ex.CancelOrder(order.Id, symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CancelOrder:", resp)
}
