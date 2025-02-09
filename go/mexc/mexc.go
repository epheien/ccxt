package mexc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	. "github.com/epheien/ccxt/go/base"
	"strings"
)

type Mexc struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Mexc, err error) {
	ex = new(Mexc)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Mexc) Describe() []byte {
	return []byte(`
{
    "id": "mexc",
    "name": "mexc",
    "countries": [
        "CN"
    ],
    "version": "3",
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
        "logo": "https://www.mexc.com/images/full-logo-light-ko.svg",
        "api": {
            "public": "https://api.mexc.com",
            "private": "https://api.mexc.com"
        },
        "www": "https://www.mexc.com",
        "doc": "https://mxcdevelop.github.io"
    },
    "api": {
        "public": {
            "get": [
                "ping",
                "time",
                "defaultSymbols",
                "exchangeInfo",
                "depth",
            ]
        },
        "private": {
            "get": [
                "account",
                "order",
                "openOrders",
                "allOrders",
                "myTrades",
                "mxDeduct/enable",
            ],
            "post": [
                "order",
                "batchOrders",
                "mxDeduct/enable",
            ],
            "delete": [
                "order",
                "openOrders",
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

func (self *Mexc) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	return &Market{
		Id:     li[0] + li[1], // 需要大写
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *Mexc) LoadMarkets() map[string]*Market {
	return nil
}

func (self *Mexc) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	response := self.ApiFunc("publicGetExchangeInfo", params, nil, nil)
	data := response["symbols"].([]interface{})
	result := []interface{}{}
	for _, market := range data {
		/*
		   {
		       "baseAsset": "AES",
		       "baseAssetPrecision": 2,
		       "baseCommissionPrecision": 2,
		       "baseSizePrecision": "0",
		       "filters": [],
		       "isMarginTradingAllowed": false,
		       "isSpotTradingAllowed": true,
		       "makerCommission": "0",
		       "maxQuoteAmount": "5000000.000000000000000000",
		       "maxQuoteAmountMarket": "500000.000000000000000000",
		       "orderTypes": [
		           "LIMIT",
		           "MARKET",
		           "LIMIT_MAKER"
		       ],
		       "permissions": [
		           "SPOT"
		       ],
		       "quoteAmountPrecision": "5.000000000000000000",
		       "quoteAmountPrecisionMarket": "5.000000000000000000",
		       "quoteAsset": "USDT",
		       "quoteAssetPrecision": 6,
		       "quoteCommissionPrecision": 6,
		       "quotePrecision": 6,
		       "status": "ENABLED",
		       "symbol": "AESUSDT",
		       "takerCommission": "0"
		   },
		*/
		id := self.SafeString(market, "symbol", "")
		baseId := self.SafeString(market, "baseAsset")
		quoteId := self.SafeString(market, "quoteAsset")
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		active := (self.SafeString(market, "status") == "ENABLED")
		precision := map[string]interface{}{
			"amount": self.SafeInteger(market, "baseAssetPrecision"),
			"price":  self.SafeInteger(market, "quotePrecision"),
		}
		limits := map[string]interface{}{
			"amount": map[string]interface{}{
				"min": self.SafeFloat(market, "baseSizePrecision"),
			},
			"cost": map[string]interface{}{
				"min": self.SafeFloat(market, "quoteAmountPrecision"),
				"max": self.SafeFloat(market, "maxQuoteAmount"),
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
			"spot":      true,
		})
	}
	return self.ToMarkets(result), nil
}

func (self *Mexc) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	marketId := self.MarketId(symbol)
	request := map[string]interface{}{
		"symbol": marketId,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFunc("publicGetDepth", self.Extend(request, params), nil, nil)
	orderbook := self.ParseOrderBook(response, 0, "bids", "asks", 0, 1)
	return orderbook, nil
}

func (self *Mexc) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	response := self.ApiFunc("privateGetAccount", params, nil, nil)
	result := map[string]interface{}{
		"info": response,
	}
	for _, one := range response["balances"].([]interface{}) {
		account := self.Account()
		free := self.SafeFloat(one, "free")
		used := self.SafeFloat(one, "locked")
		cc := self.SafeString(one, "asset")
		account["free"] = free
		account["used"] = used
		account["total"] = free + used
		result[cc] = account
	}
	return self.ParseBalance(result), nil
}

func (self *Mexc) CreateOrder(symbol string, _type string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
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
		"type":     strings.ToUpper(_type),
		"symbol":   marketId,
		"side":     strings.ToUpper(side),
		"price":    self.Float64ToString(price),
		"quantity": self.Float64ToString(amount),
	}
	response := self.ApiFunc("privatePostOrder", self.Extend(request, params), nil, nil)
	data := response
	timestamp := self.SafeInteger(response, "transactTime")
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

func (self *Mexc) ParseOrderStatus(status string) string {
	// NOTE: 类型必须为 map[string]interface{}, 否则无法使用 SafeString
	statuses := map[string]interface{}{
		"open":      "open",
		"closed":    "closed",
		"cancelled": "canceled",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Mexc) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
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

func (self *Mexc) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
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

func (self *Mexc) FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) (trades []*Trade, err error) {
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

func (self *Mexc) ParseTrade(trade interface{}, market *Market) (result *Trade) {
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

func (self *Mexc) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
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

func (self *Mexc) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
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
	response = self.ApiFunc("privateDeleteSpotOrdersOrderId", self.Extend(request, params), nil, nil)
	return response, nil
}

func (self *Mexc) genSign(query string, timestamp int64) (string, string) {
	var payload string
	if query == "" {
		payload = fmt.Sprintf("timestamp=%d", timestamp)
	} else {
		payload = fmt.Sprintf("%s&timestamp=%d", query, timestamp)
	}
	mac := hmac.New(sha256.New, []byte(self.Secret))
	fmt.Println(payload)
	mac.Write([]byte(payload))
	sign := hex.EncodeToString(mac.Sum(nil))
	return payload, sign
}

func (self *Mexc) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	url := self.Urls["api"].(map[string]interface{})[api].(string) + "/api/v3/" + self.ImplodeParams(path, params)
	query := self.Omit(params, self.ExtractParams(path))
	if api == "public" {
		if len(query) > 0 {
			url += "?" + self.Urlencode(query)
		}
	} else {
		self.CheckRequiredCredentials()
		if method == "GET" {
			queryString := self.Urlencode(query)
			queryString, sign := self.genSign(queryString, self.Milliseconds())
			url += "?" + queryString + fmt.Sprintf("&signature=%s", sign)
		} else {
			queryString := self.Urlencode(query)
			queryString, sign := self.genSign(queryString, self.Milliseconds())
			queryString += fmt.Sprintf("&signature=%s", sign)
			body = queryString
		}
		headers = map[string]interface{}{
			"X-MEXC-APIKEY": self.ApiKey,
			"Content-Type":  "application/json",
		}
	}

	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Mexc) HandleErrors(
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
