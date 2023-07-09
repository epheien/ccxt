package binance

import (
	"encoding/json"
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"math"
	"reflect"
	"strings"
)

type Binance struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Binance, err error) {
	ex = new(Binance)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Binance) Describe() []byte {
	return []byte(`
{
    "id": "binance",
    "name": "Binance",
    "countries": [
        "JP",
        "MT"
    ],
    "rateLimit": 500,
    "certified": true,
    "pro": true,
    "has": {
        "fetchDepositAddress": true,
        "CORS": false,
        "fetchBidsAsks": true,
        "fetchTickers": true,
        "fetchTime": true,
        "fetchOHLCV": true,
        "fetchMyTrades": true,
        "fetchOrder": true,
        "fetchOrders": true,
        "fetchOpenOrders": true,
        "fetchClosedOrders": "emulated",
        "withdraw": true,
        "fetchFundingFees": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchTransactions": false,
        "fetchTradingFee": true,
        "fetchTradingFees": true,
        "cancelAllOrders": true
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
        "test": {
            "fapiPublic": "https://testnet.binancefuture.com/fapi/v1",
            "fapiPrivate": "https://testnet.binancefuture.com/fapi/v1",
            "public": "https://testnet.binance.vision/api/v3",
            "private": "https://testnet.binance.vision/api/v3",
            "v3": "https://testnet.binance.vision/api/v3",
            "v1": "https://testnet.binance.vision/api/v1"
        },
        "api": {
            "wapi": "https://api.binance.com/wapi/v3",
            "sapi": "https://api.binance.com/sapi/v1",
            "fapiPublic": "https://fapi.binance.com/fapi/v1",
            "fapiPrivate": "https://fapi.binance.com/fapi/v1",
            "public": "https://api.binance.com/api/v3",
            "private": "https://api.binance.com/api/v3",
            "v3": "https://api.binance.com/api/v3",
            "v1": "https://api.binance.com/api/v1"
        },
        "www": "https://www.binance.com",
        "referral": "https://www.binance.com/?ref=10205187",
        "doc": [
            "https://binance-docs.github.io/apidocs/spot/en"
        ],
        "api_management": "https://www.binance.com/en/usercenter/settings/api-management",
        "fees": "https://www.binance.com/en/fee/schedule"
    },
    "api": {
        "sapi": {
            "get": [
                "accountSnapshot",
                "margin/asset",
                "margin/pair",
                "margin/allAssets",
                "margin/allPairs",
                "margin/priceIndex",
                "asset/assetDividend",
                "margin/loan",
                "margin/repay",
                "margin/account",
                "margin/transfer",
                "margin/interestHistory",
                "margin/forceLiquidationRec",
                "margin/order",
                "margin/openOrders",
                "margin/allOrders",
                "margin/myTrades",
                "margin/maxBorrowable",
                "margin/maxTransferable",
                "futures/transfer",
                "capital/config/getall",
                "capital/deposit/address",
                "capital/deposit/hisrec",
                "capital/deposit/subAddress",
                "capital/deposit/subHisrec",
                "capital/withdraw/history",
                "sub-account/futures/account",
                "sub-account/futures/accountSummary",
                "sub-account/futures/positionRisk",
                "sub-account/margin/account",
                "sub-account/margin/accountSummary",
                "sub-account/status",
                "sub-account/transfer/subUserHistory",
                "lending/daily/product/list",
                "lending/daily/userLeftQuota",
                "lending/daily/userRedemptionQuota",
                "lending/daily/token/position",
                "lending/union/account",
                "lending/union/purchaseRecord",
                "lending/union/redemptionRecord",
                "lending/union/interestHistory",
                "lending/project/list",
                "lending/project/position/list",
                "mining/pub/algoList",
                "mining/pub/coinList",
                "mining/worker/detail",
                "mining/worker/list",
                "mining/payment/list",
                "mining/statistics/user/status",
                "mining/statistics/user/list"
            ],
            "post": [
                "asset/dust",
                "account/disableFastWithdrawSwitch",
                "account/enableFastWithdrawSwitch",
                "capital/withdraw/apply",
                "margin/transfer",
                "margin/loan",
                "margin/repay",
                "margin/order",
                "sub-account/margin/enable",
                "sub-account/margin/enable",
                "sub-account/futures/enable",
                "userDataStream",
                "futures/transfer",
                "lending/customizedFixed/purchase",
                "lending/daily/purchase",
                "lending/daily/redeem"
            ],
            "put": [
                "userDataStream"
            ],
            "delete": [
                "margin/order",
                "userDataStream"
            ]
        },
        "wapi": {
            "post": [
                "withdraw",
                "sub-account/transfer"
            ],
            "get": [
                "depositHistory",
                "withdrawHistory",
                "depositAddress",
                "accountStatus",
                "systemStatus",
                "apiTradingStatus",
                "userAssetDribbletLog",
                "tradeFee",
                "assetDetail",
                "sub-account/list",
                "sub-account/transfer/history",
                "sub-account/assets"
            ]
        },
        "fapiPublic": {
            "get": [
                "ping",
                "time",
                "exchangeInfo",
                "depth",
                "trades",
                "historicalTrades",
                "aggTrades",
                "klines",
                "fundingRate",
                "premiumIndex",
                "ticker/24hr",
                "ticker/price",
                "ticker/bookTicker",
                "allForceOrders",
                "openInterest",
                "leverageBracket"
            ]
        },
        "fapiPrivate": {
            "get": [
                "allForceOrders",
                "allOrders",
                "openOrder",
                "openOrders",
                "order",
                "account",
                "balance",
                "positionMargin/history",
                "positionRisk",
                "positionSide/dual",
                "userTrades",
                "income"
            ],
            "post": [
                "batchOrders",
                "positionSide/dual",
                "positionMargin",
                "marginType",
                "order",
                "leverage",
                "listenKey",
                "countdownCancelAll"
            ],
            "put": [
                "listenKey"
            ],
            "delete": [
                "batchOrders",
                "order",
                "allOpenOrders",
                "listenKey"
            ]
        },
        "v3": {
            "get": [
                "ticker/price",
                "ticker/bookTicker"
            ]
        },
        "public": {
            "get": [
                "ping",
                "time",
                "depth",
                "trades",
                "aggTrades",
                "historicalTrades",
                "klines",
                "ticker/24hr",
                "ticker/price",
                "ticker/bookTicker",
                "exchangeInfo"
            ],
            "put": [
                "userDataStream"
            ],
            "post": [
                "userDataStream"
            ],
            "delete": [
                "userDataStream"
            ]
        },
        "private": {
            "get": [
                "allOrderList",
                "openOrderList",
                "orderList",
                "order",
                "openOrders",
                "allOrders",
                "account",
                "myTrades"
            ],
            "post": [
                "order/oco",
                "order",
                "order/test"
            ],
            "delete": [
                "openOrders",
                "orderList",
                "order"
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
    "commonCurrencies": {
        "BCC": "BCC",
        "YOYO": "YOYOW"
    },
    "options": {
        "fetchTradesMethod": "publicGetAggTrades",
        "fetchTickersMethod": "publicGetTicker24hr",
        "defaultTimeInForce": "GTC",
        "defaultType": "spot",
        "hasAlreadyAuthenticatedSuccessfully": false,
        "warnOnFetchOpenOrdersWithoutSymbol": true,
        "recvWindow": 5000,
        "timeDifference": 0,
        "adjustForTimeDifference": false,
        "parseOrderToPrecision": false,
        "newOrderRespType": {
            "market": "FULL",
            "limit": "RESULT"
        },
        "quoteOrderQty": true
    },
    "exceptions": {
        "API key does not exist": "AuthenticationError",
        "Order would trigger immediately.": "InvalidOrder",
        "Account has insufficient balance for requested action.": "InsufficientFunds",
        "Rest API trading is not enabled.": "ExchangeNotAvailable",
        "You don't have permission.": "PermissionDenied",
        "Market is closed.": "ExchangeNotAvailable",
        "-1000": "ExchangeNotAvailable",
        "-1003": "RateLimitExceeded",
        "-1013": "InvalidOrder",
        "-1021": "InvalidNonce",
        "-1022": "AuthenticationError",
        "-1100": "InvalidOrder",
        "-1104": "ExchangeError",
        "-1128": "ExchangeError",
        "-2010": "ExchangeError",
        "-2011": "OrderNotFound",
        "-2013": "OrderNotFound",
        "-2014": "AuthenticationError",
        "-2015": "AuthenticationError",
        "-3008": "InsufficientFunds",
        "-3010": "ExchangeError"
    }
}
`)
}

