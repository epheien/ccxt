package kucoin

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"reflect"
	"strings"
)

type Kucoin struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Kucoin, err error) {
	ex = new(Kucoin)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		return
	}

	return
}

func (self *Kucoin) InitDescribe() (err error) {
	err = json.Unmarshal(self.Child.Describe(), &self.DescribeMap)
	if err != nil {
		return
	}

	err = self.DefineRestApi()
	if err != nil {
		return
	}

	publicUrl, err := NestedMapLookup(self.DescribeMap, "urls", "api", "public")
	if err != nil {
		return
	}
	privateUrl, err := NestedMapLookup(self.DescribeMap, "urls", "api", "private")
	if err != nil {
		return
	}
	self.ApiUrls = map[string]string{
		"private": privateUrl.(string),
		"public":  publicUrl.(string),
	}

	self.Options = self.DescribeMap["options"].(map[string]interface{})
	self.Urls = self.DescribeMap["urls"].(map[string]interface{})
	return
}

func (self *Kucoin) Describe() []byte {
	return []byte(`{
    "id": "kucoin",
    "name": "KuCoin",
    "countries": [
        "SC"
    ],
    "rateLimit": 334,
    "version": "v2",
    "certified": false,
    "pro": true,
    "comment": "Platform 2.0",
    "has": {
        "CORS": false,
        "fetchStatus": true,
        "fetchTime": true,
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
        "fetchLedger": true
    },
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/57369448-3cc3aa80-7196-11e9-883e-5ebeb35e4f57.jpg",
        "referral": "https://www.kucoin.com/?rcode=E5wkqe",
        "api": {
            "public": "https://openapi-v2.kucoin.com",
            "private": "https://openapi-v2.kucoin.com"
        },
        "test": {
            "public": "https://openapi-sandbox.kucoin.com",
            "private": "https://openapi-sandbox.kucoin.com"
        },
        "www": "https://www.kucoin.com",
        "doc": [
            "https://docs.kucoin.com"
        ]
    },
    "requiredCredentials": {
        "apiKey": true,
        "secret": true,
        "password": true
    },
    "api": {
        "public": {
            "get": [
                "timestamp",
                "status",
                "symbols",
                "markets",
                "market/allTickers",
                "market/orderbook/level{level}",
                "market/orderbook/level2",
                "market/orderbook/level2_20",
                "market/orderbook/level2_100",
                "market/orderbook/level3",
                "market/histories",
                "market/candles",
                "market/stats",
                "currencies",
                "currencies/{currency}",
                "prices",
                "mark-price/{symbol}/current",
                "margin/config"
            ],
            "post": [
                "bullet-public"
            ]
        },
        "private": {
            "get": [
                "accounts",
                "accounts/{accountId}",
                "accounts/{accountId}/ledgers",
                "accounts/{accountId}/holds",
                "accounts/transferable",
                "sub/user",
                "sub-accounts",
                "sub-accounts/{subUserId}",
                "deposit-addresses",
                "deposits",
                "hist-deposits",
                "hist-orders",
                "hist-withdrawals",
                "withdrawals",
                "withdrawals/quotas",
                "orders",
                "orders/{orderId}",
                "limit/orders",
                "fills",
                "limit/fills",
                "margin/account",
                "margin/borrow",
                "margin/borrow/outstanding",
                "margin/borrow/borrow/repaid",
                "margin/lend/active",
                "margin/lend/done",
                "margin/lend/trade/unsettled",
                "margin/lend/trade/settled",
                "margin/lend/assets",
                "margin/market",
                "margin/margin/trade/last"
            ],
            "post": [
                "accounts",
                "accounts/inner-transfer",
                "accounts/sub-transfer",
                "deposit-addresses",
                "withdrawals",
                "orders",
                "orders/multi",
                "margin/borrow",
                "margin/repay/all",
                "margin/repay/single",
                "margin/lend",
                "margin/toggle-auto-lend",
                "bullet-private"
            ],
            "delete": [
                "withdrawals/{withdrawalId}",
                "orders",
                "orders/{orderId}",
                "margin/lend/{orderId}"
            ]
        }
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
        "1w": "1week"
    },
    "exceptions": {
        "exact": {
            "order not exist": "OrderNotFound",
            "order not exist.": "OrderNotFound",
            "order_not_exist": "OrderNotFound",
            "order_not_exist_or_not_allow_to_cancel": "InvalidOrder",
            "Order size below the minimum requirement.": "InvalidOrder",
            "The withdrawal amount is below the minimum requirement.": "ExchangeError",
            "400": "BadRequest",
            "401": "AuthenticationError",
            "403": "NotSupported",
            "404": "NotSupported",
            "405": "NotSupported",
            "429": "RateLimitExceeded",
            "500": "ExchangeError",
            "503": "ExchangeNotAvailable",
            "200004": "InsufficientFunds",
            "230003": "InsufficientFunds",
            "260100": "InsufficientFunds",
            "300000": "InvalidOrder",
            "400000": "BadSymbol",
            "400001": "AuthenticationError",
            "400002": "InvalidNonce",
            "400003": "AuthenticationError",
            "400004": "AuthenticationError",
            "400005": "AuthenticationError",
            "400006": "AuthenticationError",
            "400007": "AuthenticationError",
            "400008": "NotSupported",
            "400100": "BadRequest",
            "411100": "AccountSuspended",
            "415000": "BadRequest",
            "500000": "ExchangeError"
        },
        "broad": {
            "Exceeded the access frequency": "RateLimitExceeded"
        }
    },
    "fees": {
        "trading": {
            "tierBased": false,
            "percentage": true,
            "taker": 0.001,
            "maker": 0.001
        },
        "funding": {
            "tierBased": false,
            "percentage": false,
            "withdraw": {},
            "deposit": {}
        }
    },
    "commonCurrencies": {
        "HOT": "HOTNOW",
        "EDGE": "DADI",
        "WAX": "WAXP",
        "TRY": "Trias"
    },
    "options": {
        "version": "v1",
        "symbolSeparator": "-",
        "fetchMyTradesMethod": "private_get_fills",
        "fetchBalance": {
            "type": "trade"
        },
        "versions": {
            "public": {
                "GET": {
                    "status": "v1",
                    "market/orderbook/level{level}": "v1",
                    "market/orderbook/level2": "v2",
                    "market/orderbook/level2_20": "v1",
                    "market/orderbook/level2_100": "v1"
                }
            },
            "private": {
                "POST": {
                    "accounts/inner-transfer": "v2",
                    "accounts/sub-transfer": "v2"
                }
            }
        }
    }
}`)
}

