package huobipro

import (
	"fmt"
	. "github.com/epheien/ccxt/go/base"
	"math"
	"reflect"
	"strings"
)

type Huobipro struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Huobipro, err error) {
	ex = new(Huobipro)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Huobipro) Describe() []byte {
	return []byte(`{
    "id": "huobipro",
    "name": "Huobi Pro",
    "countries": [
        "CN"
    ],
    "rateLimit": 2000,
    "version": "v1",
    "accounts": null,
    "accountsById": null,
    "hostname": "api.huobi.pro",
    "pro": true,
    "has": {
        "CORS": false,
        "fetchTickers": true,
        "fetchDepositAddress": true,
        "fetchOHLCV": true,
        "fetchOrder": true,
        "fetchOrders": true,
        "fetchOpenOrders": true,
        "fetchClosedOrders": true,
        "fetchTradingLimits": true,
        "fetchMyTrades": true,
        "withdraw": true,
        "fetchCurrencies": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true
    },
    "timeframes": {
        "1m": "1min",
        "5m": "5min",
        "15m": "15min",
        "30m": "30min",
        "1h": "60min",
        "4h": "4hour",
        "1d": "1day",
        "1w": "1week",
        "1M": "1mon",
        "1y": "1year"
    },
    "urls": {
        "test": {
            "market": "https://api.testnet.huobi.pro",
            "public": "https://api.testnet.huobi.pro",
            "private": "https://api.testnet.huobi.pro"
        },
        "logo": "https://user-images.githubusercontent.com/1294454/76137448-22748a80-604e-11ea-8069-6e389271911d.jpg",
        "api": {
            "market": "https://{hostname}",
            "public": "https://{hostname}",
            "private": "https://{hostname}",
            "v2Public": "https://{hostname}",
            "v2Private": "https://{hostname}"
        },
        "www": "https://www.huobi.pro",
        "referral": "https://www.huobi.co/en-us/topic/invited/?invite_code=rwrd3",
        "doc": "https://huobiapi.github.io/docs/spot/v1/cn/",
        "fees": "https://www.huobi.pro/about/fee/"
    },
    "api": {
        "v2Public": {
            "get": [
                "reference/currencies"
            ]
        },
        "v2Private": {
            "get": [
                "account/ledger",
                "account/withdraw/quota",
                "account/deposit/address",
                "reference/transact-fee-rate"
            ],
            "post": [
                "sub-user/management"
            ]
        },
        "market": {
            "get": [
                "history/kline",
                "detail/merged",
                "depth",
                "trade",
                "history/trade",
                "detail",
                "tickers"
            ]
        },
        "public": {
            "get": [
                "common/symbols",
                "common/currencys",
                "common/timestamp",
                "common/exchange",
                "settings/currencys"
            ]
        },
        "private": {
            "get": [
                "account/accounts",
                "account/accounts/{id}/balance",
                "account/accounts/{sub-uid}",
                "account/history",
                "cross-margin/loan-info",
                "fee/fee-rate/get",
                "order/openOrders",
                "order/orders",
                "order/orders/{id}",
                "order/orders/{id}/matchresults",
                "order/orders/getClientOrder",
                "order/history",
                "order/matchresults",
                "dw/withdraw-virtual/addresses",
                "query/deposit-withdraw",
                "margin/loan-orders",
                "margin/accounts/balance",
                "points/actions",
                "points/orders",
                "subuser/aggregate-balance",
                "stable-coin/exchange_rate",
                "stable-coin/quote"
            ],
            "post": [
                "futures/transfer",
                "order/batch-orders",
                "order/orders/place",
                "order/orders/submitCancelClientOrder",
                "order/orders/batchCancelOpenOrders",
                "order/orders",
                "order/orders/{id}/place",
                "order/orders/{id}/submitcancel",
                "order/orders/batchcancel",
                "dw/balance/transfer",
                "dw/withdraw/api/create",
                "dw/withdraw-virtual/create",
                "dw/withdraw-virtual/{id}/place",
                "dw/withdraw-virtual/{id}/cancel",
                "dw/transfer-in/margin",
                "dw/transfer-out/margin",
                "margin/orders",
                "margin/orders/{id}/repay",
                "stable-coin/exchange",
                "subuser/transfer"
            ]
        }
    },
    "fees": {
        "trading": {
            "tierBased": false,
            "percentage": true,
            "maker": 0.002,
            "taker": 0.002
        }
    },
    "exceptions": {
        "exact": {
            "bad-request": "BadRequest",
            "api-not-support-temp-addr": "PermissionDenied",
            "timeout": "RequestTimeout",
            "gateway-internal-error": "ExchangeNotAvailable",
            "account-frozen-balance-insufficient-error": "InsufficientFunds",
            "invalid-amount": "InvalidOrder",
            "order-limitorder-amount-min-error": "InvalidOrder",
            "order-limitorder-amount-max-error": "InvalidOrder",
            "order-marketorder-amount-min-error": "InvalidOrder",
            "order-limitorder-price-min-error": "InvalidOrder",
            "order-limitorder-price-max-error": "InvalidOrder",
            "order-orderstate-error": "OrderNotFound",
            "order-queryorder-invalid": "OrderNotFound",
            "order-update-error": "ExchangeNotAvailable",
            "api-signature-check-failed": "AuthenticationError",
            "api-signature-not-valid": "AuthenticationError",
            "base-record-invalid": "OrderNotFound",
            "invalid symbol": "BadSymbol",
            "invalid-parameter": "BadRequest",
            "base-symbol-trade-disabled": "BadSymbol"
        }
    },
    "options": {
        "fetchOrdersByStatesMethod": "privateGetOrderOrders",
        "fetchOpenOrdersMethod": "fetch_open_orders_v1",
        "createMarketBuyOrderRequiresPrice": true,
        "fetchMarketsMethod": "publicGetCommonSymbols",
        "fetchBalanceMethod": "privateGetAccountAccountsIdBalance",
        "createOrderMethod": "privatePostOrderOrdersPlace",
        "language": "en-US"
    },
    "commonCurrencies": {
        "GET": "Themis",
        "HOT": "Hydro Protocol"
    }
}`)
}

