package futures_gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	. "github.com/epheien/ccxt/go/base"
	"math"
	urllib "net/url"
	"strconv"
	"strings"
)

type FuturesGateio struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *FuturesGateio, err error) {
	ex = new(FuturesGateio)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *FuturesGateio) Describe() []byte {
	return []byte(`
{
    "id": "futures_gateio",
    "name": "Gate.io",
    "countries": "JP",
    "rateLimit": 500,
	"version": "v4",
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
        "logo": "",
        "api": {
            "public": "https://api.gateio.ws/api/v4",
            "private": "https://api.gateio.ws/api/v4",
        },
        "www": "https://www.gate.io",
        "doc": "https://www.gate.io/docs/developers/apiv4/en/#futures",
        "fees": [
        ]
    },
    "api": {
        "public": {
            "get": [
				"futures/usdt/contracts",
				"futures/usdt/contracts/{contract}",
				"futures/usdt/order_book",
				"futures/usdt/trades",
            ]
        },
        "private": {
            "get": [
				"futures/usdt/accounts",
				"futures/usdt/orders/{order_id}",
				"futures/usdt/orders",
				"futures/usdt/positions",
				"futures/usdt/positions/{contract}",
            ],
            "post": [
				"futures/usdt/orders",
            ],
            "delete": [
				"futures/usdt/orders/{order_id}",
            ],
            "put": [
				"futures/usdt/orders/{order_id}",
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

func (self *FuturesGateio) LoadMarkets() map[string]*Market {
	return nil
}

func (self *FuturesGateio) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	response := self.ApiFuncReturnList("publicGetFuturesUsdtContracts", params, nil, nil)
	data := response
	result := []interface{}{}
	for i := 0; i < self.Length(data); i++ {
		market := self.Member(data, i)
		id := self.SafeString(market, "name", "")
		baseId, quoteId := self.Unpack2(strings.Split(id, "_"))
		base := self.SafeCurrencyCode(baseId)
		quote := self.SafeCurrencyCode(quoteId)
		symbol := base + "/" + quote
		active := !self.SafeBool(market, "in_delisting")
		multiplier := self.SafeFloat(market, "quanto_multiplier")
		precision := map[string]interface{}{
			"base":   self.PrecisionFromString(self.SafeString(market, "quanto_multiplier")),
			"quote":  self.PrecisionFromString(self.SafeString(market, "order_price_round")),
			"amount": self.PrecisionFromString(self.SafeString(market, "quanto_multiplier")),
			"price":  self.PrecisionFromString(self.SafeString(market, "order_price_round")),
		}
		limits := map[string]interface{}{
			"amount": map[string]interface{}{
				"min": self.SafeFloat(market, "order_size_min") * multiplier,
				"max": self.SafeFloat(market, "order_size_max") * multiplier,
			},
		}
		result = append(result, map[string]interface{}{
			"futures":        true,
			"id":             id,
			"symbol":         symbol,
			"baseId":         baseId,
			"quoteId":        quoteId,
			"base":           base,
			"quote":          quote,
			"active":         active,
			"precision":      precision,
			"limits":         limits,
			"baseMultiplier": self.SafeFloat(market, "quanto_multiplier"),
			"info":           market,
		})
	}
	return self.ToMarkets(result), nil
}

// BTC/USDT => BTC_USDT
func (self *FuturesGateio) Market(symbol string) *Market {
	li := strings.Split(symbol, "/")
	if len(li) != 2 {
		return &Market{}
	}
	return &Market{
		Id:     li[0] + "_" + li[1],
		Symbol: symbol,
		Base:   li[0],
		Quote:  li[1],
	}
}

func (self *FuturesGateio) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()

	response := self.ApiFunc("privateGetFuturesUsdtAccounts", params, nil, nil)

	result := map[string]interface{}{
		"info": response,
	}

	balance := response
	account := self.Account()
	// 钱包余额是用户累计充值提现和盈亏结算(包括已实现盈亏, 资金费用,手续费及推荐返佣)之后的余额,
	// 不包含未实现盈亏. total = SUM(history_dnw, history_pnl, history_fee, history_refr, history_fund)
	//total := self.SafeFloat(balance, "total")
	free := self.SafeFloat(balance, "available") // NOTE: 不包含未实现盈亏
	used := self.SafeFloat(balance, "order_margin") + self.SafeFloat(balance, "position_margin")
	unrealPnl := self.SafeFloat(balance, "unrealised_pnl")
	account["free"] = free + unrealPnl // 我们的 free 包含未实验盈亏
	account["used"] = used
	account["total"] = free + used
	account["unrealPnl"] = unrealPnl
	currency := self.SafeString(balance, "currency")
	result[currency] = account

	// point
	account = self.Account()
	account["free"] = self.SafeFloat(balance, "point")
	result["point"] = account

	return self.ParseBalance(result), nil
}

func (self *FuturesGateio) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"contract": market.Id,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	response := self.ApiFunc("publicGetFuturesUsdtOrderBook", self.Extend(request, params), nil, nil)
	orderBook = new(OrderBook)
	orderBook.Timestamp = int64(ToFloat(response["update"]) * 1000)
	orderBook.Datetime = self.Iso8601(orderBook.Timestamp)
	for _, one := range response["bids"].([]interface{}) {
		bid := one.(map[string]interface{})
		price := ToFloat(bid["p"])
		amount := ToFloat(bid["s"])
		orderBook.Bids = append(orderBook.Bids, [2]float64{price, amount})
	}
	for _, one := range response["asks"].([]interface{}) {
		ask := one.(map[string]interface{})
		price := ToFloat(ask["p"])
		amount := ToFloat(ask["s"])
		orderBook.Asks = append(orderBook.Asks, [2]float64{price, amount})
	}
	return orderBook, nil
}

func (self *FuturesGateio) ParseTicker(response interface{}) (ticker *Ticker) {
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

func (self *FuturesGateio) FetchTicker(symbol string, params map[string]interface{}) (ticker *Ticker, err error) {
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

func (self *FuturesGateio) ParseOHLCV(response interface{}) *OHLCV {
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

func (self *FuturesGateio) FetchOHLCV(symbol, timeframe string, since int64, limit int64, params map[string]interface{}) (klines []*OHLCV, err error) {
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

func (self *FuturesGateio) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"open":     "open",
		"finished": "closed",
	}
	return self.SafeString(statuses, status, status)
}

func (self *FuturesGateio) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	symbol := ""
	if market != nil {
		symbol = market.(*Market).Symbol
	}
	clientOid := ""
	orderId := fmt.Sprintf("%d", self.SafeInteger(order, "id"))
	timestamp := int64(self.SafeFloat(order, "create_time", 0) * 1000)
	datetime := self.Iso8601(timestamp)
	price := self.SafeFloat(order, "price", 0)
	amount := self.SafeFloat(order, "size", 0)
	side := "buy"
	if amount < 0 {
		side = "sell"
	}
	amount = math.Abs(amount)
	remaining := math.Abs(self.SafeFloat(order, "left", 0))
	status := self.ParseOrderStatus(self.SafeString(order, "status"))
	average := self.SafeFloat(order, "fill_price")
	return map[string]interface{}{
		"clientOrderId": clientOid,
		"id":            orderId,
		"symbol":        symbol,
		"type":          "",
		"side":          side,
		"amount":        amount,
		"price":         price,
		"filled":        amount - remaining,
		"remaining":     remaining,
		"average":       average,
		"timestamp":     timestamp,
		"datetime":      datetime,
		"status":        status,
		"info":          order,
	}
}

func (self *FuturesGateio) CreateOrder(symbol string, type_ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"contract": market.Id,
		"price":    self.Float64ToString(price),
	}
	if side == "buy" {
		request["size"] = int64(amount)
	} else {
		request["size"] = -int64(amount)
	}
	response := self.ApiFunc("privatePostFuturesUsdtOrders", self.Extend(request, params), nil, nil)
	return &Order{
		Id:     fmt.Sprintf("%d", self.SafeInteger(response, "id")),
		Symbol: symbol,
		Type:   type_,
		Side:   side,
		Status: "open",
		Info:   response,
	}, nil
}

func (self *FuturesGateio) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"order_id": id,
	}
	response := self.ApiFunc("privateGetFuturesUsdtOrdersOrderId", self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *FuturesGateio) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{
		"status": "open",
	}
	var market *Market
	if symbol != "" {
		market = self.Market(symbol)
		request["contract"] = market.Id
	}
	response := self.ApiFuncReturnList("privateGetFuturesUsdtOrders", self.Extend(request, params), nil, nil)
	return self.ToOrders(self.ParseOrders(response, market, since, limit)), nil
}

func (self *FuturesGateio) FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) (trades []*Trade, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"contract": market.Id,
	}
	if limit > 0 {
		request["limit"] = limit
	}
	if since > 0 {
		request["from"] = since
	}
	response := self.ApiFuncReturnList("publicGetFuturesUsdtTrades", self.Extend(request, params), nil, nil)
	trades = self.ParseTrades(response, market, since, limit)
	trades = self.ReverseTrades(trades)
	return
}

func (self *FuturesGateio) ParseTrade(trade interface{}, market *Market) (result *Trade) {
	result = &Trade{
		Id:        strconv.FormatInt(self.SafeInteger(trade, "id"), 10),
		Timestamp: int64(self.SafeFloat(trade, "create_time") * 1000),
		Price:     self.SafeFloat(trade, "price"),
		Info:      trade,
	}
	amount := self.SafeFloat(trade, "size")
	result.Amount = math.Abs(amount)
	if amount >= 0 {
		result.Side = "buy"
	} else {
		result.Side = "sell"
	}
	result.Datetime = self.Iso8601(result.Timestamp)
	if market != nil {
		result.Symbol = market.Symbol
	}
	return
}

func (self *FuturesGateio) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{
		"order_id": id,
	}
	response = self.ApiFunc("privateDeleteFuturesUsdtOrdersOrderId", self.Extend(request, params), nil, nil)
	return response, nil
}

func (self *FuturesGateio) FetchMarkPrice(symbol string, params map[string]interface{}) (markPrice *MarkPrice, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"contract": market.Id,
	}
	response := self.ApiFunc("publicGetFuturesUsdtContractsContract", self.Extend(request, params), nil, nil)
	data := response
	return &MarkPrice{
		Symbol:     symbol,
		MarkPrice:  self.SafeFloat(data, "mark_price", 0),
		IndexPrice: self.SafeFloat(data, "index_price", 0),
		//Timestamp:  self.SafeInteger(data, "time", 0),
		Info: data,
	}, nil
}

func (self *FuturesGateio) FetchPositions(symbol string, params map[string]interface{}) (result []*Position, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	request := map[string]interface{}{}
	market := self.Market(symbol)
	request["contract"] = market.Id
	response := self.ApiFunc("privateGetFuturesUsdtPositionsContract", self.Extend(request, params), nil, nil)

	item := response
	amount := self.SafeFloat(item, "size")
	pos := &Position{
		Symbol:     symbol,
		Side:       "long",
		Leverage:   self.SafeFloat(item, "cross_leverage_limit", 0),
		Amount:     math.Abs(amount),
		UsedAmount: 0,
		Price:      self.SafeFloat(item, "entry_price", 0),
		RealPnl:    self.SafeFloat(item, "realised_pnl"),
		UnrealPnl:  self.SafeFloat(item, "unrealised_pnl", 0),
		Info:       item,
	}
	if amount < 0 {
		pos.Side = "short"
	}
	result = append(result, pos)

	return
}

func (self *FuturesGateio) genSign(method, url, query, body string) map[string]interface{} {
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

func (self *FuturesGateio) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
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

func (self *FuturesGateio) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
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
