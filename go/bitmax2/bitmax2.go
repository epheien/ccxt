package bitmax2

import (
	"fmt"
	. "github.com/epheien/ccxt/go/base"
	"reflect"
	"strings"
)

type Bitmax2 struct {
	Exchange
	accountGroup string
}

func New(config *ExchangeConfig) (ex *Bitmax2, err error) {
	ex = new(Bitmax2)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Bitmax2) Describe() []byte {
	return []byte(`{
    "id": "ascendex",
    "name": "AscendEX",
    "countries": [
        "SG"
    ],
    "rateLimit": 500,
    "certified": true,
    "has": {
        "CORS": false,
        "fetchMarkets": true,
        "fetchCurrencies": true,
        "fetchOrderBook": true,
        "fetchTicker": true,
        "fetchTickers": true,
        "fetchOHLCV": true,
        "fetchTrades": true,
        "fetchAccounts": true,
        "fetchBalance": true,
        "createOrder": true,
        "cancelOrder": true,
        "cancelAllOrders": true,
        "fetchDepositAddress": true,
        "fetchTransactions": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchOrder": true,
        "fetchOrders": true,
        "fetchOpenOrders": true,
        "fetchClosedOrders": true
    },
    "timeframes": {
        "1m": "1",
        "5m": "5",
        "15m": "15",
        "30m": "30",
        "1h": "60",
        "2h": "120",
        "4h": "240",
        "6h": "360",
        "12h": "720",
        "1d": "1d",
        "1w": "1w",
        "1M": "1m"
    },
    "version": "v1",
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/112027508-47984600-8b48-11eb-9e17-d26459cc36c6.jpg",
        "api": "https://ascendex.com",
        "test": "https://bitmax-test.io",
        "www": "https://ascendex.com",
        "doc": [
            "https://bitmax-exchange.github.io/bitmax-pro-api/#bitmax-pro-api-documentation"
        ],
        "fees": "https://ascendex.com/en/feerate/transactionfee-traderate",
        "referral": {
            "url": "https://ascendex.com/en-us/register?inviteCode=EL6BXBQM",
            "discount": 0.25
        }
    },
    "api": {
        "public": {
            "get": [
                "assets",
                "products",
                "ticker",
                "barhist/info",
                "barhist",
                "depth",
                "trades",
                "cash/assets",
                "cash/products",
                "margin/assets",
                "margin/products",
                "futures/collateral",
                "futures/contracts",
                "futures/ref-px",
                "futures/market-data",
                "futures/funding-rates"
            ]
        },
        "accountCategory": {
            "get": [
                "balance",
                "order/open",
                "order/status",
                "order/hist/current",
                "risk"
            ],
            "post": [
                "order",
                "order/batch"
            ],
            "delete": [
                "order",
                "order/all",
                "order/batch"
            ]
        },
        "accountGroup": {
            "get": [
                "cash/balance",
                "margin/balance",
                "margin/risk",
                "transfer",
                "futures/collateral-balance",
                "futures/position",
                "futures/risk",
                "futures/funding-payments",
                "order/hist"
            ],
            "post": [
                "futures/transfer/deposit",
                "futures/transfer/withdraw"
            ]
        },
        "private": {
            "get": [
                "info",
                "wallet/transactions",
                "wallet/deposit/address"
            ]
        }
    },
    "fees": {
        "trading": {
            "feeSide": "get",
            "tierBased": true,
            "percentage": true
        }
    },
    "precisionMode": "TICK_SIZE",
    "options": {
        "account-category": "cash",
        "account-group": null,
        "fetchClosedOrders": {
            "method": "accountGroupGetOrderHist"
        }
    },
    "exceptions": {
        "exact": {
            "1900": "BadRequest",
            "2100": "AuthenticationError",
            "5002": "BadSymbol",
            "6001": "BadSymbol",
            "6010": "InsufficientFunds",
            "60060": "InvalidOrder",
            "600503": "InvalidOrder",
            "100001": "BadRequest",
            "100002": "BadRequest",
            "100003": "BadRequest",
            "100004": "BadRequest",
            "100005": "BadRequest",
            "100006": "BadRequest",
            "100007": "BadRequest",
            "100008": "BadSymbol",
            "100009": "AuthenticationError",
            "100010": "BadRequest",
            "100011": "BadRequest",
            "100012": "BadRequest",
            "100013": "BadRequest",
            "100101": "ExchangeError",
            "150001": "BadRequest",
            "200001": "AuthenticationError",
            "200002": "ExchangeError",
            "200003": "ExchangeError",
            "200004": "ExchangeError",
            "200005": "ExchangeError",
            "200006": "ExchangeError",
            "200007": "ExchangeError",
            "200008": "ExchangeError",
            "200009": "ExchangeError",
            "200010": "AuthenticationError",
            "200011": "ExchangeError",
            "200012": "ExchangeError",
            "200013": "ExchangeError",
            "200014": "PermissionDenied",
            "200015": "PermissionDenied",
            "300001": "InvalidOrder",
            "300002": "InvalidOrder",
            "300003": "InvalidOrder",
            "300004": "InvalidOrder",
            "300005": "InvalidOrder",
            "300006": "InvalidOrder",
            "300007": "InvalidOrder",
            "300008": "InvalidOrder",
            "300009": "InvalidOrder",
            "300011": "InsufficientFunds",
            "300012": "BadSymbol",
            "300013": "InvalidOrder",
            "300020": "InvalidOrder",
            "300021": "InvalidOrder",
            "300031": "InvalidOrder",
            "310001": "InsufficientFunds",
            "310002": "InvalidOrder",
            "310003": "InvalidOrder",
            "310004": "BadSymbol",
            "310005": "InvalidOrder",
            "510001": "ExchangeError",
            "900001": "ExchangeError"
        },
        "broad": {}
    },
    "commonCurrencies": {
        "BOND": "BONDED",
        "BTCBEAR": "BEAR",
        "BTCBULL": "BULL",
        "BYN": "Beyond Finance"
    }
}`)
}