func (self *Binance) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	defaultType := self.SafeString2(self.Options, "fetchMarkets", "defaultType", "spot")
	typ := self.SafeString(params, "type", defaultType)
	query := self.Omit(params, "type")
	if self.ToBool(typ != "spot" && typ != "future" && typ != "margin") {
		self.RaiseException("ExchangeError", self.Id+" does not support "+typ+" type, set exchange.options[defaultType] to spot, margin or future")
	}
	method := self.IfThenElse(self.ToBool(typ == "future"), "fapiPublicGetExchangeInfo", "publicGetExchangeInfo").(string)
	response := self.ApiFunc(method, query, nil, nil)
	if self.ToBool(self.Member(self.Options, "adjustForTimeDifference")) {
		// TODO, false
		//self.LoadTimeDifference()
	}
	markets := self.SafeValue(response, "symbols", nil)
	result := []interface{}{}
	for i := 0; i < self.Length(markets); i++ {
		market := self.Member(markets, i)
		future := self.InMap("maintMarginPercent", market)
		spot := !self.ToBool(future)
		marketType := self.IfThenElse(self.ToBool(spot), "spot", "future")
		id := self.SafeString(market, "symbol", "")
		lowercaseId := self.SafeStringLower(market, "symbol", "")
		baseId := self.SafeString(market, "baseAsset", "")
		quoteId := self.SafeString(market, "quoteAsset", "")
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		filters := self.SafeValue(market, "filters", []interface{}{})
		filtersByType := self.IndexBy(filters, "filterType")
		precision := map[string]interface{}{
			"base":   self.SafeInteger(market, "baseAssetPrecision", 0),
			"quote":  self.SafeInteger(market, "quotePrecision", 0),
			"amount": self.SafeInteger(market, "baseAssetPrecision", 0),
			"price":  self.SafeInteger(market, "quotePrecision", 0),
		}
		status := self.SafeString(market, "status", "")
		active := status == "TRADING"
		margin := self.SafeValue(market, "isMarginTradingAllowed", future)
		entry := map[string]interface{}{
			"id":          id,
			"lowercaseId": lowercaseId,
			"symbol":      symbol,
			"base":        base,
			"quote":       quote,
			"baseId":      baseId,
			"quoteId":     quoteId,
			"info":        market,
			"type":        marketType,
			"spot":        spot,
			"margin":      margin,
			"future":      future,
			"active":      active,
			"precision":   precision,
			"limits": map[string]interface{}{
				"amount": map[string]interface{}{
					"min": math.Pow10(int(-precision["amount"].(int64))),
					"max": nil,
				},
				"price": map[string]interface{}{
					"min": nil,
					"max": nil,
				},
				"cost": map[string]interface{}{
					"min": nil,
					"max": nil,
				},
			},
		}
		if self.ToBool(self.InMap("PRICE_FILTER", filtersByType)) {
			filter := self.SafeValue(filtersByType, "PRICE_FILTER", map[string]interface{}{})
			self.SetValue(self.Member(entry, "limits"), "price", map[string]interface{}{
				"min": self.SafeFloat(filter, "minPrice", 0),
				"max": nil,
			})
			maxPrice := self.SafeFloat(filter, "maxPrice", 0)
			if self.ToBool(!self.TestNil(maxPrice) && maxPrice > 0) {
				self.SetValue(self.Member(self.Member(entry, "limits"), "price"), "max", maxPrice)
			}
			self.SetValue(self.Member(entry, "precision"), "price", self.PrecisionFromString(self.Member(filter, "tickSize").(string)))
		}
		if self.ToBool(self.InMap("LOT_SIZE", filtersByType)) {
			filter := self.SafeValue(filtersByType, "LOT_SIZE", map[string]interface{}{})
			stepSize := self.SafeString(filter, "stepSize", "")
			self.SetValue(self.Member(entry, "precision"), "amount", self.PrecisionFromString(stepSize))
			self.SetValue(self.Member(entry, "limits"), "amount", map[string]interface{}{
				"min": self.SafeFloat(filter, "minQty", 0),
				"max": self.SafeFloat(filter, "maxQty", 0),
			})
		}
		if self.ToBool(self.InMap("MARKET_LOT_SIZE", filtersByType)) {
			filter := self.SafeValue(filtersByType, "MARKET_LOT_SIZE", map[string]interface{}{})
			self.SetValue(self.Member(entry, "limits"), "market", map[string]interface{}{
				"min": self.SafeFloat(filter, "minQty", 0),
				"max": self.SafeFloat(filter, "maxQty", 0),
			})
		}
		if self.ToBool(self.InMap("MIN_NOTIONAL", filtersByType)) {
			filter := self.SafeValue(filtersByType, "MIN_NOTIONAL", map[string]interface{}{})
			self.SetValue(self.Member(self.Member(entry, "limits"), "cost"), "min", self.SafeFloat(filter, "minNotional", 0))
		}
		result = append(result, entry)
	}
	return self.ToMarkets(result), nil
}

