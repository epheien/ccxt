package futures_binance

import (
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"math"
	"strings"
)

type FuturesBinance struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *FuturesBinance, err error) {
	ex = new(FuturesBinance)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *FuturesBinance) Describe() []byte {
	return []byte(`
{
    "id": "futures_binance",
    "name": "Binance",
    "countries": "JP",
    "rateLimit": 500,
    "has": {
        "CORS": false,
        "fetchBidsAsks": true,
        "fetchTickers": true,
        "fetchOHLCV": true,
        "fetchMyTrades": true,
        "fetchOrder": true,
        "fetchOrders": true,
        "fetchOpenOrders": true,
        "fetchClosedOrders": true
    },
    "timeframes": {
        "1m": "1m",
        "3m": "3m",
        "5m": "5m",
        "15m": "15m",
        "30m": "30m",
        "1h": "1h",
        "2h": "2h",
        "4h": "4h",
        "6h": "6h",
        "8h": "8h",
        "12h": "12h",
        "1d": "1d",
        "3d": "3d",
        "1w": "1w",
        "1M": "1M"
    },
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/29604020-d5483cdc-87ee-11e7-94c7-d1a8d9169293.jpg",
        "api": {
            "public": "https://fapi.binance.com/fapi/v1",
            "private": "https://fapi.binance.com/fapi/v1"
        },
        "www": "https://www.binance.com",
        "doc": "https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md",
        "fees": [
            "https://binance.zendesk.com/hc/en-us/articles/115000429332",
            "https://support.binance.com/hc/en-us/articles/115000583311"
        ]
    },
    "api": {
        "public": {
            "get": [
                "ping",
                "time",
                "exchangeInfo",
                "depth",
                "trades",
                "historicalTrades",
                "aggTrades",
                "klines",
                "premiumIndex",
                "fundingRate",
                "ticker/24hr",
                "ticker/price",
                "ticker/bookTicker"
            ]
        },
        "private": {
            "get": [
                "order",
                "openOrders",
                "allOrders",
                "balance",
                "account",
                "positionRisk",
                "userTrades",
                "income"
            ],
            "post": [
                "order",
                "order/test",
                "leverage",
                "positionRisk",
                "listenKey"
            ],
            "delete": [
                "order",
                "listenKey"
            ],
            "put": [
                "listenKey"
            ]
        }
    },
    "fees": {
        "trading": {
            "tierBased": false,
            "percentage": true,
            "taker": 0.001,
            "maker": 0.001
        }
    },
    "options": {
        "warnOnFetchOpenOrdersWithoutSymbol": true,
        "recvWindow": 5000,
        "timeDifference": 0,
        "adjustForTimeDifference": false
    },
    "exceptions": {
		"exact": {
			"-1121": "InvalidSymbol",
			"-1013": "InvalidOrder",
			"-1021": "InvalidNonce",
			"-1100": "InvalidOrder",
			"-2010": "InsufficientFunds",
			"-2011": "CancelRejected",
			"-2013": "OrderNotFound",
			"-2015": "AuthenticationError"
		},
		"broad": {
			"Price * QTY is zero or less": "InvalidOrder",
			"LOT_SIZE": "InvalidOrder",
			"PRICE_FILTER": "InvalidOrder",
			"Order does not exist": "OrderNotFound",
		},
    }
}
`)
}

func (self *FuturesBinance) LoadMarkets() map[string]*Market {
	return nil
}

func (self *FuturesBinance) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	if len(li) != 2 {
		return &Market{}
	}
	return &Market{
		Id:     li[0] + li[1],
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *FuturesBinance) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()

	response := self.ApiFunc("privateGetAccount", params, nil, nil)

	result := map[string]interface{}{
		"info": response,
	}

	balances := response["assets"].([]interface{})
	for _, balance := range balances {
		account := self.Account()
		total := self.SafeFloat(balance, "marginBalance")
		used := self.SafeFloat(balance, "maintMargin") + self.SafeFloat(balance, "initialMargin")
		account["total"] = total
		account["used"] = used
		account["free"] = total - used
		account["unrealPnl"] = self.SafeFloat(balance, "unrealizedProfit")
		currency := self.SafeString(balance, "asset")
		result[currency] = account
	}

	return self.ParseBalance(result), nil
}

func (self *FuturesBinance) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFunc("publicGetDepth", self.Extend(request, params), nil, nil)
	orderBook = self.ParseOrderBook(response, ToInteger(response["T"]), "bids", "asks", 0, 1)
	return orderBook, nil
}

