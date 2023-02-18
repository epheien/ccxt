package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	urllib "net/url"
	"strings"
)

type Gateio struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Gateio, err error) {
	ex = new(Gateio)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Gateio) Describe() []byte {
	return []byte(`
{
    "id": "gateio",
    "name": "Gate.io",
    "countries": [
        "CN"
    ],
    "version": "4",
    "rateLimit": 1000,
    "pro": true,
    "has": {
        "CORS": false,
        "createMarketOrder": false,
        "fetchCurrencies": true,
        "fetchTickers": true,
        "withdraw": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchTransactions": true,
        "createDepositAddress": true,
        "fetchDepositAddress": true,
        "fetchClosedOrders": false,
        "fetchOHLCV": true,
        "fetchOpenOrders": true,
        "fetchOrderTrades": true,
        "fetchOrders": true,
        "fetchOrder": true,
        "fetchMyTrades": true
    },
    // 值只支持字符串形式
    "timeframes": {
    //     "1m": 60,
    //     "5m": 300,
    //     "10m": 600,
    //     "15m": 900,
    //     "30m": 1800,
    //     "1h": 3600,
    //     "2h": 7200,
    //     "4h": 14400,
    //     "6h": 21600,
    //     "12h": 43200,
    //     "1d": 86400,
    //     "1w": 604800
    },
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/31784029-0313c702-b509-11e7-9ccc-bc0da6a0e435.jpg",
        "api": {
            "public": "https://api.gateio.ws/api/v4",
            "private": "https://api.gateio.ws/api/v4"
        },
        "www": "https://gate.io/",
        "doc": "https://www.gate.io/docs/apiv4/zh_CN/index.html"
    },
    "api": {
        "public": {
            "get": [
                "spot/order_book"
            ]
        },
        "private": {
            "get": [
                "spot/accounts",
                "spot/orders",
                "spot/orders/{order_id}",
                "margin/accounts"
            ],
            "post": [
                "spot/orders",
                "wallet/transfers",
                "wallet/sub_account_transfers"
            ],
            "delete": [
                "spot/orders/{order_id}"
            ]
        }
    },
    "fees": {
        "trading": {
            "tierBased": true,
            "percentage": true,
            "maker": 0.002,
            "taker": 0.002
        }
    },
    "exceptions": {
        "exact": {
            "BALANCE_NOT_ENOUGH": "InsufficientFunds",
            "MARGIN_BALANCE_NOT_ENOUGH": "InsufficientFunds",
            "FUTURES_BALANCE_NOT_ENOUGH": "InsufficientFunds",
            "ORDER_NOT_FOUND": "OrderNotFound"
        }
    },
    "options": {
        "fetchTradesMethod": "public_get_tradehistory_id",
        "limits": {
            "cost": {
                "min": {
                    "BTC": 0.0001,
                    "ETH": 0.001,
                    "USDT": 1
                }
            }
        },
        "account": "spot"
    },
}
`)
}

func (self *Gateio) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	return &Market{
		Id:     li[0] + "_" + li[1], // 需要大写
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *Gateio) LoadMarkets() map[string]*Market {
	return nil
}

func (self *Gateio) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	marketId := self.MarketId(symbol)
	request := map[string]interface{}{
		"currency_pair": marketId,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFunc("publicGetSpotOrderBook", self.Extend(request, params), nil, nil)
	timestamp := self.SafeInteger(response, "update")
	orderbook := self.ParseOrderBook(response, timestamp, "bids", "asks", 0, 1)
	return orderbook, nil
}

func (self *Gateio) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	response := self.ApiFuncReturnList("privateGetSpotAccounts", params, nil, nil)
	result := map[string]interface{}{
		"info": response,
	}
	for _, one := range response {
		account := self.Account()
		free := self.SafeFloat(one, "available")
		used := self.SafeFloat(one, "locked")
		cc := self.SafeString(one, "currency")
		account["free"] = free
		account["used"] = used
		account["total"] = free + used
		result[cc] = account
	}
	return self.ParseBalance(result), nil
}

func (self *Gateio) genSign(method, url, query, body string) map[string]interface{} {
	timestamp := self.Milliseconds() / 1000
	m := sha512.New()
	if body != "" {
		m.Write([]byte(body))
	}
	hashedPayload := hex.EncodeToString(m.Sum(nil))
	s := fmt.Sprintf("%s\n%s\n%s\n%s\n%d", method, url, query, hashedPayload, timestamp)
	mac := hmac.New(sha512.New, []byte(self.Secret))
	mac.Write([]byte(s))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]interface{}{
		"KEY":       self.ApiKey,
		"Timestamp": fmt.Sprint(timestamp),
		"SIGN":      sign,
	}
}

func (self *Gateio) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	url := self.Urls["api"].(map[string]interface{})[api].(string) + "/" + self.ImplodeParams(path, params)
	query := self.Omit(params, self.ExtractParams(path))
	if api == "public" {
		if len(query) > 0 {
			url += "?" + self.Urlencode(query)
		}
	} else {
		self.CheckRequiredCredentials()
		u, _ := urllib.Parse(url)
		if method == "GET" || method == "DELETE" {
			queryString := self.Urlencode(query)
			headers = self.genSign(method, u.Path, queryString, "")
			if len(query) > 0 {
				url += "?" + self.Urlencode(query)
			}
		} else {
			bodyString := self.Json(query)
			headers = self.genSign(method, u.Path, "", bodyString)
		}
		headers.(map[string]interface{})["Content-Type"] = "application/json"
	}

	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}