func (self *Binance) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	defaultType := self.SafeString2(self.Options, "fetchBalance", "defaultType", "spot")
	typ := self.SafeString(params, "type", defaultType)
	method := "privateGetAccount"
	if self.ToBool(typ == "future") {
		method = "fapiPrivateGetAccount"
	} else if self.ToBool(typ == "margin") {
		method = "sapiGetMarginAccount"
	}
	query := self.Omit(params, "type")
	response := self.ApiFunc(method, query, nil, nil)
	result := map[string]interface{}{
		"info": response,
	}
	if self.ToBool(typ == "spot" || typ == "margin") {
		balances := self.SafeValue2(response, "balances", "userAssets", []interface{}{})
		for i := 0; i < self.Length(balances); i++ {
			balance := self.Member(balances, i)
			currencyId := self.SafeString(balance, "asset", "")
			code := self.SafeCurrencyCode(currencyId)
			account := self.Account()
			self.SetValue(account, "free", self.SafeFloat(balance, "free", 0))
			self.SetValue(account, "used", self.SafeFloat(balance, "locked", 0))
			self.SetValue(result, code, account)
		}
	} else {
		balances := self.SafeValue(response, "assets", []interface{}{})
		for i := 0; i < self.Length(balances); i++ {
			balance := self.Member(balances, i)
			currencyId := self.SafeString(balance, "asset", "")
			code := self.SafeCurrencyCode(currencyId)
			account := self.Account()
			self.SetValue(account, "used", self.SafeFloat(balance, "initialMargin", 0))
			self.SetValue(account, "total", self.SafeFloat(balance, "marginBalance", 0))
			self.SetValue(result, code, account)
		}
	}
	return self.ParseBalance(result), nil
}

