package binance

import (
	"encoding/json"
	"fmt"
	"github.com/georgexdz/ccxt/go/generated/binance"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// api.json 需要放到和此文件同一目录
func loadApiKey(ex *binance.Binance) {
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
}

func TestFetchOrderBook(t *testing.T) {
	ex, err := binance.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ex.Verbose = true
	loadApiKey(ex)
	fmt.Println(ex.ApiKey)

	// @ FetchOrderBook
	orderbook, err := ex.FetchOrderBook("BTC/USDT", 5, nil)
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
	order, err := ex.CreateOrder("BTC/USDT", "limit", "buy", 0.002, 8000., nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### CreateOrder:", order.Id)

	// @ FetchOrder
	o, err := ex.FetchOrder(order.Id, "BTC/USDT", nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOrder:", ex.Json(o))

	// @ FetchOpenOrders
	openOrders, err := ex.FetchOpenOrders("BTC/USDT", 0, 1000, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("##### FetchOpenOrders:", ex.Json(openOrders))

	// @ CancelOrder
	for _, order := range openOrders {
		resp, err := ex.CancelOrder(order.Id, "BTC/USDT", nil)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("##### CancelOrder:", resp)
	}

}