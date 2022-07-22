package futures_kucoin

import (
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"strings"
)

type FuturesKucoin struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *FuturesKucoin, err error) {
	ex = new(FuturesKucoin)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *FuturesKucoin) Describe() []byte {
	return []byte(`{
    "id": "futures_kumex",
    "name": "KuMEX",
    "countries": ["SC"],
    "version": "v1",
    "has": {
        "fetchMarkets": true,
        "fetchCurrencies": true,
        "fetchTicker": true,
        "fetchTickers": true,
        "fetchOrderBook": true,
        "fetchOrder": true,
        "fetchClosedOrders": true,
        "fetchOpenOrders": true,
        "fetchDepositAddress": true,
        "createDepositAddress": true,
        "withdraw": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchBalance": true,
        "fetchTrades": true,
        "fetchMyTrades": true,
        "createOrder": true,
        "cancelOrder": true,
        "fetchAccounts": true,
        "fetchFundingFee": true,
        "fetchOHLCV": true,
    },
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/51909432-b0a72780-23dd-11e9-99ba-73d23c8d4eed.jpg",
        "api": {
            "public": "https://api-futures.kucoin.com",
            "private": "https://api-futures.kucoin.com",
        },
        "test": {
            "public": "https://sandbox-api.kumex.com",
            "private": "https://sandbox-api.kumex.com",
        },
        "www": "https://www.kumex.com",
        "doc": [
            "https://docs.kumex.com",
        ],
    },
    "requiredCredentials": {
        "apiKey": true,
        "secret": true,
        "password": true,
    },
    "api": {
        "public": {
            "get": [
                "contracts/active",
                "contracts/{symbol}",
                "ticker",
                "level2/snapshot",
                "level2/depth20",
                "level2/depth100",
                "level2/message/query",
                "level3/snapshot",
                "v2/level3/snapshot",
                "mark-price/{symbol}/current",
                "funding-rate/{symbol}/current",
            ],
            "post": [
                "bullet-public",
            ],
        },
        "private": {
            "get": [
                "account-overview",
                "transaction-history",
                "orders",
                "stopOrders",
                "recentDoneOrders",
                "orders/{orderId}",
                "position",
                "positions",
                "funding-history",
                "transfer-list",
            ],
            "post": [
                "orders",
                "position/margin/auto-deposit-status",
                "position/margin/deposit-margin",
                "bullet-private",
                "transfer-out",
            ],
            "delete": [
                "orders/{orderId}",
                "orders",
                "stopOrders",
                "cancel/transfer-out",
            ],
        },
    },
    "timeframes": {
        "1m": "1min",
        "3m": "3min",
        "5m": "5min",
        "15m": "15min",
        "30m": "30min",
        "1h": "1hour",
        "2h": "2hour",
        "4h": "4hour",
        "6h": "6hour",
        "8h": "8hour",
        "12h": "12hour",
        "1d": "1day",
        "1w": "1week",
    },
    "exceptions": {
        "400": "BadRequest",
        "401": "AuthenticationError",
        "403": "NotSupported",
        "404": "NotSupported",
        "405": "NotSupported",
        "429": "DDoSProtection",
        "500": "ExchangeError",
        "503": "ExchangeNotAvailable",
        "200004": "InsufficientFunds",
        "300000": "InvalidOrder",
        "300009": "InsufficientFunds",
        "400001": "AuthenticationError",
        "400002": "InvalidNonce",
        "400003": "AuthenticationError",
        "400004": "AuthenticationError",
        "400005": "AuthenticationError",
        "400006": "AuthenticationError",
        "400007": "AuthenticationError",
        "400008": "NotSupported",
        "400100": "ArgumentsRequired",
        "411100": "AccountSuspended",
        "500000": "ExchangeError",
        "order_not_exist": "OrderNotFound",  
    },
    "fees": {
        "trading": {
            "tierBased": false,
            "percentage": true,
            "taker": 0.001,
            "maker": 0.001,
        },
        "funding": {
            "tierBased": false,
            "percentage": false,
            "withdraw": {},
            "deposit": {},
        },
    },
    "options": {
        "version": "v1",
        "symbolSeparator": "-",
    },
    "markets_by_id": {},
}`)
}

func (self *FuturesKucoin) LoadMarkets() map[string]*Market {
	return nil
}

func (self *FuturesKucoin) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	return &Market{
		Id:     li[0] + li[1],
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *FuturesKucoin) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()

	if params != nil && params["symbol"] != nil {
		params["currency"] = params["symbol"]
		if params["currency"] == "USDTM" {
			params["currency"] = "USDT"
		}
		delete(params, "symbol")
	}

	response := self.ApiFunc("privateGetAccountOverview", params, nil, nil)
	responseData := response["data"]

	result := map[string]interface{}{
		"info": response,
	}
	entry := responseData
	account := self.Account()
	code := self.SafeString(entry, "currency", "")
	if code == "USDT" {
		code = "USDTM"
	}
	account["total"] = self.SafeFloat(entry, "marginBalance", 0)
	account["free"] = self.SafeFloat(entry, "availableBalance", 0)
	account["used"] = account["total"].(float64) - account["free"].(float64)
	account["realPnl"] = 0.0
	account["unrealPnl"] = self.SafeFloat(entry, "unrealisedPNL", 0.0)
	result[code] = account

	return self.ParseBalance(result), nil
}