func (self *Huobipro) FetchMarkets(params map[string]interface{}) []interface{} {
	method := self.Member(self.Options, "fetchMarketsMethod")
	response := self.ApiFunc(method.(string), params, nil, nil)
	markets := self.SafeValue(response, "data", nil)
	numMarkets := self.Length(markets)
	if self.ToBool(numMarkets < 1) {
		self.RaiseException("ExchangeError", self.Id+" publicGetCommonSymbols returned empty response: "+self.Json(markets))
	}
	result := []interface{}{}
	for i := 0; i < self.Length(markets); i++ {
		market := self.Member(markets, i)
		baseId := self.SafeString(market, "base-currency", "")
		quoteId := self.SafeString(market, "quote-currency", "")
		id := baseId + quoteId
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		precision := map[string]interface{}{
			"amount": self.Member(market, "amount-precision"),
			"price":  self.Member(market, "price-precision"),
		}
		maker := self.IfThenElse(self.ToBool(base == "OMG"), 0., 0.2/100)
		taker := self.IfThenElse(self.ToBool(base == "OMG"), 0., 0.2/100)
		minAmount := self.SafeFloat(market, "min-order-amt", math.Pow10(int(-precision["amount"].(float64))))
		maxAmount := self.SafeFloat(market, "max-order-amt", 0)
		minCost := self.SafeFloat(market, "min-order-value", 0)
		state := self.SafeString(market, "state", "")
		active := state == "online"
		result = append(result, map[string]interface{}{
			"id":        id,
			"symbol":    symbol,
			"base":      base,
			"quote":     quote,
			"baseId":    baseId,
			"quoteId":   quoteId,
			"active":    active,
			"precision": precision,
			"taker":     taker,
			"maker":     maker,
			"limits": map[string]interface{}{
				"amount": map[string]interface{}{
					"min": minAmount,
					"max": maxAmount,
				},
				"price": map[string]interface{}{
					"min": math.Pow10(-int(precision["price"].(float64))),
					"max": nil,
				},
				"cost": map[string]interface{}{
					"min": minCost,
					"max": nil,
				},
			},
			"info": market,
		})
	}
	return result
}

