package bybit

import (
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"strings"
)

type Bybit struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Bybit, err error) {
	ex = new(Bybit)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Bybit) Describe() []byte {
	return []byte(`{
    "id": "bybit",
    "name": "Bybit",
    "countries": [
        "VG"
    ],
    "version": "v3",
    "userAgent": null,
    "rateLimit": 20,
    "hostname": "bybit.com",
    "pro": true,
    "certified": true,
    "has": {
        "CORS": true,
        "spot": true,
        "margin": true,
        "swap": true,
        "future": true,
        "option": null,
        "cancelAllOrders": true,
        "cancelOrder": true,
        "createOrder": true,
        "createStopLimitOrder": true,
        "createStopMarketOrder": true,
        "createStopOrder": true,
        "editOrder": true,
        "fetchBalance": true,
        "fetchBorrowInterest": true,
        "fetchBorrowRate": false,
        "fetchBorrowRateHistories": false,
        "fetchBorrowRateHistory": false,
        "fetchBorrowRates": false,
        "fetchClosedOrders": true,
        "fetchCurrencies": true,
        "fetchDepositAddress": true,
        "fetchDepositAddresses": false,
        "fetchDepositAddressesByNetwork": true,
        "fetchDeposits": true,
        "fetchFundingRate": true,
        "fetchFundingRateHistory": false,
        "fetchIndexOHLCV": true,
        "fetchLedger": true,
        "fetchMarketLeverageTiers": true,
        "fetchMarkets": true,
        "fetchMarkOHLCV": true,
        "fetchMyTrades": true,
        "fetchOHLCV": true,
        "fetchOpenInterestHistory": true,
        "fetchOpenOrders": true,
        "fetchOrder": true,
        "fetchOrderBook": true,
        "fetchOrders": true,
        "fetchOrderTrades": true,
        "fetchPositions": true,
        "fetchPremiumIndexOHLCV": true,
        "fetchTicker": true,
        "fetchTickers": true,
        "fetchTime": true,
        "fetchTrades": true,
        "fetchTradingFee": false,
        "fetchTradingFees": false,
        "fetchTransactions": null,
        "fetchTransfers": true,
        "fetchWithdrawals": true,
        "setLeverage": true,
        "setMarginMode": true,
        "setPositionMode": true,
        "transfer": true,
        "withdraw": true
    },
    "timeframes": {
        "1m": "1",
        "3m": "3",
        "5m": "5",
        "15m": "15",
        "30m": "30",
        "1h": "60",
        "2h": "120",
        "4h": "240",
        "6h": "360",
        "12h": "720",
        "1d": "D",
        "1w": "W",
        "1M": "M",
        "1y": "Y"
    },
    "urls": {
        "test": {
            "spot": "https://api-testnet.{hostname}",
            "futures": "https://api-testnet.{hostname}",
            "v2": "https://api-testnet.{hostname}",
            "public": "https://api-testnet.{hostname}",
            "private": "https://api-testnet.{hostname}"
        },
        "logo": "https://user-images.githubusercontent.com/51840849/76547799-daff5b80-649e-11ea-87fb-3be9bac08954.jpg",
        "api": {
            "spot": "https://api.{hostname}",
            "futures": "https://api.{hostname}",
            "v2": "https://api.{hostname}",
            "public": "https://api.{hostname}",
            "private": "https://api.{hostname}"
        },
        "www": "https://www.bybit.com",
        "doc": [
            "https://bybit-exchange.github.io/docs/inverse/",
            "https://bybit-exchange.github.io/docs/linear/",
            "https://github.com/bybit-exchange"
        ],
        "fees": "https://help.bybit.com/hc/en-us/articles/360039261154",
        "referral": "https://partner.bybit.com/b/ccxt"
    },
    "api": {
        "public": {
            "get": [
				"public/symbols",
				"public/quote/depth",
				"public/quote/depth/merged",
				"public/quote/trades",
				"public/quote/kline",
				"public/quote/ticker/24hr",
				"public/quote/ticker/price",
				"public/quote/ticker/bookTicker",
				"public/server-time",
            ]
        },
        "private": {
            "get": [
				"private/order",
				"private/open-orders",
				"private/history-orders",
				"private/my-trades",
				"private/account",
            ],
            "post": [
				"private/order",
				"private/cancel-order",
				"private/cancel-orders",
				"private/cancel-orders-by-ids",
            ],
        }
    },
    "httpExceptions": {
        "403": "RateLimitExceeded"
    },
    "exceptions": {
        "exact": {
            "-10009": "BadRequest",
            "-1004": "BadRequest",
            "-1021": "BadRequest",
            "-1103": "BadRequest",
            "-1140": "InvalidOrder",
            "-1197": "InvalidOrder",
            "-2013": "InvalidOrder",
            "-2015": "AuthenticationError",
            "-6017": "BadRequest",
            "-6025": "BadRequest",
            "-6029": "BadRequest",
            "7001": "BadRequest",
            "10001": "BadRequest",
            "10002": "InvalidNonce",
            "10003": "AuthenticationError",
            "10004": "AuthenticationError",
            "10005": "PermissionDenied",
            "10006": "RateLimitExceeded",
            "10007": "AuthenticationError",
            "10010": "PermissionDenied",
            "10016": "ExchangeError",
            "10017": "BadRequest",
            "10018": "RateLimitExceeded",
            "12213": "OrderNotFound",
            "20001": "OrderNotFound",
            "20003": "InvalidOrder",
            "20004": "InvalidOrder",
            "20005": "InvalidOrder",
            "20006": "InvalidOrder",
            "20007": "InvalidOrder",
            "20008": "InvalidOrder",
            "20009": "InvalidOrder",
            "20010": "InvalidOrder",
            "20011": "InvalidOrder",
            "20012": "InvalidOrder",
            "20013": "InvalidOrder",
            "20014": "InvalidOrder",
            "20015": "InvalidOrder",
            "20016": "InvalidOrder",
            "20017": "InvalidOrder",
            "20018": "InvalidOrder",
            "20019": "InvalidOrder",
            "20020": "InvalidOrder",
            "20021": "InvalidOrder",
            "20022": "BadRequest",
            "20023": "BadRequest",
            "20031": "BadRequest",
            "20070": "BadRequest",
            "20071": "BadRequest",
            "20084": "BadRequest",
            "30001": "BadRequest",
            "30003": "InvalidOrder",
            "30004": "InvalidOrder",
            "30005": "InvalidOrder",
            "30007": "InvalidOrder",
            "30008": "InvalidOrder",
            "30009": "ExchangeError",
            "30010": "InsufficientFunds",
            "30011": "PermissionDenied",
            "30012": "PermissionDenied",
            "30013": "PermissionDenied",
            "30014": "InvalidOrder",
            "30015": "InvalidOrder",
            "30016": "ExchangeError",
            "30017": "InvalidOrder",
            "30018": "InvalidOrder",
            "30019": "InvalidOrder",
            "30020": "InvalidOrder",
            "30021": "InvalidOrder",
            "30022": "InvalidOrder",
            "30023": "InvalidOrder",
            "30024": "InvalidOrder",
            "30025": "InvalidOrder",
            "30026": "InvalidOrder",
            "30027": "InvalidOrder",
            "30028": "InvalidOrder",
            "30029": "InvalidOrder",
            "30030": "InvalidOrder",
            "30031": "InsufficientFunds",
            "30032": "InvalidOrder",
            "30033": "RateLimitExceeded",
            "30034": "OrderNotFound",
            "30035": "RateLimitExceeded",
            "30036": "ExchangeError",
            "30037": "InvalidOrder",
            "30041": "ExchangeError",
            "30042": "InsufficientFunds",
            "30043": "InvalidOrder",
            "30044": "InvalidOrder",
            "30045": "InvalidOrder",
            "30049": "InsufficientFunds",
            "30050": "ExchangeError",
            "30051": "ExchangeError",
            "30052": "ExchangeError",
            "30054": "ExchangeError",
            "30057": "ExchangeError",
            "30063": "ExchangeError",
            "30067": "InsufficientFunds",
            "30068": "ExchangeError",
            "30074": "InvalidOrder",
            "30075": "InvalidOrder",
            "30078": "ExchangeError",
            "33004": "AuthenticationError",
            "34026": "ExchangeError",
            "34036": "BadRequest",
            "35015": "BadRequest",
            "130006": "InvalidOrder",
            "130021": "InsufficientFunds",
            "130074": "InvalidOrder",
            "3100116": "BadRequest",
            "3100198": "BadRequest",
            "3200300": "InsufficientFunds"
        },
        "broad": {
            "unknown orderInfo": "OrderNotFound",
            "invalid api_key": "AuthenticationError",
            "oc_diff": "InsufficientFunds",
            "new_oc": "InsufficientFunds"
        }
    },
    "precisionMode": "TICK_SIZE",
    "options": {
        "createMarketBuyOrderRequiresPrice": true,
        "defaultType": "swap",
        "defaultSubType": "linear",
        "defaultSettle": "USDT",
        "code": "BTC",
        "recvWindow": 5000,
        "timeDifference": 0,
        "adjustForTimeDifference": false,
        "brokerId": "CCXT",
        "accountsByType": {
            "spot": "SPOT",
            "margin": "SPOT",
            "future": "CONTRACT",
            "swap": "CONTRACT",
            "option": "OPTION"
        },
        "accountsById": {
            "SPOT": "spot",
            "MARGIN": "spot",
            "CONTRACT": "contract",
            "OPTION": "option"
        }
    },
    "fees": {
        "trading": {
            "feeSide": "get",
            "tierBased": true,
            "percentage": true,
            "taker": 0.00075,
            "maker": 0.0001
        },
        "funding": {
            "tierBased": false,
            "percentage": false,
            "withdraw": {},
            "deposit": {}
        }
    }
}`)
}