func (self *FuturesBinance) ParseTicker(response interface{}) (ticker *Ticker) {
	timestamp := self.SafeInteger(response, "closeTime")
	datetime := self.Iso8601(timestamp)
	last := self.SafeFloat(response, "lastPrice")
	ticker = &Ticker{
		Timestamp:   timestamp,
		Datetime:    datetime,
		Last:        last,
		Open:        self.SafeFloat(response, "openPrice"),
		High:        self.SafeFloat(response, "highPrice"),
		Low:         self.SafeFloat(response, "lowPrice"),
		Close:       last,
		BaseVolume:  self.SafeFloat(response, "volume"),
		QuoteVolume: self.SafeFloat(response, "quoteVolume"),
		Change:      self.SafeFloat(response, "priceChange"),
		Percentage:  self.SafeFloat(response, "priceChangePercent"),
		Vwap:        self.SafeFloat(response, "weightedAvgPrice"),
		Info:        response,
	}
	return
}

func (self *FuturesBinance) FetchTicker(symbol string, params map[string]interface{}) (ticker *Ticker, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": self.Member(market, "id"),
	}
	response := self.ApiFunc("publicGetTicker24hr", self.Extend(request, params), nil, nil)
	ticker = self.ParseTicker(response)
	ticker.Symbol = symbol
	return ticker, nil
}

func (self *FuturesBinance) ParseOHLCV(response interface{}) *OHLCV {
	data := response.([]interface{})
	return &OHLCV{
		Timestamp: ToInteger(data[0]),
		Open:      ToFloat(data[1]),
		High:      ToFloat(data[2]),
		Low:       ToFloat(data[3]),
		Close:     ToFloat(data[4]),
		Volume:    ToFloat(data[5]),
		Info:      response,
	}
}

func (self *FuturesBinance) FetchOHLCV(symbol, timeframe string, since int64, limit int64, params map[string]interface{}) (klines []*OHLCV, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol":   market.Id,
		"interval": self.Timeframes[timeframe],
	}
	if since > 0 {
		request["startTime"] = since
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFuncReturnList("publicGetKlines", self.Extend(request, params), nil, nil)
	for _, item := range response {
		klines = append(klines, self.ParseOHLCV(item))
	}
	return klines, nil
}

func (self *FuturesBinance) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"NEW":              "open",
		"PARTIALLY_FILLED": "open",
		"FILLED":           "closed",
		"CANCELED":         "canceled",
		"REJECTED":         "canceled",
		"EXPIRED":          "canceled", // 订单过期(根据timeInForce参数规则)
	}
	return self.SafeString(statuses, status, status)
}

func (self *FuturesBinance) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	symbol := ""
	if market != nil {
		symbol = market.(*Market).Symbol
	}
	clientOid := self.SafeString(order, "clientOrderId", "")
	orderId := fmt.Sprintf("%d", self.SafeInteger(order, "orderId"))
	_type := self.SafeStringLower(order, "type", "")
	timestamp := self.SafeInteger(order, "time", 0)
	datetime := self.Iso8601(timestamp)
	price := self.SafeFloat(order, "price", 0)
	side := self.SafeStringLower(order, "side", "")
	amount := self.SafeFloat(order, "origQty", 0)
	filled := self.SafeFloat(order, "executedQty", 0)
	status := self.ParseOrderStatus(self.SafeString(order, "status"))
	return map[string]interface{}{
		"clientOrderId": clientOid,
		"id":            orderId,
		"symbol":        symbol,
		"type":          _type,
		"side":          side,
		"amount":        amount,
		"price":         price,
		"filled":        filled,
		"remaining":     amount - filled,
		"timestamp":     timestamp,
		"datetime":      datetime,
		"status":        status,
		"info":          order,
	}
}

func (self *FuturesBinance) CreateOrder(symbol string, type_ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol":   market.Id,
		"quantity": amount,
		"side":     strings.ToUpper(side),
		"type":     strings.ToUpper(type_),
	}
	if type_ == "limit" {
		request["price"] = price
		request["timeInForce"] = "GTC"
	}
	response := self.ApiFunc("privatePostOrder", self.Extend(request, params), nil, nil)
	return &Order{
		Id:            fmt.Sprintf("%d", self.SafeInteger(response, "orderId")),
		Symbol:        symbol,
		Type:          type_,
		Side:          side,
		Status:        "open",
		ClientOrderId: self.SafeString(response, "clientOrderId"),
		Info:          response,
	}, nil
}

