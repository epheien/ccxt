package kucoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// api.json 需要放到和此文件同一目录
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
	ex.Password = data["password"].(string)
}

func TestFetchOrderBook(t *testing.T) {
	ex, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ex.Verbose = true
	loadApiKey(ex)

	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook("BTC/USDT", 5, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("FetchOrderBook:", orderbook)

	// @ FetchBalance
	balance, err := ex.FetchBalance(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("FetchBalance:", ex.Json(balance))

	// @ CreateOrder
	order, err := ex.CreateOrder("ETH/BTC", "limit", "buy", 0.0001, 0.01, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("CreateOrder:", order.Id)

	// @ FetchOrder
	o, err := ex.FetchOrder(order.Id, "ETH/BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("FetchOrder:", ex.Json(o))

	// @ FetchOpenOrders
	openOrders, err := ex.FetchOpenOrders("ETH/BTC", 0, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("FetchOpenOrders:", ex.Json(openOrders))

	// @ CancelOrder
	resp, err := ex.CancelOrder(order.Id, "ETH/BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("CancelOrder:", resp)
}