func (self *Bybit) LoadMarkets() map[string]*Market {
	return nil
}

func (self *Bybit) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	return &Market{
		Id:     li[0] + li[1],
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *Bybit) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFunc("publicGetPublicQuoteDepth", self.Extend(request, params), nil, nil)
	result := response["result"].(map[string]interface{})
	timestamp := self.SafeInteger(result, "time", 0)
	return self.ParseOrderBook(result, timestamp, "bids", "asks", 0, 1), nil
}

func (self *Bybit) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	response := self.ApiFunc("privateGetPrivateAccount", params, nil, nil)
	balances := response["result"].(map[string]interface{})["balances"].([]interface{})
	result := map[string]interface{}{
		"info": response,
	}
	for _, balance := range balances {
		account := self.Account()
		account["free"] = self.SafeFloat(balance, "free", 0)
		account["used"] = self.SafeFloat(balance, "locked", 0)
		account["total"] = self.SafeFloat(balance, "total", 0)
		currencyId := self.SafeString(balance, "coinId", "")
		result[currencyId] = account
	}
	return self.ParseBalance(result), nil
}

func (self *Bybit) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
		"limit":  "500",
	}
	response := self.ApiFunc("privateGetPrivateOpenOrders", self.Extend(request, params), nil, nil)
	rs := self.SafeValue(response, "result", map[string]interface{}{})
	orders := self.SafeValue(rs, "list", []interface{}{})
	return self.ToOrders(self.ParseOrders(orders, market, since, limit)), nil
}