func (self *FuturesBinance) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	if id != "" {
		request["orderId"] = id
	}
	response := self.ApiFunc("privateGetOrder", self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *FuturesBinance) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{}
	var market *Market
	if symbol != "" {
		market = self.Market(symbol)
		request["symbol"] = market.Id
	}
	response := self.ApiFuncReturnList("privateGetOpenOrders", self.Extend(request, params), nil, nil)
	return self.ToOrders(self.ParseOrders(response, market, since, limit)), nil
}

func (self *FuturesBinance) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	if id != "" {
		request["orderId"] = id
	}
	response = self.ApiFunc("privateDeleteOrder", self.Extend(request, params), nil, nil)
	return response, nil
}

func (self *FuturesBinance) FetchMarkPrice(symbol string, params map[string]interface{}) (markPrice *MarkPrice, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	response := self.ApiFunc("publicGetPremiumIndex", self.Extend(request, params), nil, nil)
	data := response
	return &MarkPrice{
		Symbol:     symbol,
		MarkPrice:  self.SafeFloat(data, "markPrice", 0),
		IndexPrice: self.SafeFloat(data, "indexPrice", 0),
		Timestamp:  self.SafeInteger(data, "time", 0),
		Info:       data,
	}, nil
}

func (self *FuturesBinance) FetchPosition(symbol string, params map[string]interface{}) (result []*Position, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{}
	market := self.Market(symbol)
	request["symbol"] = market.Id
	response := self.ApiFuncReturnList("privateGetPositionRisk", self.Extend(request, params), nil, nil)

	for _, item := range response {
		if self.SafeString(item, "symbol") != market.Id {
			continue
		}
		amount := self.SafeFloat(item, "positionAmt")
		pos := &Position{
			Symbol:     symbol,
			Side:       "long",
			Leverage:   self.SafeFloat(item, "leverage", 0),
			Amount:     math.Abs(amount),
			UsedAmount: 0,
			Price:      self.SafeFloat(item, "entryPrice", 0),
			RealPnl:    0,
			UnrealPnl:  self.SafeFloat(item, "unRealizedProfit", 0),
			Info:       item,
		}
		if amount < 0 {
			pos.Side = "short"
		}
		result = append(result, pos)
	}
	return
}

func (self *FuturesBinance) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	url := self.Urls["api"].(map[string]interface{})[api].(string) + "/" + path
	if path == "userDataStream" {
		body = self.Urlencode(params)
		headers = map[string]interface{}{
			"X-MBX-APIKEY": self.ApiKey,
			"Content-Type": "application/x-www-form-urlencoded",
		}
	} else if api == "private" {
		self.CheckRequiredCredentials()
		query := self.Urlencode(self.Extend(map[string]interface{}{
			"timestamp":  self.Nonce(),
			"recvWindow": self.Options["recvWindow"],
		}, params))
		signature := self.Hmac(query, self.Secret, "sha256", "hex")
		query += fmt.Sprintf("&signature=%s", signature)
		headers = map[string]interface{}{
			"X-MBX-APIKEY": self.ApiKey,
		}
		if method == "GET" {
			url += "?" + query
		} else {
			body = query
			headers.(map[string]interface{})["Content-Type"] = "application/x-www-form-urlencoded"
		}
	} else {
		// api == "public"
		if len(params) > 0 {
			url += "?" + self.Urlencode(params)
		}
	}
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *FuturesBinance) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if httpCode < 300 {
		return
	}
	if httpCode < 400 {
		return
	}
	if httpCode == 418 || httpCode == 429 {
		self.RaiseException("DDoSProtection", fmt.Sprintf("%s %d %s %s", self.Id, httpCode, reason, body))
	}

	if httpCode >= 400 {
		self.ThrowExactlyMatchedException(self.Exceptions["broad"], body, body)
	}

	code := self.SafeInteger(response, "code")
	if code == 0 {
		return
	}
	msg := self.SafeString(response, "msg")
	self.ThrowExactlyMatchedException(self.Exceptions["exact"], fmt.Sprintf("%d", code), self.Id+" "+msg)
	self.ThrowBroadlyMatchedException(self.Exceptions["broad"], msg, self.Id+" "+msg)

	self.RaiseException("ExchangeError", fmt.Sprintf("%s %s", self.Id, body))
}