func (self *Huobipro) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": self.Member(market, "id"),
		"type":   "step0",
	}
	response := self.ApiFunc("marketGetDepth", self.Extend(request, params), nil, nil)
	if self.ToBool(self.InMap("tick", response)) {
		if self.ToBool(!self.ToBool(self.Member(response, "tick"))) {
			self.RaiseException("ExchangeError", self.Id+" fetchOrderBook() returned empty response: "+self.Json(response))
		}
		tick := self.SafeValue(response, "tick", nil)
		timestamp := self.SafeInteger(tick, "ts", self.SafeInteger(response, "ts", 0))
		result := self.ParseOrderBook(tick, timestamp, "bids", "asks", 0, 1)
		self.SetValue(result, "nonce", self.SafeInteger(tick, "version", 0))
		return result, nil
	}
	self.RaiseException("ExchangeError", self.Id+" fetchOrderBook() returned unrecognized response: "+self.Json(response))
	return
}

func (self *Huobipro) FetchCurrencies(params map[string]interface{}) map[string]interface{} {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	request := map[string]interface{}{
		"language": self.Member(self.Options, "language"),
	}
	response := self.ApiFunc("publicGetSettingsCurrencys", self.Extend(request, params), nil, nil)
	currencies := self.SafeValue(response, "data", nil)
	result := map[string]interface{}{}
	for i := 0; i < self.Length(currencies); i++ {
		currency := self.Member(currencies, i)
		id := self.SafeValue(currency, "name", nil)
		precision := self.SafeInteger(currency, "withdraw-precision", 0)
		code := self.SafeCurrencyCode(id)
		active := self.Member(currency, "visible").(bool) && self.Member(currency, "deposit-enabled").(bool) && self.Member(currency, "withdraw-enabled").(bool)
		name := self.SafeString(currency, "display-name", "")
		self.SetValue(result, code, map[string]interface{}{
			"id":        id,
			"code":      code,
			"type":      "crypto",
			"name":      name,
			"active":    active,
			"fee":       nil,
			"precision": precision,
			"limits": map[string]interface{}{
				"amount": map[string]interface{}{
					"min": math.Pow10(-int(precision)),
					"max": math.Pow10(int(precision)),
				},
				"price": map[string]interface{}{
					"min": math.Pow10(-int(precision)),
					"max": math.Pow10(int(precision)),
				},
				"cost": map[string]interface{}{
					"min": nil,
					"max": nil,
				},
				"deposit": map[string]interface{}{
					"min": self.SafeFloat(currency, "deposit-min-amount", 0),
					"max": math.Pow10(int(precision)),
				},
				"withdraw": map[string]interface{}{
					"min": self.SafeFloat(currency, "withdraw-min-amount", 0),
					"max": math.Pow10(int(precision)),
				},
			},
			"info": currency,
		})
	}
	return result
}

func (self *Huobipro) FetchAccounts(params map[string]interface{}) []interface{} {
	self.LoadMarkets()
	response := self.ApiFunc("privateGetAccountAccounts", params, nil, nil)
	return self.Member(response, "data").([]interface{})
}

func (self *Huobipro) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	method := self.Member(self.Options, "fetchBalanceMethod").(string)
	request := map[string]interface{}{
		"id": self.Member(self.Member(self.Accounts, 0), "id"),
	}
	response := self.ApiFunc(method, request, nil, nil)
	balances := self.SafeValue(self.Member(response, "data"), "list", []interface{}{})
	result := map[string]interface{}{
		"info": response,
	}
	for i := 0; i < self.Length(balances); i++ {
		balance := self.Member(balances, i)
		currencyId := self.SafeString(balance, "currency", "")
		code := self.SafeCurrencyCode(currencyId)
		var account interface{}
		if self.ToBool(self.InMap(code, result)) {
			account = self.Member(result, code)
		} else {
			account = self.Account()
		}
		if self.ToBool(self.Member(balance, "type") == "trade") {
			self.SetValue(account, "free", self.SafeFloat(balance, "balance", 0))
		}
		if self.ToBool(self.Member(balance, "type") == "frozen") {
			self.SetValue(account, "used", self.SafeFloat(balance, "balance", 0))
		}
		self.SetValue(result, code, account)
	}
	return self.ParseBalance(result), nil
}

func (self *Huobipro) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	request := map[string]interface{}{
		"id": id,
	}
	response := self.ApiFunc("privateGetOrderOrdersId", self.Extend(request, params), nil, nil)
	order := self.SafeValue(response, "data", nil)
	return self.ToOrder(self.ParseOrder(order, nil)), nil
}

func (self *Huobipro) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	method := self.SafeString(self.Options, "fetchOpenOrdersMethod", "fetch_open_orders_v1")
	if method == "fetch_open_orders_v1" {
		return self.ToOrders(self.fetch_open_orders_v1(symbol, since, limit, params)), nil
	} else {
		self.RaiseInternalException("unsported method: " + method)
	}
	return
}