func (self *Bitmax2) FetchAccounts(params map[string]interface{}) []interface{} {
	accountGroup := self.accountGroup
	var response interface{}
	if self.ToBool(self.TestNil(accountGroup)) {
		response = self.ApiFunc("privateGetInfo", params, nil, nil)
		data := self.SafeValue(response, "data", map[string]interface{}{})
		accountGroup = self.SafeString(data, "accountGroup", "")
		self.accountGroup = accountGroup
	}
	return []interface{}{map[string]interface{}{
		"id":       accountGroup,
		"type":     nil,
		"currency": nil,
		"info":     response,
	}}
}

func (self *Bitmax2) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	defaultAccountCategory := self.SafeString(self.Options, "account-category", "cash")
	options := self.SafeValue(self.Options, "fetchBalance", map[string]interface{}{})
	accountCategory := self.SafeString(options, "account-category", defaultAccountCategory)
	accountCategory = self.SafeString(params, "account-category", accountCategory)
	params = self.Omit(params, "account-category")
	account := self.SafeValue(self.Accounts, 0, map[string]interface{}{})
	accountGroup := self.SafeString(account, "id", "")
	request := map[string]interface{}{
		"account-group": accountGroup,
	}
	method := "accountCategoryGetBalance"
	if self.ToBool(accountCategory == "futures") {
		method = "accountGroupGetFuturesCollateralBalance"
	} else {
		self.SetValue(request, "account-category", accountCategory)
	}
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	result := map[string]interface{}{
		"info":      response,
		"timestamp": nil,
		"datetime":  nil,
	}
	balances := self.SafeValue(response, "data", []interface{}{})
	for i := 0; i < self.Length(balances); i++ {
		balance := self.Member(balances, i)
		code := self.SafeCurrencyCode(self.SafeString(balance, "asset", ""))
		account := self.Account()
		free := self.SafeFloat(balance, "availableBalance", 0)
		total := self.SafeFloat(balance, "totalBalance", 0)
		account["free"] = free
		account["total"] = total
		account["used"] = total - free
		result[code] = account
	}
	return self.ParseBalance(result), nil
}