func (self *Bybit) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	symbol := self.SafeString(order, "symbol", "")
	if market != nil {
		symbol = market.(*Market).Symbol
	}
	timestamp := self.SafeInteger(order, "createTime", 0)
	id := self.SafeString(order, "orderId", "")
	price := self.SafeFloat(order, "orderPrice", 0)
	average := self.SafeFloat(order, "avgPrice", 0)
	amount := self.SafeFloat(order, "orderQty", 0)
	filled := self.SafeFloat(order, "execQty", 0)
	remaining := amount - filled
	raw_status := self.SafeString(order, "status", "")
	status := self.ParseOrderStatus(raw_status)
	side := self.SafeStringLower(order, "side", "")
	clientOrderId := self.SafeString(order, "orderLinkId", "")
	type_ := self.SafeStringLower(order, "orderType", "")
	return map[string]interface{}{
		"info":          order,
		"id":            id,
		"clientOrderId": clientOrderId,
		"timestamp":     timestamp,
		"datetime":      self.Iso8601(timestamp),
		"symbol":        symbol,
		"side":          side,
		"price":         price,
		"amount":        amount,
		"average":       average,
		"filled":        filled,
		"remaining":     remaining,
		"type":          type_,
		"status":        status,
	}
}

func (self *Bybit) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"REJECTED":         "rejected",
		"NEW":              "open",
		"PARTIALLY_FILLED": "open",
		"FILLED":           "closed",
		"CANCELED":         "canceled",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Bybit) CreateOrder(symbol string, type_ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol":     market.Id,
		"orderQty":   self.Float64ToString(amount),
		"side":       self.Capitalize(side),
		"orderType":  strings.ToUpper(type_),
		"orderPrice": self.Float64ToString(price),
	}
	response := self.ApiFunc("privatePostPrivateOrder", self.Extend(request, params), nil, nil)
	data := response["result"]
	order := map[string]interface{}{
		"id":            self.SafeString(data, "orderId", ""),
		"symbol":        symbol,
		"type":          type_,
		"side":          side,
		"price":         self.SafeFloat(data, "orderPrice", 0),
		"amount":        self.SafeFloat(data, "orderQty", 0),
		"timestamp":     self.SafeInteger(data, "createTime", 0),
		"clientOrderId": self.SafeString(data, "orderLinkId", ""),
		"info":          data,
	}
	return self.ToOrder(order), nil
}