func (self *FuturesKucoin) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	response := self.ApiFunc("publicGetLevel2Depth20", self.Extend(request, params), nil, nil)
	orderBook = self.ParseOrderBook(response["data"], 0, "bids", "asks", 0, 1)
	orderBook.Timestamp = int64(response["data"].(map[string]interface{})["ts"].(float64) / 1000000)
	orderBook.Datetime = self.Iso8601(orderBook.Timestamp)
	return orderBook, nil
}

func (self *FuturesKucoin) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"NEW":              "open",
		"PARTIALLY_FILLED": "open",
		"FILLED":           "closed",
		"CANCELED":         "canceled",
		"PENDING_CANCEL":   "canceling",
		"REJECTED":         "rejected",
		"EXPIRED":          "canceled",
	}
	return self.SafeString(statuses, status, status)
}

func (self *FuturesKucoin) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	symbol := ""
	if market != nil {
		symbol = market.(*Market).Symbol
	}
	clientOid := self.SafeString(order, "clientOid", "")
	orderId := self.SafeString(order, "id", "")
	_type := self.SafeString(order, "type", "")
	timestamp := self.SafeInteger(order, "createdAt", 0)
	datetime := self.Iso8601(timestamp)
	price := self.SafeFloat(order, "price", 0)
	side := self.SafeString(order, "side", "")
	amount := self.SafeFloat(order, "size", 0)
	filled := self.SafeFloat(order, "dealSize", 0)
	status := "closed"
	if order.(map[string]interface{})["isActive"].(bool) {
		status = "open"
	}
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

func (self *FuturesKucoin) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	clientOid := self.Uuid()
	request := map[string]interface{}{
		"clientOid": clientOid,
		"price":     fmt.Sprint(price),
		"side":      side,
		"size":      fmt.Sprint(amount),
		"symbol":    market.Id,
		"type":      typ,
		"leverage":  5,
	}
	response := self.ApiFunc("privatePostOrders", self.Extend(request, params), nil, nil)
	responseData := response["data"]
	return &Order{
		Id:            responseData.(map[string]interface{})["orderId"].(string),
		Symbol:        symbol,
		Type:          typ,
		Side:          side,
		Status:        "open",
		ClientOrderId: clientOid,
		Info:          responseData,
	}, nil
}

func (self *FuturesKucoin) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"orderId": id,
	}
	response := self.ApiFunc("privateGetOrdersOrderId", self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response["data"], market)), nil
}

func (self *FuturesKucoin) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{
		"status": "active",
	}
	var market *Market
	if symbol != "" {
		market = self.Market(symbol)
		request["symbol"] = market.Id
	}
	if since > 0 {
		request["startAt"] = since
		if limit > 0 {
			request["endAt"] = since + limit
		}
	}
	response := self.ApiFunc("privateGetOrders", self.Extend(request, params), nil, nil)
	return self.ToOrders(self.ParseOrders(response["data"].(map[string]interface{})["items"], market, since, limit)), nil
}

func (self *FuturesKucoin) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{
		"orderId": id,
	}
	response = self.ApiFunc("privateDeleteOrdersOrderId", self.Extend(request, params), nil, nil)
	return response, nil
}

func (self *FuturesKucoin) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	//
	// the v2 URL is https://openapi-v2.kucoin.com/api/v1/endpoint
	//                                †                 ↑
	//
	var endpoint string
	if strings.HasPrefix(path, "v2/") {
		endpoint = "/api/" + self.ImplodeParams(path, params)
	} else {
		endpoint = "/api/" + self.Options["version"].(string) + "/" + self.ImplodeParams(path, params)
	}
	query := self.Omit(params, self.ExtractParams(path))
	endpart := ""
	if len(headers.(map[string]interface{})) == 0 {
		headers = map[string]interface{}{}
	}
	if len(query) > 0 {
		if method == "GET" {
			endpoint += "?" + self.Urlencode(query)
		} else {
			endpart = self.Json(query)
			self.SetValue(headers, "Content-Type", "application/json")
		}
	}
	url := self.Member(self.Urls["api"], api).(string) + endpoint
	if api == "private" {
		self.CheckRequiredCredentials()
		timestamp := fmt.Sprintf("%d", self.Nonce())
		headers = self.Extend(map[string]interface{}{
			"KC-API-KEY":        self.ApiKey,
			"KC-API-TIMESTAMP":  timestamp,
			"KC-API-PASSPHRASE": self.Password,
		}, headers)
		payload := timestamp + method + endpoint + endpart
		signature := self.Hmac(self.Encode(payload), self.Encode(self.Secret), "sha256", "base64")
		self.SetValue(headers, "KC-API-SIGN", self.Decode(signature))
		// v2 apiKey
		headers.(map[string]interface{})["KC-API-KEY-VERSION"] = "2"
		password := self.Hmac(self.Encode(self.Password), self.Encode(self.Secret), "sha256", "base64")
		headers.(map[string]interface{})["KC-API-PASSPHRASE"] = self.Decode(password)
	}
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    endpart,
		"headers": headers,
	}
}

func (self *FuturesKucoin) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if response == nil {
		self.ThrowBroadlyMatchedException(self.Member(self.Exceptions, "broad"), body, body)
		return
	}
	errorCode := self.SafeString(response, "code", "")
	message := self.SafeString(response, "msg", "")
	self.ThrowExactlyMatchedException(self.Exceptions, errorCode, message)
	self.ThrowExactlyMatchedException(self.Exceptions, message, message)
	if errorCode != "200000" {
		self.RaiseException("ExchangeError", fmt.Sprintf("%s %s", self.Id, body))
	}
}