func (self *Bitmax2) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
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
	response := self.ApiFunc("publicGetDepth", self.Extend(request, params), nil, nil)
	data := self.SafeValue(response, "data", map[string]interface{}{})
	orderbook := self.SafeValue(data, "data", map[string]interface{}{})
	timestamp := self.SafeInteger(orderbook, "ts", 0)
	result := self.ParseOrderBook(orderbook, timestamp, "bids", "asks", 0, 1)
	self.SetValue(result, "nonce", self.SafeInteger(orderbook, "seqnum", 0))
	return result, nil
}

func (self *Bitmax2) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"PendingNew":      "open",
		"New":             "open",
		"PartiallyFilled": "open",
		"Filled":          "closed",
		"Canceled":        "canceled",
		"Rejected":        "rejected",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Bitmax2) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	status := self.ParseOrderStatus(self.SafeString(order, "status", ""))
	marketId := self.SafeString(order, "symbol", "")
	symbol := self.SafeSymbol(marketId, market, "/")
	timestamp := self.SafeInteger2(order, "timestamp", "sendingTime", 0)
	lastTradeTimestamp := self.SafeInteger(order, "lastExecTime", 0)
	price := self.SafeNumber(order, "price")
	amount := self.SafeNumber(order, "orderQty")
	average := self.SafeNumber(order, "avgPx")
	filled := self.SafeNumber2(order, "cumFilledQty", "cumQty")
	id := self.SafeString(order, "orderId", "")
	clientOrderId := self.SafeString(order, "id", "")
	if self.ToBool(!self.TestNil(clientOrderId)) {
		if self.ToBool(self.Length(clientOrderId) < 1) {
			clientOrderId = ""
		}
	}
	typ := self.SafeStringLower(order, "orderType", "")
	side := self.SafeStringLower(order, "side", "")
	feeCost := self.SafeNumber(order, "cumFee")
	var fee interface{}
	if self.ToBool(!self.TestNil(feeCost)) {
		feeCurrencyId := self.SafeString(order, "feeAsset", "")
		feeCurrencyCode := self.SafeCurrencyCode(feeCurrencyId)
		fee = map[string]interface{}{
			"cost":     feeCost,
			"currency": feeCurrencyCode,
		}
	}
	stopPrice := self.SafeNumber(order, "stopPrice")
	return map[string]interface{}{
		"info":               order,
		"id":                 id,
		"clientOrderId":      nil,
		"timestamp":          timestamp,
		"datetime":           self.Iso8601(timestamp),
		"lastTradeTimestamp": lastTradeTimestamp,
		"symbol":             symbol,
		"type":               typ,
		"timeInForce":        nil,
		"postOnly":           nil,
		"side":               side,
		"price":              price,
		"stopPrice":          stopPrice,
		"amount":             amount,
		"cost":               nil,
		"average":            average,
		"filled":             filled,
		"remaining":          nil,
		"status":             status,
		"fee":                fee,
		"trades":             nil,
	}
}