func (self *Binance) ParseTicker(response interface{}) (ticker *Ticker) {
	timestamp := self.SafeInteger(response, "closeTime")
	datetime := self.Iso8601(timestamp)
	last := self.SafeFloat(response, "lastPrice")
	ticker = &Ticker{
		Timestamp:   timestamp,
		Datetime:    datetime,
		Last:        last,
		Ask:         self.SafeFloat(response, "askPrice"),
		AskQty:      self.SafeFloat(response, "askQty"),
		Bid:         self.SafeFloat(response, "bidPrice"),
		BidQty:      self.SafeFloat(response, "bidQty"),
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

func (self *Binance) FetchTicker(symbol string, params map[string]interface{}) (ticker *Ticker, err error) {
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

func (self *Binance) ParseOHLCV(response interface{}) *OHLCV {
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

func (self *Binance) FetchOHLCV(symbol, timeframe string, since int64, limit int64, params map[string]interface{}) (klines []*OHLCV, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol":   self.Member(market, "id"),
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

func (self *Binance) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
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
	if self.ToBool(!self.TestNil(limit)) {
		self.SetValue(request, "limit", limit)
	}
	method := self.IfThenElse(self.ToBool(self.Member(market, "spot")), "publicGetDepth", "fapiPublicGetDepth").(string)
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	orderbook := self.ParseOrderBook(response, 0, "bids", "asks", 0, 1)
	self.SetValue(orderbook, "nonce", self.SafeInteger(response, "lastUpdateId", 0))
	return orderbook, nil
}

func (self *Binance) ParseOrderStatus(status string) string {
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

func (self *Binance) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	status := self.ParseOrderStatus(self.SafeString(order, "status", ""))
	var symbol interface{}
	marketId := self.SafeString(order, "symbol", "")
	if self.ToBool(self.InMap(marketId, self.MarketsById)) {
		market = self.Member(self.MarketsById, marketId)
	}
	if self.ToBool(!self.TestNil(market)) {
		symbol = self.Member(market, "symbol")
	}
	var timestamp interface{}
	if self.ToBool(self.InMap("time", order)) {
		timestamp = self.SafeInteger(order, "time", 0)
	} else if self.ToBool(self.InMap("transactTime", order)) {
		timestamp = self.SafeInteger(order, "transactTime", 0)
	}
	price := self.SafeFloat(order, "price", 0)
	amount := self.SafeFloat(order, "origQty", 0)
	filled := self.SafeFloat(order, "executedQty", 0)
	var remaining interface{}
	cost := self.SafeFloat2(order, "cummulativeQuoteQty", "cumQuote", 0.0)
	if self.ToBool(!self.TestNil(filled)) {
		if self.ToBool(!self.TestNil(amount)) {
			remaining = amount - filled
			if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
				remaining = ToFloat(self.AmountToPrecision(symbol.(string), remaining.(float64)))
			}
			remaining = math.Max(remaining.(float64), 0.0)
		}
		if self.ToBool(!self.TestNil(price)) {
			if self.ToBool(self.TestNil(cost)) {
				cost = price * filled
			}
		}
	}
	// NOTE: json 解析后的类型为 float64, 转为字符串的时候就不能显示为整数了
	id := fmt.Sprint(self.SafeInteger(order, "orderId", 0))
	typ := self.SafeStringLower(order, "type", "")
	if self.ToBool(typ == "market") {
		if self.ToBool(price == 0) {
			if self.ToBool(!self.TestNil(cost) && !self.TestNil(filled)) {
				if self.ToBool(cost > 0 && filled > 0) {
					price = cost / filled
					if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
						price = ToFloat(self.PriceToPrecision(symbol.(string), price))
					}
				}
			}
		}
	} else if self.ToBool(typ == "limit_maker") {
		typ = "limit"
	}
	side := self.SafeStringLower(order, "side", "")
	var fee interface{}
	var trades interface{}
	// TODO, fetchOrder返回的fee 没用到，没必要实现这么复杂
	/*
		fills := self.SafeValue(order, "fills", nil)
		if self.ToBool(!self.TestNil(fills)) {
			trades = self.ParseTrades(fills, market)
			numTrades := self.Length(trades)
			if self.ToBool(numTrades > 0) {
				cost = self.Member(self.Member(trades, 0), "cost").(float64)
				fee = map[string]interface{}{
					"cost":     self.Member(self.Member(self.Member(trades, 0), "fee"), "cost"),
					"currency": self.Member(self.Member(self.Member(trades, 0), "fee"), "currency"),
				}
				for i := 1; i < self.Length(trades); i++ {
					cost += self.Member(self.Member(trades, i), "cost").(float64)
					fee["cost"] = fee["cost"].(float64) + self.Member(self.Member(self.Member(trades, i), "fee"), "cost").(float64)
				}
			}
		}
	*/
	var average interface{}
	if self.ToBool(!self.TestNil(cost)) {
		if self.ToBool(filled) {
			average = cost / filled
			if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
				average = ToFloat(self.PriceToPrecision(symbol.(string), average.(float64)))
			}
		}
		if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
			//cost = ToFloat(self.CostToPrecision(symbol.(string), cost))
		}
	}
	clientOrderId := self.SafeString(order, "clientOrderId", "")
	return map[string]interface{}{
		"info":          order,
		"id":            id,
		"clientOrderId": clientOrderId,
		"timestamp":     timestamp,
		//"datetime":           self.Iso8601(timestamp.(int64)),
		"lastTradeTimestamp": nil,
		"symbol":             symbol,
		"type":               typ,
		"side":               side,
		"price":              price,
		"amount":             amount,
		"cost":               cost,
		"average":            average,
		"filled":             filled,
		"remaining":          remaining,
		"status":             status,
		"fee":                fee,
		"trades":             trades,
	}
}