func (self *Bybit) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	request := map[string]interface{}{}
	if id != "" {
		request["orderId"] = id
	}
	response := self.ApiFunc("privateGetPrivateOrder", self.Extend(request, params), nil, nil)
	order := self.ParseOrder(response["result"], nil)
	return self.ToOrder(order), nil
}

func (self *Bybit) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	request := map[string]interface{}{}
	if id != "" {
		request["orderId"] = id
	}
	response = self.ApiFunc("privatePostPrivateCancelOrder", self.Extend(request, params), nil, nil)
	return response, nil
}

func (self *Bybit) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	//url := self.ImplodeHostname(self.Member(self.Member(self.Urls, "api"), api).(string)) + "/spot/v3/" + path
	url := self.ImplodeHostname(self.DescribeJson.Get("urls.api").Get(api).String()) + "/spot/" + self.Version + "/" + path
	if api == "public" {
		if len(params) > 0 {
			url += "?" + self.Urlencode(params)
		}
	} else if api == "private" {
		recvWindow := self.DescribeJson.Get("options.recvWindow").String()
		timestamp := self.Nonce()
		payload := fmt.Sprintf("%d%s%s", timestamp, self.ApiKey, recvWindow)
		if method == "GET" {
			queryString := self.Urlencode(params)
			url += "?" + queryString
			payload += queryString
		} else {
			// POST
			body = self.Json(params)
			payload += body.(string)
		}
		signature := self.Hmac(payload, self.Secret, "sha256", "hex")
		headers = self.Extend(map[string]interface{}{
			"X-BAPI-SIGN-TYPE":   "2",
			"X-BAPI-SIGN":        signature,
			"X-BAPI-API-KEY":     self.ApiKey,
			"X-BAPI-TIMESTAMP":   fmt.Sprint(timestamp),
			"X-BAPI-RECV-WINDOW": recvWindow,
		}, headers)
	}
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Bybit) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if !self.ToBool(response) {
		return
	}
	errorCode := self.SafeString(response, "retCode", "0")
	if errorCode != "0" {
		feedback := self.Id + " " + body
		self.ThrowBroadlyMatchedException(self.Member(self.Exceptions, "broad"), body, feedback)
		self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), errorCode, feedback)
		self.RaiseException("ExchangeError", feedback)
	}
}