func (self *Bitmax2) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	market := self.Market(symbol)
	defaultAccountCategory := self.SafeString(self.Options, "account-category", "cash")
	options := self.SafeValue(self.Options, "createOrder", map[string]interface{}{})
	accountCategory := self.SafeString(options, "account-category", defaultAccountCategory)
	accountCategory = self.SafeString(params, "account-category", accountCategory)
	params = self.Omit(params, "account-category")
	account := self.SafeValue(self.Accounts, 0, map[string]interface{}{})
	accountGroup := self.SafeValue(account, "id", nil)
	clientOrderId := self.SafeString2(params, "clientOrderId", "id", "")
	request := map[string]interface{}{
		"account-group":    accountGroup,
		"account-category": accountCategory,
		"symbol":           market.Id,
		"time":             self.Milliseconds(),
		"orderQty":         self.AmountToPrecision(symbol, amount),
		"orderType":        typ,
		"side":             side,
	}
	if self.ToBool(!self.TestNil(clientOrderId)) {
		self.SetValue(request, "id", clientOrderId)
		params = self.Omit(params, []interface{}{"clientOrderId", "id"})
	}
	if self.ToBool(typ == "limit" || typ == "stop_limit") {
		self.SetValue(request, "orderPrice", self.PriceToPrecision(symbol, price))
	}
	if self.ToBool(typ == "stop_limit" || typ == "stop_market") {
		stopPrice := self.SafeNumber(params, "stopPrice")
		if self.ToBool(self.TestNil(stopPrice)) {
			self.RaiseException("InvalidOrder", self.Id+" createOrder() requires a stopPrice parameter for "+typ+" orders")
		} else {
			self.SetValue(request, "stopPrice", self.PriceToPrecision(symbol, stopPrice))
			params = self.Omit(params, "stopPrice")
		}
	}
	response := self.ApiFunc("accountCategoryPostOrder", self.Extend(request, params), nil, nil)
	data := self.SafeValue(response, "data", map[string]interface{}{})
	info := self.SafeValue(data, "info", map[string]interface{}{})
	return self.ToOrder(self.ParseOrder(info, market)), nil
}

func (self *Bitmax2) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	defaultAccountCategory := self.SafeString(self.Options, "account-category", "cash")
	options := self.SafeValue(self.Options, "fetchOrder", map[string]interface{}{})
	accountCategory := self.SafeString(options, "account-category", defaultAccountCategory)
	accountCategory = self.SafeString(params, "account-category", accountCategory)
	params = self.Omit(params, "account-category")
	account := self.SafeValue(self.Accounts, 0, map[string]interface{}{})
	accountGroup := self.SafeValue(account, "id", nil)
	request := map[string]interface{}{
		"account-group":    accountGroup,
		"account-category": accountCategory,
		"orderId":          id,
	}
	response := self.ApiFunc("accountCategoryGetOrderStatus", self.Extend(request, params), nil, nil)
	data := self.SafeValue(response, "data", map[string]interface{}{})
	return self.ToOrder(self.ParseOrder(data, nil)), nil
}

func (self *Bitmax2) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	self.LoadAccounts()
	var market interface{}
	if self.ToBool(!self.TestNil(symbol)) {
		market = self.Market(symbol)
	}
	defaultAccountCategory := self.SafeString(self.Options, "account-category", "cash")
	options := self.SafeValue(self.Options, "fetchOpenOrders", map[string]interface{}{})
	accountCategory := self.SafeString(options, "account-category", defaultAccountCategory)
	accountCategory = self.SafeString(params, "account-category", accountCategory)
	params = self.Omit(params, "account-category")
	account := self.SafeValue(self.Accounts, 0, map[string]interface{}{})
	accountGroup := self.SafeValue(account, "id", nil)
	request := map[string]interface{}{
		"account-group":    accountGroup,
		"account-category": accountCategory,
	}
	response := self.ApiFunc("accountCategoryGetOrderOpen", self.Extend(request, params), nil, nil)
	data := self.SafeValue(response, "data", []interface{}{})
	if self.ToBool(accountCategory == "futures") {
		return self.ToOrders(self.ParseOrders(data, market, since, limit)), nil
	}
	orders := []interface{}{}
	for i := 0; i < self.Length(data); i++ {
		order := self.ParseOrder(self.Member(data, i), market)
		orders = append(orders, order)
	}
	return self.ToOrders(self.FilterBySymbolSinceLimit(orders, symbol, since, limit)), nil
}