func (self *Binance) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	defaultType := self.SafeString2(self.Options, "createOrder", "defaultType", market.Type)
	orderType := self.SafeString(params, "type", defaultType)
	clientOrderId := self.SafeString2(params, "newClientOrderId", "clientOrderId", "")
	params = self.Omit(params, []interface{}{"type", "newClientOrderId", "clientOrderId"})
	method := "privatePostOrder"
	if self.ToBool(orderType == "future") {
		method = "fapiPrivatePostOrder"
	} else if self.ToBool(orderType == "margin") {
		method = "sapiPostMarginOrder"
	}
	if self.ToBool(self.Member(market, "spot")) {
		test := self.SafeValue(params, "test", false)
		if self.ToBool(test) {
			method += "Test"
		}
		params = self.Omit(params, "test")
	}
	uppercaseType := strings.ToUpper(typ)
	request := map[string]interface{}{
		"symbol": self.Member(market, "id"),
		"type":   uppercaseType,
		"side":   strings.ToUpper(side),
	}
	if self.ToBool(!self.TestNil(clientOrderId)) {
		self.SetValue(request, "newClientOrderId", clientOrderId)
	}
	if self.ToBool(self.Member(market, "spot")) {
		self.SetValue(request, "newOrderRespType", self.SafeValue(self.Member(self.Options, "newOrderRespType"), typ, "RESULT"))
	}
	timeInForceIsRequired := false
	priceIsRequired := false
	stopPriceIsRequired := false
	quantityIsRequired := false
	if self.ToBool(uppercaseType == "MARKET") {
		quoteOrderQty := self.SafeValue(self.Options, "quoteOrderQty", false)
		if self.ToBool(quoteOrderQty) {
			quoteOrderQty := self.SafeFloat(params, "quoteOrderQty", 0)
			precision := self.Member(self.Member(market, "precision"), "price")
			if self.ToBool(!self.TestNil(quoteOrderQty)) {
				x, _ := DecimalToPrecision(quoteOrderQty, Truncate, precision.(int), DecimalPlaces, NoPadding)
				self.SetValue(request, "quoteOrderQty", x)
				params = self.Omit(params, "quoteOrderQty")
			} else if self.ToBool(!self.TestNil(price)) {
				x, _ := DecimalToPrecision(amount*price, Truncate, precision.(int), DecimalPlaces, NoPadding)
				self.SetValue(request, "quoteOrderQty", x)
			} else {
				quantityIsRequired = true
			}
		} else {
			quantityIsRequired = true
		}
	} else if self.ToBool(uppercaseType == "LIMIT") {
		priceIsRequired = true
		timeInForceIsRequired = true
		quantityIsRequired = true
	} else if self.ToBool(uppercaseType == "STOP_LOSS" || uppercaseType == "TAKE_PROFIT") {
		stopPriceIsRequired = true
		quantityIsRequired = true
		if self.ToBool(self.Member(market, "future")) {
			priceIsRequired = true
		}
	} else if self.ToBool(uppercaseType == "STOP_LOSS_LIMIT" || uppercaseType == "TAKE_PROFIT_LIMIT") {
		quantityIsRequired = true
		stopPriceIsRequired = true
		priceIsRequired = true
		timeInForceIsRequired = true
	} else if self.ToBool(uppercaseType == "LIMIT_MAKER") {
		priceIsRequired = true
		quantityIsRequired = true
	} else if self.ToBool(uppercaseType == "STOP") {
		quantityIsRequired = true
		stopPriceIsRequired = true
		priceIsRequired = true
	} else if self.ToBool(uppercaseType == "STOP_MARKET" || uppercaseType == "TAKE_PROFIT_MARKET") {
		closePosition := self.SafeValue(params, "closePosition", nil)
		if self.ToBool(self.TestNil(closePosition)) {
			quantityIsRequired = true
		}
		stopPriceIsRequired = true
	}
	if self.ToBool(quantityIsRequired) {
		self.SetValue(request, "quantity", self.AmountToPrecision(symbol, amount))
	}
	if self.ToBool(priceIsRequired) {
		if self.ToBool(self.TestNil(price)) {
			self.RaiseException("InvalidOrder", self.Id+" createOrder method requires a price argument for a "+typ+" order")
		}
		self.SetValue(request, "price", self.PriceToPrecision(symbol, price))
	}
	if self.ToBool(timeInForceIsRequired) {
		self.SetValue(request, "timeInForce", self.Member(self.Options, "defaultTimeInForce"))
	}
	if self.ToBool(stopPriceIsRequired) {
		stopPrice := self.SafeFloat(params, "stopPrice", 0)
		if self.ToBool(self.TestNil(stopPrice)) {
			self.RaiseException("InvalidOrder", self.Id+" createOrder method requires a stopPrice extra param for a "+typ+" order")
		} else {
			params = self.Omit(params, "stopPrice")
			self.SetValue(request, "stopPrice", self.PriceToPrecision(symbol, stopPrice))
		}
	}
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Binance) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if self.ToBool(self.TestNil(symbol)) {
		self.RaiseException("ArgumentsRequired", self.Id+" fetchOrder requires a symbol argument")
	}
	self.LoadMarkets()
	market := self.Market(symbol)
	defaultType := self.SafeString2(self.Options, "fetchOrder", "defaultType", market.Type)
	typ := self.SafeString(params, "type", defaultType)
	method := "privateGetOrder"
	if self.ToBool(typ == "future") {
		method = "fapiPrivateGetOrder"
	} else if self.ToBool(typ == "margin") {
		method = "sapiGetMarginOrder"
	}
	request := map[string]interface{}{
		"symbol": self.Member(market, "id"),
	}
	clientOrderId := self.SafeValue2(params, "origClientOrderId", "clientOrderId", "")
	if self.ToBool(!self.TestNil(clientOrderId)) {
		self.SetValue(request, "origClientOrderId", clientOrderId)
	} else {
		self.SetValue(request, "orderId", ToInteger(id))
	}
	query := self.Omit(params, []interface{}{"type", "clientOrderId", "origClientOrderId"})
	response := self.ApiFunc(method, self.Extend(request, query), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Binance) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	var market *Market
	var query interface{}
	var typ interface{}
	request := map[string]interface{}{}
	if self.ToBool(!self.TestNil(symbol)) {
		market = self.Market(symbol)
		self.SetValue(request, "symbol", self.Member(market, "id"))
		defaultType := self.SafeString2(self.Options, "fetchOpenOrders", "defaultType", market.Type)
		typ = self.SafeString(params, "type", defaultType)
		query = self.Omit(params, "type")
	} else if self.ToBool(self.Member(self.Options, "warnOnFetchOpenOrdersWithoutSymbol")) {
		symbols := self.Symbols
		numSymbols := self.Length(symbols)
		fetchOpenOrdersRateLimit := ToInteger(numSymbols / 2)
		self.RaiseException("ExchangeError", self.Id+" fetchOpenOrders WARNING: fetching open orders without specifying a symbol is rate-limited to one call per "+fmt.Sprintf("%v", fetchOpenOrdersRateLimit)+" seconds. Do not call this method frequently to avoid ban. Set "+self.Id+".options[warnOnFetchOpenOrdersWithoutSymbol] = false to suppress this warning message.")
	} else {
		defaultType := self.SafeString2(self.Options, "fetchOpenOrders", "defaultType", "spot")
		typ = self.SafeString(params, "type", defaultType)
		query = self.Omit(params, "type")
	}
	method := "privateGetOpenOrders"
	if self.ToBool(typ == "future") {
		method = "fapiPrivateGetOpenOrders"
	} else if self.ToBool(typ == "margin") {
		method = "sapiGetMarginOpenOrders"
	}
	response := self.ApiFuncReturnList(method, self.Extend(request, query), nil, nil)
	return self.ToOrders(self.ParseOrders(response, market, since, limit)), nil
}