func (self *Huobipro) FetchOrdersByStates(states string, symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}) {
	self.LoadMarkets()
	request := map[string]interface{}{
		"states": states,
	}
	var market interface{}
	if symbol != "" {
		market = self.Market(symbol)
		self.SetValue(request, "symbol", self.Member(market, "id"))
	}
	method := self.SafeString(self.Options, "fetchOrdersByStatesMethod", "privateGetOrderOrders")
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	return self.ParseOrders(self.Member(response, "data"), market, since, limit)
}

func (self *Huobipro) fetch_open_orders_v1(symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}) {
	if symbol == "" {
		self.RaiseInternalException(self.Id + " fetchOpenOrdersV1 requires a symbol argument")
	}
	return self.FetchOrdersByStates("pre-submitted,submitted,partial-filled", symbol, since, limit, params)
}

func (self *Huobipro) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"partial-filled":   "open",
		"partial-canceled": "canceled",
		"filled":           "closed",
		"canceled":         "canceled",
		"submitted":        "open",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Huobipro) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	id := ""
	if order.(map[string]interface{})["id"] != nil {
		id = fmt.Sprintf("%v", int64(order.(map[string]interface{})["id"].(float64)))
	}
	//id := self.SafeString(order, "id", "")
	var side interface{}
	var typ interface{}
	var status interface{}
	if self.ToBool(self.InMap("type", order)) {
		orderType := strings.Split(self.Member(order, "type").(string), "-")
		side = self.Member(orderType, 0)
		typ = self.Member(orderType, 1)
		status = self.ParseOrderStatus(self.SafeString(order, "state", ""))
	}
	var symbol interface{}
	if self.ToBool(self.TestNil(market)) {
		if self.ToBool(self.InMap("symbol", order)) {
			if self.ToBool(self.InMap(self.Member(order, "symbol"), self.MarketsById)) {
				marketId := self.Member(order, "symbol")
				market = self.Member(self.MarketsById, marketId)
			}
		}
	}
	if self.ToBool(!self.TestNil(market)) {
		symbol = market.(*Market).Symbol
	}
	timestamp := self.SafeInteger(order, "created-at", 0)
	amount := self.SafeFloat(order, "amount", 0)
	filled := self.SafeFloat2(order, "filled-amount", "field-amount", 0.0)
	if typ == "market" && side == "buy" {
		if status == "closed" {
			amount = filled
		} else {
			amount = 0.
		}
	}
	var price interface{}
	price = self.SafeFloat(order, "price", 0)
	if price == 0 {
		price = nil
	}
	cost := self.SafeFloat2(order, "filled-cash-amount", "field-cash-amount", 0.0)
	var remaining interface{}
	var average interface{}
	if self.ToBool(!self.TestNil(filled)) {
		if self.ToBool(!self.TestNil(amount)) {
			remaining = amount - filled
		}
		if self.ToBool(!self.TestNil(cost) && filled > 0) {
			average = cost / filled
		}
	}
	feeCost := self.SafeFloat2(order, "filled-fees", "field-fees", 0.0)
	var fee interface{}
	if self.ToBool(!self.TestNil(feeCost)) {
		var feeCurrency interface{}
		if self.ToBool(!self.TestNil(market)) {
			feeCurrency = self.IfThenElse(self.ToBool(side == "sell"), self.Member(market, "quote"), self.Member(market, "base"))
		}
		fee = map[string]interface{}{
			"cost":     feeCost,
			"currency": feeCurrency,
		}
	}
	return map[string]interface{}{
		"info":               order,
		"id":                 id,
		"clientOrderId":      nil,
		"timestamp":          timestamp,
		"datetime":           self.Iso8601(timestamp),
		"lastTradeTimestamp": nil,
		"symbol":             symbol,
		"type":               typ,
		"side":               side,
		"price":              price,
		"average":            average,
		"cost":               cost,
		"amount":             amount,
		"filled":             filled,
		"remaining":          remaining,
		"status":             status,
		"fee":                fee,
		"trades":             nil,
	}
}