func (self *Bitmax2) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if self.ToBool(self.TestNil(symbol)) {
		self.RaiseException("ArgumentsRequired", self.Id+" cancelOrder() requires a symbol argument")
	}
	self.LoadMarkets()
	self.LoadAccounts()
	market := self.Market(symbol)
	defaultAccountCategory := self.SafeString(self.Options, "account-category", "cash")
	options := self.SafeValue(self.Options, "cancelOrder", map[string]interface{}{})
	accountCategory := self.SafeString(options, "account-category", defaultAccountCategory)
	accountCategory = self.SafeString(params, "account-category", accountCategory)
	params = self.Omit(params, "account-category")
	account := self.SafeValue(self.Accounts, 0, map[string]interface{}{})
	accountGroup := self.SafeValue(account, "id", nil)
	clientOrderId := self.SafeString2(params, "clientOrderId", "id", "")
	request := map[string]interface{}{
		"account-group":    accountGroup,
		"account-category": accountCategory,
		"symbol":           market.Id,
		"time":             self.Milliseconds(),
		"id":               "foobar",
	}
	if self.ToBool(self.TestNil(clientOrderId)) {
		self.SetValue(request, "orderId", id)
	} else {
		self.SetValue(request, "id", clientOrderId)
		params = self.Omit(params, []interface{}{"clientOrderId", "id"})
	}
	response = self.ApiFunc("accountCategoryDeleteOrder", self.Extend(request, params), nil, nil)
	data := self.SafeValue(response, "data", map[string]interface{}{})
	info := self.SafeValue(data, "info", map[string]interface{}{})
	return self.ParseOrder(info, market), nil
}

func (self *Bitmax2) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	url := ""
	query := params
	accountCategory := api == "accountCategory"
	if self.ToBool(accountCategory || api == "accountGroup") {
		url += self.ImplodeParams("/{account-group}", params)
		query = self.Omit(params, "account-group")
	}
	request := self.ImplodeParams(path, query)
	url += "/api/pro/" + self.Version
	if self.ToBool(accountCategory) {
		url += self.ImplodeParams("/{account-category}", query)
		query = self.Omit(query, "account-category")
	}
	url += "/" + request
	query = self.Omit(query, self.ExtractParams(path))
	if self.ToBool(api == "public") {
		if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
			url += "?" + self.Urlencode(query)
		}
	} else {
		self.CheckRequiredCredentials()
		timestamp := fmt.Sprintf("%v", self.Milliseconds())
		payload := timestamp + "+" + request
		hmac := self.Hmac(self.Encode(payload), self.Encode(self.Secret), "sha256", "base64")
		headers = map[string]interface{}{
			"x-auth-key":       self.ApiKey,
			"x-auth-timestamp": timestamp,
			"x-auth-signature": hmac,
		}
		if self.ToBool(method == "GET") {
			if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
				url += "?" + self.Urlencode(query)
			}
		} else {
			self.SetValue(headers, "Content-Type", "application/json")
			body = self.Json(query)
		}
	}
	url = self.Urls["api"].(string) + url
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Bitmax2) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if self.ToBool(self.TestNil(response)) {
		return
	}
	code := self.SafeString(response, "code", "")
	message := self.SafeString(response, "message", "")
	error := !self.TestNil(code) && code != "0"
	if self.ToBool(error || !self.TestNil(message)) {
		feedback := self.Id + " " + body
		self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), code, feedback)
		self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), message, feedback)
		self.ThrowBroadlyMatchedException(self.Member(self.Exceptions, "broad"), message, feedback)
		self.RaiseException("ExchangeError", feedback)
	}
}

func (self *Bitmax2) LoadMarkets() map[string]*Market {
	return nil
}

func (self *Bitmax2) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	return &Market{
		Id:     symbol,
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}