func (self *Binance) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if self.ToBool(self.TestNil(symbol)) {
		self.RaiseException("ArgumentsRequired", self.Id+" cancelOrder requires a symbol argument")
	}
	self.LoadMarkets()
	market := self.Market(symbol)
	defaultType := self.SafeString2(self.Options, "fetchOpenOrders", "defaultType", market.Type)
	typ := self.SafeString(params, "type", defaultType)
	origClientOrderId := self.SafeValue2(params, "origClientOrderId", "clientOrderId", "")
	request := map[string]interface{}{
		"symbol": self.Member(market, "id"),
	}
	if self.ToBool(self.TestNil(origClientOrderId)) {
		self.SetValue(request, "orderId", ToInteger(id))
	} else {
		self.SetValue(request, "origClientOrderId", origClientOrderId)
	}
	method := "privateDeleteOrder"
	if self.ToBool(typ == "future") {
		method = "fapiPrivateDeleteOrder"
	} else if self.ToBool(typ == "margin") {
		method = "sapiDeleteMarginOrder"
	}
	query := self.Omit(params, []interface{}{"type", "origClientOrderId", "clientOrderId"})
	response = self.ApiFunc(method, self.Extend(request, query), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Binance) FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) (trades []*Trade, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"symbol": market.Id,
		"limit":  1000, // 默认 500
	}
	if limit > 0 {
		request["limit"] = limit
	}
	if since > 0 {
		request["startTime"] = since
	}
	response := self.ApiFuncReturnList("publicGetAggTrades", self.Extend(request, params), nil, nil)
	trades = self.ParseTrades(response, market, since, limit)
	return
}

