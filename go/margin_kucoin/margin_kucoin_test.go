package margin_kucoin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// api.json 需要放到和此文件同一目录
func loadApiKey(ex *MarginKucoin) {
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
	return

	// @ CreateOrder
	order, err := ex.CreateOrder(symbol, "limit", "buy", 0.001 /*amount*/, 0.1 /*price*/, nil)
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