func (self *Huobipro) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"account-id": self.Member(self.Member(self.Accounts, 0), "id"),
		"symbol":     market.Id,
		"type":       side + "-" + typ,
	}
	if self.ToBool(typ == "market" && side == "buy") {
		if self.ToBool(self.Member(self.Options, "createMarketBuyOrderRequiresPrice")) {
			if self.ToBool(self.TestNil(price)) {
				self.RaiseException("InvalidOrder", self.Id+" market buy order requires price argument to calculate cost (total amount of quote currency to spend for buying, amount * price). To switch off this warning exception and specify cost in the amount argument, set .options[createMarketBuyOrderRequiresPrice] = false. Make sure you know what youre doing.")
			} else {
				self.SetValue(request, "amount", self.CostToPrecision(symbol, amount*price))
			}
		} else {
			self.SetValue(request, "amount", self.CostToPrecision(symbol, amount))
		}
	} else {
		self.SetValue(request, "amount", self.AmountToPrecision(symbol, amount))
	}
	if self.ToBool(typ == "limit" || typ == "ioc" || typ == "limit-maker") {
		self.SetValue(request, "price", self.PriceToPrecision(symbol, price))
	}
	method := self.Member(self.Options, "createOrderMethod")
	response := self.ApiFunc(method.(string), self.Extend(params, request), nil, nil)
	timestamp := self.Milliseconds()
	id := self.SafeString(response, "data", "")
	return self.ToOrder(map[string]interface{}{
		"info":               response,
		"id":                 id,
		"timestamp":          timestamp,
		"datetime":           self.Iso8601(timestamp),
		"lastTradeTimestamp": nil,
		"status":             nil,
		"symbol":             symbol,
		"type":               typ,
		"side":               side,
		"price":              price,
		"amount":             amount,
		"filled":             nil,
		"remaining":          nil,
		"cost":               nil,
		"trades":             nil,
		"fee":                nil,
		"clientOrderId":      nil,
		"average":            nil,
	}), nil
}

func (self *Huobipro) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	response = self.ApiFunc("privatePostOrderOrdersIdSubmitcancel", map[string]interface{}{
		"id": id,
	}, nil, nil)
	return self.Extend(self.ParseOrder(response, nil), map[string]interface{}{
		"id":     id,
		"status": "canceled",
	}), nil
}

func (self *Huobipro) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	url := "/"
	if self.ToBool(api == "market") {
		url += api
	} else if self.ToBool(api == "public" || api == "private") {
		url += self.Version
	} else if self.ToBool(api == "v2Public" || api == "v2Private") {
		url += "v2"
	}
	url += "/" + self.ImplodeParams(path, params)
	query := self.Omit(params, self.ExtractParams(path))
	if self.ToBool(api == "private" || api == "v2Private") {
		timestamp := self.Ymdhms(self.Milliseconds(), "T")
		request := map[string]interface{}{
			"SignatureMethod":  "HmacSHA256",
			"SignatureVersion": "2",
			"AccessKeyId":      self.ApiKey,
			"Timestamp":        timestamp,
		}
		if self.ToBool(method != "POST") {
			request = self.Extend(request, query).(map[string]interface{})
		}
		// TODO, Urlencode iterate by sorted key
		// request = self.Keysort(request)
		auth := self.Urlencode(request)
		payload := strings.Join([]string{method, self.Hostname, url, auth}, "\n")
		signature := self.Hmac(self.Encode(payload), self.Encode(self.Secret), "sha256", "base64")
		auth += "&" + self.Urlencode(map[string]interface{}{
			"Signature": signature,
		})
		url += "?" + auth
		if self.ToBool(method == "POST") {
			body = self.Json(query)
			headers = map[string]interface{}{
				"Content-Type": "application/json",
			}
		} else {
			headers = map[string]interface{}{
				"Content-Type": "application/x-www-form-urlencoded",
			}
		}
	} else {
		if self.ToBool(self.Length(reflect.ValueOf(params).MapKeys())) {
			url += "?" + self.Urlencode(params)
		}
	}
	url = self.ImplodeParams(self.Member(self.Member(self.Urls, "api"), api).(string), map[string]interface{}{
		"hostname": self.Hostname,
	}) + url
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Huobipro) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if self.ToBool(self.TestNil(response)) {
		return
	}
	if self.ToBool(self.InMap("status", response)) {
		status := self.SafeString(response, "status", "")
		if self.ToBool(status == "error") {
			code := self.SafeString(response, "err-code", "")
			feedback := self.Id + " " + body
			self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), code, feedback)
			message := self.SafeString(response, "err-msg", "")
			self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), message, feedback)
			self.RaiseException("ExchangeError", feedback)
		}
	}
}