func (self *Binance) ParseTrade(trade interface{}, market *Market) (result *Trade) {
	result = &Trade{
		Id:        fmt.Sprint(self.SafeInteger(trade, "a")),
		Timestamp: self.SafeInteger(trade, "T"),
		Price:     self.SafeFloat(trade, "p"),
		Amount:    self.SafeFloat(trade, "q"),
		Info:      trade,
	}
	result.Datetime = self.Iso8601(result.Timestamp)
	if self.SafeBool(trade, "m") {
		result.Side = "sell"
	} else {
		result.Side = "buy"
	}
	if market != nil {
		result.Symbol = market.Symbol
	}
	return
}

func (self *Binance) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	if self.ToBool(!self.ToBool(self.InMap(api, self.Member(self.Urls, "api")))) {
		self.RaiseException("NotSupported", self.Id+" does not have a testnet/sandbox URL for "+api+" endpoints")
	}
	url := self.Member(self.Member(self.Urls, "api"), api).(string)
	url += "/" + path
	if self.ToBool(api == "wapi") {
		url += ".html"
	}
	userDataStream := path == "userDataStream" || path == "listenKey"
	if self.ToBool(path == "historicalTrades") {
		if self.ToBool(self.ApiKey) {
			headers = map[string]interface{}{
				"X-MBX-APIKEY": self.ApiKey,
			}
		} else {
			self.RaiseException("AuthenticationError", self.Id+" historicalTrades endpoint requires `apiKey` credential")
		}
	} else if self.ToBool(userDataStream) {
		if self.ToBool(self.ApiKey) {
			body = self.Urlencode(params)
			headers = map[string]interface{}{
				"X-MBX-APIKEY": self.ApiKey,
				"Content-Type": "application/x-www-form-urlencoded",
			}
		} else {
			self.RaiseException("AuthenticationError", self.Id+" userDataStream endpoint requires `apiKey` credential")
		}
	}
	if self.ToBool(api == "private" || api == "sapi" || api == "wapi" && path != "systemStatus" || api == "fapiPrivate") {
		self.CheckRequiredCredentials()
		var query string
		if self.ToBool(api == "sapi" && path == "asset/dust") {
			query = self.UrlencodeWithArrayRepeat(self.Extend(map[string]interface{}{
				"timestamp":  self.Nonce(),
				"recvWindow": self.Member(self.Options, "recvWindow"),
			}, params))
		} else {
			/*else if self.ToBool(path == "batchOrders") {
				query = self.Rawencode(self.Extend(map[string]interface{}{
					"timestamp":  self.Nonce(),
					"recvWindow": self.Member(self.Options, "recvWindow"),
				}, params))
			}
			*/
			query = self.Urlencode(self.Extend(map[string]interface{}{
				"timestamp":  self.Nonce(),
				"recvWindow": self.Member(self.Options, "recvWindow"),
			}, params))
		}
		signature := self.Hmac(self.Encode(query), self.Encode(self.Secret), "sha256", "hex")
		query += "&" + "signature=" + signature
		headers = map[string]interface{}{
			"X-MBX-APIKEY": self.ApiKey,
		}
		if self.ToBool(method == "GET" || method == "DELETE" || api == "wapi") {
			url += "?" + query
		} else {
			body = query
			self.SetValue(headers, "Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		if self.ToBool(!self.ToBool(userDataStream)) {
			if self.ToBool(self.Length(reflect.ValueOf(params).MapKeys())) {
				url += "?" + self.Urlencode(params)
			}
		}
	}
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Binance) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	if self.ToBool(httpCode == 418 || httpCode == 429) {
		self.RaiseException("DDoSProtection", self.Id+" "+fmt.Sprintf("%v", httpCode)+" "+reason+" "+body)
	}
	if self.ToBool(httpCode >= 400) {
		if strings.Contains(body, "Price * QTY is zero or less") {
			self.RaiseException("InvalidOrder", self.Id+" order cost = amount * price is zero or less "+body)
		}
		if strings.Contains(body, "LOT_SIZE") {
			self.RaiseException("InvalidOrder", self.Id+" order amount should be evenly divisible by lot size "+body)
		}
		if strings.Contains(body, "PRICE_FILTER") {
			self.RaiseException("InvalidOrder", self.Id+" order price is invalid, i.e. exceeds allowed price precision, exceeds min price or max price limits or is invalid float value in general, use this.priceToPrecision (symbol, amount) "+body)
		}
	}
	if self.ToBool(self.TestNil(response)) {
		return
	}
	success := self.SafeValue(response, "success", true)
	if !success.(bool) {
		message := self.SafeString(response, "msg", "")
		var parsedMessage map[string]interface{}
		if message != "" {
			if err := json.Unmarshal([]byte(message), &parsedMessage); err != nil {
				response = parsedMessage
			}
		}
	}
	message := self.SafeString(response, "msg", "")
	if self.ToBool(!self.TestNil(message)) {
		self.ThrowExactlyMatchedException(self.Exceptions, message, self.Id+" "+message)
	}
	errorStr := self.SafeString(response, "code", "")
	if errorStr != "" {
		if self.ToBool(errorStr == "200") {
			return
		}
		if errorStr == "-2015" && self.Options["hasAlreadyAuthenticatedSuccessfully"].(bool) {
			self.RaiseException("DDoSProtection", self.Id+" temporary banned: "+body)
		}
		feedback := self.Id + " " + body
		self.ThrowExactlyMatchedException(self.Exceptions, errorStr, feedback)
		self.RaiseException("ExchangeError", feedback)
	}
	if self.ToBool(!self.ToBool(success)) {
		self.RaiseException("ExchangeError", self.Id+" "+body)
	}
}
