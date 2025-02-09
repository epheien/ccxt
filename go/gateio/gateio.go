package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	. "github.com/epheien/ccxt/go/base"
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
                "spot/order_book",
                "spot/currencies",
                "spot/currency_pairs",
                "spot/trades",
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

func (self *Gateio) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	response := self.ApiFuncReturnList("publicGetSpotCurrencyPairs", params, nil, nil)
	data := response
	result := []interface{}{}
	for i := 0; i < self.Length(data); i++ {
		/*
			{
				"amount_precision": 0,
				"base": "100X",
				"buy_start": 1622793600,
				"fee": "0.2",
				"id": "100X_USDT",
				"min_quote_amount": "1",
				"precision": 11,
				"quote": "USDT",
				"sell_start": 1608782400,
				"trade_status": "untradable"
			}
		*/
		market := self.Member(data, i)
		id := self.SafeString(market, "id", "")
		baseId, quoteId := self.Unpack2(strings.Split(id, "_"))
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		active := (self.SafeString(market, "trade_status") == "tradable")
		precision := map[string]interface{}{
			"amount": self.SafeInteger(market, "amount_precision"),
			"price":  self.SafeInteger(market, "precision"),
		}
		limits := map[string]interface{}{
			"cost": map[string]interface{}{
				"min": self.SafeFloat(market, "min_quote_amount"),
			},
		}
		result = append(result, map[string]interface{}{
			"id":        id,
			"symbol":    symbol,
			"baseId":    baseId,
			"quoteId":   quoteId,
			"base":      base,
			"quote":     quote,
			"active":    active,
			"precision": precision,
			"limits":    limits,
			"info":      market,
		})
	}
	return self.ToMarkets(result), nil
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

func (self *Gateio) CreateOrder(symbol string, _type string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if _type != "limit" {
		self.RaiseException("ExchangeError", self.Id+" allows limit orders only")
	}
	marketId := self.MarketId(symbol)
	request := map[string]interface{}{
		"account":       self.Options["account"],
		"currency_pair": marketId,
		"side":          side,
		"price":         self.Float64ToString(price),
		"amount":        self.Float64ToString(amount),
	}
	response := self.ApiFunc("privatePostSpotOrders", self.Extend(request, params), nil, nil)
	data := response
	timestamp := self.SafeInteger(response, "create_time_ms")
	order := map[string]interface{}{
		"id":        self.SafeString(data, "id"),
		"symbol":    symbol,
		"type":      _type,
		"side":      side,
		"price":     price,
		"amount":    amount,
		"cost":      nil,
		"filled":    nil,
		"remaining": nil,
		"timestamp": timestamp,
		"datetime":  self.Iso8601(timestamp),
		"fee":       nil,
		"status":    "open",
		"info":      data,
	}
	return self.ToOrder(order), nil
}

func (self *Gateio) ParseOrderStatus(status string) string {
	// NOTE: 类型必须为 map[string]interface{}, 否则无法使用 SafeString
	statuses := map[string]interface{}{
		"open":      "open",
		"closed":    "closed",
		"cancelled": "canceled",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Gateio) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	var symbol string
	if market != nil {
		symbol = market.(*Market).Symbol
	}
	orderId := self.SafeString(order, "id")
	timestamp := self.SafeInteger(order, "create_time_ms")
	status := self.ParseOrderStatus(self.SafeString(order, "status"))
	side := self.SafeString(order, "side")
	price := self.SafeFloat(order, "price")
	amount := self.SafeFloat(order, "amount")
	remaining := self.SafeFloat(order, "left")
	filled := amount - remaining

	return map[string]interface{}{
		"id":                 orderId,
		"symbol":             symbol,
		"side":               side,
		"amount":             amount,
		"price":              price,
		"filled":             filled,
		"remaining":          remaining,
		"timestamp":          timestamp,
		"datetime":           self.Iso8601(timestamp),
		"status":             status,
		"info":               order,
		"lastTradeTimestamp": nil,
		"average":            nil,
		"trades":             nil,
	}
}

func (self *Gateio) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if symbol == "" {
		self.RaiseException("ArgumentsRequired", "symbol")
	}
	market := self.Market(symbol)
	request := map[string]interface{}{
		"currency_pair": market.Id,
		"status":        "open",
		"limit":         100,
	}
	response := self.ApiFuncReturnList("privateGetSpotOrders", self.Extend(request, params), nil, nil)
	orders := []interface{}{}
	at := self.Options["account"].(string)
	for _, one := range response {
		if self.SafeString(one, "account") != at {
			continue
		}
		orders = append(orders, one)
	}

	return self.ToOrders(self.ParseOrders(orders, market, since, limit)), nil
}

func (self *Gateio) FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) (trades []*Trade, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"currency_pair": market.Id,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	if since > 0 {
		request["from"] = since
	}
	response := self.ApiFuncReturnList("publicGetSpotTrades", self.Extend(request, params), nil, nil)
	trades = self.ParseTrades(response, market, since, limit)
	trades = self.ReverseTrades(trades)
	return
}

func (self *Gateio) ParseTrade(trade interface{}, market *Market) (result *Trade) {
	result = &Trade{
		Id:        self.SafeString(trade, "id"),
		Timestamp: self.SafeInteger(trade, "create_time_ms"),
		Price:     self.SafeFloat(trade, "price"),
		Amount:    self.SafeFloat(trade, "amount"),
		Side:      self.SafeString(trade, "side"),
		Info:      trade,
	}
	result.Datetime = self.Iso8601(result.Timestamp)
	if market != nil {
		result.Symbol = market.Symbol
	}
	return
}

func (self *Gateio) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if symbol == "" {
		self.RaiseException("ArgumentsRequired", "symbol")
	}
	market := self.Market(symbol)
	request := map[string]interface{}{
		"order_id":      id,
		"currency_pair": market.Id,
	}
	response := self.ApiFunc("privateGetSpotOrdersOrderId", self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Gateio) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if symbol == "" {
		self.RaiseException("ArgumentsRequired", "symbol")
	}

	market := self.Market(symbol)
	request := map[string]interface{}{
		"order_id":      id,
		"currency_pair": market.Id,
	}
	// NOTE: 撤掉的返回类型有时候是 []interface{} 有时候是 map[string]interface{}, 暂时不管
	response = self.ApiFuncRaw("privateDeleteSpotOrdersOrderId", self.Extend(request, params).(map[string]interface{}), nil, nil)
	return response, nil
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
			body = self.Json(query)
			headers = self.genSign(method, u.Path, "", body.(string))
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

func (self *Gateio) HandleErrors(
	code int64, reason string, url string, method string, headers interface{}, body string, response interface{},
	requestHeaders interface{}, requestBody interface{},
) {
	if response == nil {
		return
	}
	if _, ok := response.(map[string]interface{}); !ok {
		return
	}
	errorCode := self.SafeString(response, "label")
	message := self.SafeString(response, "message")
	self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), errorCode, message)
}