func (self *Kucoin) FetchMarkets(params map[string]interface{}) (ret interface{}, err error) {

	response, err := self.ApiFunc("publicGetSymbols", params, nil, nil)
	if err != nil {
		return nil, err
	}

	data := self.Member(response, "data")

	result := []interface{}{}

	for i := 0; i < self.Length(data); i++ {
		market := self.Member(data, i)
		id := self.SafeString(market, "symbol", "")
		baseId, quoteId := self.Unpack2(strings.Split(id, "-"))
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		active := self.SafeValue(market, "enableTrading", nil)
		baseMaxSize := self.SafeFloat(market, "baseMaxSize", 0)
		baseMinSize := self.SafeFloat(market, "baseMinSize", 0)
		quoteMaxSize := self.SafeFloat(market, "quoteMaxSize", 0)
		quoteMinSize := self.SafeFloat(market, "quoteMinSize", 0)
		precision := map[string]interface{}{
			"amount": self.PrecisionFromString(self.SafeString(market, "baseIncrement", "")),
			"price":  self.PrecisionFromString(self.SafeString(market, "priceIncrement", "")),
		}
		limits := map[string]interface{}{
			"amount": map[string]interface{}{
				"min": baseMinSize,
				"max": baseMaxSize,
			},
			"price": map[string]interface{}{
				"min": self.SafeFloat(market, "priceIncrement", 0),
				"max": quoteMaxSize / baseMinSize,
			},
			"cost": map[string]interface{}{
				"min": quoteMinSize,
				"max": quoteMaxSize,
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

	return result, nil

}

func (self *Kucoin) FetchCurrencies(params map[string]interface{}) (ret interface{}, err error) {

	response, err := self.ApiFunc("publicGetCurrencies", params, nil, nil)
	if err != nil {
		return nil, err
	}

	responseData := self.Member(response, "data")

	result := map[string]interface{}{}

	for i := 0; i < self.Length(responseData); i++ {
		entry := self.Member(responseData, i)
		id := self.SafeString(entry, "currency", "")
		name := self.SafeString(entry, "fullName", "")
		code := self.SafeCurrencyCode(id)
		precision := self.SafeInteger(entry, "precision", 0)
		self.SetValue(result, code, map[string]interface{}{
			"id":        id,
			"name":      name,
			"code":      code,
			"precision": precision,
			"info":      entry,
			"active":    nil,
			"fee":       nil,
			"limits":    self.Limits,
		})
	}

	return result, nil

}

func (self *Kucoin) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook interface{}, err error) {

	level := self.SafeInteger(params, "level", 2)

	levelLimit := fmt.Sprintf("%v", level)

	if self.ToBool(levelLimit == "2") {
		if self.ToBool(!self.TestNil(limit)) {
			if self.ToBool(limit != 20 && limit != 100) {
				err = errors.New(self.Id + " fetchOrderBook limit argument must be undefined, 20 or 100")
				return
			}
			levelLimit += "_" + fmt.Sprintf("%v", limit)
		}
	}

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	marketId, err := self.MarketId(symbol)
	if err != nil {
		return nil, err
	}

	request := map[string]interface{}{
		"symbol": marketId,
		"level":  levelLimit,
	}

	response, err := self.ApiFunc("publicGetMarketOrderbookLevelLevel", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	data := self.SafeValue(response, "data", map[string]interface{}{})

	timestamp := self.SafeInteger(data, "time", 0)

	orderbook, err := self.ParseOrderBook(data, timestamp, "bids", "asks", level-2, level-1)
	if err != nil {
		return nil, err
	}

	self.SetValue(orderbook, "nonce", self.SafeInteger(data, "sequence", 0))

	return orderbook, nil

}

func (self *Kucoin) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (order map[string]interface{}, err error) {

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	marketId, err := self.MarketId(symbol)
	if err != nil {
		return nil, err
	}

	clientOrderId := self.SafeString2(params, "clientOid", "clientOrderId", self.Uuid())

	params = self.Omit(params, []interface{}{"clientOid", "clientOrderId"})

	request := map[string]interface{}{
		"clientOid": clientOrderId,
		"side":      side,
		"symbol":    marketId,
		"type":      typ,
	}

	if self.ToBool(typ != "market") {
		self.SetValue(request, "price", self.PriceToPrecision(symbol, price))
		self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
	} else {
		if self.ToBool(self.SafeValue(params, "quoteAmount", nil)) {
			self.SetValue(request, "funds", self.AmountToPrecision(symbol, amount))
		} else {
			self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
		}
	}

	response, err := self.ApiFunc("privatePostOrders", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	data := self.SafeValue(response, "data", map[string]interface{}{})

	timestamp := self.Milliseconds()

	order = map[string]interface{}{
		"id":            self.SafeString(data, "orderId", ""),
		"symbol":        symbol,
		"type":          typ,
		"side":          side,
		"price":         price,
		"cost":          nil,
		"filled":        nil,
		"remaining":     nil,
		"timestamp":     timestamp,
		"datetime":      self.Iso8601(timestamp),
		"fee":           nil,
		"status":        "open",
		"clientOrderId": clientOrderId,
		"info":          data,
	}

	if self.ToBool(!self.ToBool(self.SafeValue(params, "quoteAmount", nil))) {
		self.SetValue(order, "amount", amount)
	}

	return order, nil

}

func (self *Kucoin) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {

	request := map[string]interface{}{
		"orderId": id,
	}

	response, err = self.ApiFunc("privateDeleteOrdersOrderId", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (self *Kucoin) FetchOrdersByStatus(status string, symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}, err error) {

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"status": status,
	}

	var market interface{}

	if self.ToBool(!self.TestNil(symbol)) {
		market = self.Market(symbol)
		self.SetValue(request, "symbol", self.Member(market, "id"))
	}

	if self.ToBool(!self.TestNil(since)) {
		self.SetValue(request, "startAt", since)
	}

	if self.ToBool(!self.TestNil(limit)) {
		self.SetValue(request, "pageSize", limit)
	}

	response, err := self.ApiFunc("privateGetOrders", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	responseData := self.SafeValue(response, "data", map[string]interface{}{})

	orders = self.SafeValue(responseData, "items", []interface{}{})

	return self.ParseOrders(orders, market, since, limit), nil

}

func (self *Kucoin) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}, err error) {

	return self.FetchOrdersByStatus("active", symbol, since, limit, params)

}

func (self *Kucoin) FetchOrder(id string, symbol string, params map[string]interface{}) (order interface{}, err error) {

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"orderId": id,
	}

	var market interface{}

	if self.ToBool(!self.TestNil(symbol)) {
		market = self.Market(symbol)
	}

	response, err := self.ApiFunc("privateGetOrdersOrderId", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	responseData := self.Member(response, "data")

	return self.ParseOrder(responseData, market), nil

}

func (self *Kucoin) ParseOrder(order interface{}, market interface{}) (result interface{}) {

	var symbol interface{}

	marketId := self.SafeString(order, "symbol", "")

	if self.ToBool(!self.TestNil(marketId)) {
		if self.ToBool(self.InMap(marketId, self.MarketsById)) {
			market = self.Member(self.MarketsById, marketId)
			symbol = self.Member(market, "symbol")
		} else {
			baseId, quoteId := self.Unpack2(strings.Split(marketId, "-"))
			base := self.SafeCurrencyCode(baseId)
			quote := self.SafeCurrencyCode(quoteId)
			symbol = base + "/" + quote
		}
		market = self.SafeValue(self.MarketsById, marketId, nil)
	}

	if self.ToBool(self.TestNil(symbol)) {
		if self.ToBool(!self.TestNil(market)) {
			symbol = self.Member(market, "symbol")
		}
	}

	orderId := self.SafeString(order, "id", "")

	typ := self.SafeString(order, "type", "")

	timestamp := self.SafeInteger(order, "createdAt", 0)

	datetime := self.Iso8601(timestamp)

	price := self.SafeFloat(order, "price", 0)

	side := self.SafeString(order, "side", "")

	feeCurrencyId := self.SafeString(order, "feeCurrency", "")

	feeCurrency := self.SafeCurrencyCode(feeCurrencyId)

	feeCost := self.SafeFloat(order, "fee", 0)

	amount := self.SafeFloat(order, "size", 0)

	filled := self.SafeFloat(order, "dealSize", 0)

	cost := self.SafeFloat(order, "dealFunds", 0)

	remaining := amount - filled

	status := self.IfThenElse(self.ToBool(self.Member(order, "isActive")), "open", "closed")

	status = self.IfThenElse(self.ToBool(self.Member(order, "cancelExist")), "canceled", status)

	fee := map[string]interface{}{
		"currency": feeCurrency,
		"cost":     feeCost,
	}

	if self.ToBool(typ == "market") {
		if self.ToBool(price == 0.0) {
			if self.ToBool(!self.TestNil(cost) && !self.TestNil(filled)) {
				if self.ToBool(cost > 0 && filled > 0) {
					price = cost / filled
				}
			}
		}
	}

	clientOrderId := self.SafeString(order, "clientOid", "")

	return map[string]interface{}{
		"id":                 orderId,
		"clientOrderId":      clientOrderId,
		"symbol":             symbol,
		"type":               typ,
		"side":               side,
		"amount":             amount,
		"price":              price,
		"cost":               cost,
		"filled":             filled,
		"remaining":          remaining,
		"timestamp":          timestamp,
		"datetime":           datetime,
		"fee":                fee,
		"status":             status,
		"info":               order,
		"lastTradeTimestamp": nil,
		"average":            nil,
		"trades":             nil,
	}

}

func (self *Kucoin) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	var typ interface{}

	request := map[string]interface{}{}

	if self.ToBool(self.InMap("type", params)) {
		typ = self.Member(params, "type")
		if self.ToBool(!self.TestNil(typ)) {
			self.SetValue(request, "type", typ)
		}
		params = self.Omit(params, "type")
	} else {
		options := self.SafeValue(self.Options, "fetchBalance", map[string]interface{}{})
		typ = self.SafeString(options, "type", "trade")
	}

	response, err := self.ApiFunc("privateGetAccounts", self.Extend(request, params), nil, nil)
	if err != nil {
		return nil, err
	}

	data := self.SafeValue(response, "data", []interface{}{})

	result := map[string]interface{}{
		"info": response,
	}

	for i := 0; i < self.Length(data); i++ {
		balance := self.Member(data, i)
		balanceType := self.SafeString(balance, "type", "")
		if self.ToBool(balanceType == typ) {
			currencyId := self.SafeString(balance, "currency", "")
			code := self.SafeCurrencyCode(currencyId)
			account := self.Account()
			self.SetValue(account, "total", self.SafeFloat(balance, "balance", 0))
			self.SetValue(account, "free", self.SafeFloat(balance, "available", 0))
			self.SetValue(account, "used", self.SafeFloat(balance, "holds", 0))
			self.SetValue(result, code, account)
		}
	}

	return self.ParseBalance(result), nil

}

func (self *Kucoin) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}, err error) {

	versions := self.SafeValue(self.Options, "versions", map[string]interface{}{})

	apiVersions := self.SafeValue(versions, api, nil)

	methodVersions := self.SafeValue(apiVersions, method, map[string]interface{}{})

	defaultVersion := self.SafeString(methodVersions, path, self.Member(self.Options, "version"))

	version := self.SafeString(params, "version", defaultVersion)

	params = self.Omit(params, "version")

	endpoint := "/api/" + version + "/" + self.ImplodeParams(path, params)

	query := self.Omit(params, self.ExtractParams(path))

	endpart := ""

	headers = self.IfThenElse(self.ToBool(!self.TestNil(headers)), headers, map[string]interface{}{})

	if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
		if self.ToBool(method != "GET") {
			endpart = self.Json(query)
			self.SetValue(headers, "Content-Type", "application/json")
		} else {
			endpoint += "?" + self.Urlencode(query)
		}
	}

	url := fmt.Sprintf("%v", self.Member(self.Member(self.Urls, "api"), api)) + endpoint

	if self.ToBool(api == "private") {
		self.CheckRequiredCredentials()
		timestamp := fmt.Sprintf("%v", self.Nonce())
		headers = self.Extend(map[string]interface{}{
			"KC-API-KEY":        self.ApiKey,
			"KC-API-TIMESTAMP":  timestamp,
			"KC-API-PASSPHRASE": self.Password,
		}, headers)
		payload := timestamp + method + endpoint + endpart
		signature, err := self.Hmac(self.Encode(payload), self.Encode(self.Secret), "sha256", "base64")
		if err != nil {
			return nil, err
		}
		self.SetValue(headers, "KC-API-SIGN", self.Decode(signature))
	}

	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    endpart,
		"headers": headers,
	}, nil

}