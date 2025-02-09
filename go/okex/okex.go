package okex

import (
	"fmt"
	. "github.com/epheien/ccxt/go/base"
	"math"
	"reflect"
	"strings"
)

type Okex struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Okex, err error) {
	ex = new(Okex)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *Okex) Describe() []byte {
	return []byte(`{
    "id": "okex",
    "name": "OKEX",
    "countries": [
        "CN",
        "US"
    ],
    "version": "v3",
    "rateLimit": 1000,
    "pro": true,
    "has": {
        "CORS": false,
        "fetchOHLCV": true,
        "fetchOrder": true,
        "fetchOrders": false,
        "fetchOpenOrders": true,
        "fetchClosedOrders": true,
        "fetchCurrencies": false,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchTime": true,
        "fetchTransactions": false,
        "fetchMyTrades": true,
        "fetchDepositAddress": true,
        "fetchOrderTrades": true,
        "fetchTickers": true,
        "fetchLedger": true,
        "withdraw": true,
        "futures": true
    },
    "timeframes": {
        "1m": "60",
        "3m": "180",
        "5m": "300",
        "15m": "900",
        "30m": "1800",
        "1h": "3600",
        "2h": "7200",
        "4h": "14400",
        "6h": "21600",
        "12h": "43200",
        "1d": "86400",
        "1w": "604800"
    },
    "hostname": "okex.com",
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/32552768-0d6dd3c6-c4a6-11e7-90f8-c043b64756a7.jpg",
        "api": {
            "rest": "https://www.{hostname}"
        },
        "www": "https://www.okex.com",
        "doc": "https://www.okex.com/docs/en/",
        "fees": "https://www.okex.com/pages/products/fees.html",
        "referral": "https://www.okex.com/join/1888677",
        "test": {
            "rest": "https://testnet.okex.com"
        }
    },
    "api": {
        "general": {
            "get": [
                "time"
            ]
        },
        "account": {
            "get": [
                "wallet",
                "sub-account",
                "asset-valuation",
                "wallet/{currency}",
                "withdrawal/history",
                "withdrawal/history/{currency}",
                "ledger",
                "deposit/address",
                "deposit/history",
                "deposit/history/{currency}",
                "currencies",
                "withdrawal/fee"
            ],
            "post": [
                "transfer",
                "withdrawal"
            ]
        },
        "spot": {
            "get": [
                "accounts",
                "accounts/{currency}",
                "accounts/{currency}/ledger",
                "orders",
                "orders_pending",
                "orders/{order_id}",
                "orders/{client_oid}",
                "trade_fee",
                "fills",
                "algo",
                "instruments",
                "instruments/{instrument_id}/book",
                "instruments/ticker",
                "instruments/{instrument_id}/ticker",
                "instruments/{instrument_id}/trades",
                "instruments/{instrument_id}/candles"
            ],
            "post": [
                "order_algo",
                "orders",
                "batch_orders",
                "cancel_orders/{order_id}",
                "cancel_orders/{client_oid}",
                "cancel_batch_algos",
                "cancel_batch_orders"
            ]
        },
        "margin": {
            "get": [
                "accounts",
                "accounts/{instrument_id}",
                "accounts/{instrument_id}/ledger",
                "accounts/availability",
                "accounts/{instrument_id}/availability",
                "accounts/borrowed",
                "accounts/{instrument_id}/borrowed",
                "orders",
                "accounts/{instrument_id}/leverage",
                "orders/{order_id}",
                "orders/{client_oid}",
                "orders_pending",
                "fills",
                "instruments/{instrument_id}/mark_price"
            ],
            "post": [
                "accounts/borrow",
                "accounts/repayment",
                "orders",
                "batch_orders",
                "cancel_orders",
                "cancel_orders/{order_id}",
                "cancel_orders/{client_oid}",
                "cancel_batch_orders",
                "accounts/{instrument_id}/leverage"
            ]
        },
        "futures": {
            "get": [
                "position",
                "{instrument_id}/position",
                "accounts",
                "accounts/{underlying}",
                "accounts/{underlying}/leverage",
                "accounts/{underlying}/ledger",
                "order_algo/{instrument_id}",
                "orders/{instrument_id}",
                "orders/{instrument_id}/{order_id}",
                "orders/{instrument_id}/{client_oid}",
                "fills",
                "trade_fee",
                "accounts/{instrument_id}/holds",
                "order_algo/{instrument_id}",
                "instruments",
                "instruments/{instrument_id}/book",
                "instruments/ticker",
                "instruments/{instrument_id}/ticker",
                "instruments/{instrument_id}/trades",
                "instruments/{instrument_id}/candles",
                "instruments/{instrument_id}/index",
                "rate",
                "instruments/{instrument_id}/estimated_price",
                "instruments/{instrument_id}/open_interest",
                "instruments/{instrument_id}/price_limit",
                "instruments/{instrument_id}/mark_price",
                "instruments/{instrument_id}/liquidation"
            ],
            "post": [
                "accounts/{underlying}/leverage",
                "order",
                "orders",
                "cancel_order/{instrument_id}/{order_id}",
                "cancel_order/{instrument_id}/{client_oid}",
                "cancel_batch_orders/{instrument_id}",
                "accounts/margin_mode",
                "close_position",
                "cancel_all",
                "order_algo",
                "cancel_algos"
            ]
        },
        "swap": {
            "get": [
                "position",
                "{instrument_id}/position",
                "accounts",
                "{instrument_id}/accounts",
                "accounts/{instrument_id}/settings",
                "accounts/{instrument_id}/ledger",
                "orders/{instrument_id}",
                "orders/{instrument_id}/{order_id}",
                "orders/{instrument_id}/{client_oid}",
                "fills",
                "accounts/{instrument_id}/holds",
                "trade_fee",
                "order_algo/{instrument_id}",
                "instruments",
                "instruments/{instrument_id}/depth",
                "instruments/ticker",
                "instruments/{instrument_id}/ticker",
                "instruments/{instrument_id}/trades",
                "instruments/{instrument_id}/candles",
                "instruments/{instrument_id}/index",
                "rate",
                "instruments/{instrument_id}/open_interest",
                "instruments/{instrument_id}/price_limit",
                "instruments/{instrument_id}/liquidation",
                "instruments/{instrument_id}/funding_time",
                "instruments/{instrument_id}/mark_price",
                "instruments/{instrument_id}/historical_funding_rate"
            ],
            "post": [
                "accounts/{instrument_id}/leverage",
                "order",
                "orders",
                "cancel_order/{instrument_id}/{order_id}",
                "cancel_order/{instrument_id}/{client_oid}",
                "cancel_batch_orders/{instrument_id}",
                "order_algo",
                "cancel_algos"
            ]
        },
        "option": {
            "get": [
                "accounts",
                "{underlying}/position",
                "accounts/{underlying}",
                "orders/{underlying}",
                "fills/{underlying}",
                "accounts/{underlying}/ledger",
                "trade_fee",
                "orders/{underlying}/{order_id}",
                "orders/{underlying}/{client_oid}",
                "underlying",
                "instruments/{underlying}",
                "instruments/{underlying}/summary",
                "instruments/{underlying}/summary/{instrument_id}",
                "instruments/{instrument_id}/book",
                "instruments/{instrument_id}/trades",
                "instruments/{instrument_id}/ticker",
                "instruments/{instrument_id}/candles"
            ],
            "post": [
                "order",
                "orders",
                "cancel_order/{underlying}/{order_id}",
                "cancel_order/{underlying}/{client_oid}",
                "cancel_batch_orders/{underlying}",
                "amend_order/{underlying}",
                "amend_batch_orders/{underlying}"
            ]
        },
        "index": {
            "get": [
                "{instrument_id}/constituents"
            ]
        }
    },
    "fees": {
        "trading": {
            "taker": 0.0015,
            "maker": 0.001
        },
        "spot": {
            "taker": 0.0015,
            "maker": 0.001
        },
        "futures": {
            "taker": 0.0005,
            "maker": 0.0002
        },
        "swap": {
            "taker": 0.00075,
            "maker": 0.0002
        }
    },
    "requiredCredentials": {
        "apiKey": true,
        "secret": true,
        "password": true
    },
    "exceptions": {
        "exact": {
            "1": "ExchangeError",
            "failure to get a peer from the ring-balancer": "ExchangeNotAvailable",
            "4010": "PermissionDenied",
            "4001": "ExchangeError",
            "4002": "ExchangeError",
            "30001": "AuthenticationError",
            "30002": "AuthenticationError",
            "30003": "AuthenticationError",
            "30004": "AuthenticationError",
            "30005": "InvalidNonce",
            "30006": "AuthenticationError",
            "30007": "BadRequest",
            "30008": "RequestTimeout",
            "30009": "ExchangeError",
            "30010": "AuthenticationError",
            "30011": "PermissionDenied",
            "30012": "AuthenticationError",
            "30013": "AuthenticationError",
            "30014": "DDoSProtection",
            "30015": "AuthenticationError",
            "30016": "ExchangeError",
            "30017": "ExchangeError",
            "30018": "ExchangeError",
            "30019": "ExchangeNotAvailable",
            "30020": "BadRequest",
            "30021": "BadRequest",
            "30022": "PermissionDenied",
            "30023": "BadRequest",
            "30024": "BadSymbol",
            "30025": "BadRequest",
            "30026": "DDoSProtection",
            "30027": "AuthenticationError",
            "30028": "PermissionDenied",
            "30029": "AccountSuspended",
            "30030": "ExchangeError",
            "30031": "BadRequest",
            "30032": "BadSymbol",
            "30033": "BadRequest",
            "30034": "ExchangeError",
            "30035": "ExchangeError",
            "30036": "ExchangeError",
            "30037": "ExchangeNotAvailable",
            "30038": "OnMaintenance",
            "32001": "AccountSuspended",
            "32002": "PermissionDenied",
            "32003": "CancelPending",
            "32004": "ExchangeError",
            "32005": "InvalidOrder",
            "32006": "InvalidOrder",
            "32007": "InvalidOrder",
            "32008": "InvalidOrder",
            "32009": "InvalidOrder",
            "32010": "ExchangeError",
            "32011": "ExchangeError",
            "32012": "ExchangeError",
            "32013": "ExchangeError",
            "32014": "ExchangeError",
            "32015": "ExchangeError",
            "32016": "ExchangeError",
            "32017": "ExchangeError",
            "32018": "ExchangeError",
            "32019": "ExchangeError",
            "32020": "ExchangeError",
            "32021": "ExchangeError",
            "32022": "ExchangeError",
            "32023": "ExchangeError",
            "32024": "ExchangeError",
            "32025": "ExchangeError",
            "32026": "ExchangeError",
            "32027": "ExchangeError",
            "32028": "ExchangeError",
            "32029": "ExchangeError",
            "32030": "InvalidOrder",
            "32031": "ArgumentsRequired",
            "32038": "AuthenticationError",
            "32040": "ExchangeError",
            "32044": "ExchangeError",
            "32045": "ExchangeError",
            "32046": "ExchangeError",
            "32047": "ExchangeError",
            "32048": "InvalidOrder",
            "32049": "ExchangeError",
            "32050": "InvalidOrder",
            "32051": "InvalidOrder",
            "32052": "ExchangeError",
            "32053": "ExchangeError",
            "32057": "ExchangeError",
            "32054": "ExchangeError",
            "32055": "InvalidOrder",
            "32056": "ExchangeError",
            "32058": "ExchangeError",
            "32059": "InvalidOrder",
            "32060": "InvalidOrder",
            "32061": "InvalidOrder",
            "32062": "InvalidOrder",
            "32063": "InvalidOrder",
            "32064": "ExchangeError",
            "32065": "ExchangeError",
            "32066": "ExchangeError",
            "32067": "ExchangeError",
            "32068": "ExchangeError",
            "32069": "ExchangeError",
            "32070": "ExchangeError",
            "32071": "ExchangeError",
            "32072": "ExchangeError",
            "32073": "ExchangeError",
            "32074": "ExchangeError",
            "32075": "ExchangeError",
            "32076": "ExchangeError",
            "32077": "ExchangeError",
            "32078": "ExchangeError",
            "32079": "ExchangeError",
            "32080": "ExchangeError",
            "32083": "ExchangeError",
            "33001": "PermissionDenied",
            "33002": "AccountSuspended",
            "33003": "InsufficientFunds",
            "33004": "ExchangeError",
            "33005": "ExchangeError",
            "33006": "ExchangeError",
            "33007": "ExchangeError",
            "33008": "InsufficientFunds",
            "33009": "ExchangeError",
            "33010": "ExchangeError",
            "33011": "ExchangeError",
            "33012": "ExchangeError",
            "33013": "InvalidOrder",
            "33014": "OrderNotFound",
            "33015": "InvalidOrder",
            "33016": "ExchangeError",
            "33017": "InsufficientFunds",
            "33018": "ExchangeError",
            "33020": "ExchangeError",
            "33021": "BadRequest",
            "33022": "InvalidOrder",
            "33023": "ExchangeError",
            "33024": "InvalidOrder",
            "33025": "InvalidOrder",
            "33026": "ExchangeError",
            "33027": "InvalidOrder",
            "33028": "InvalidOrder",
            "33029": "InvalidOrder",
            "33034": "ExchangeError",
            "33035": "ExchangeError",
            "33036": "ExchangeError",
            "33037": "ExchangeError",
            "33038": "ExchangeError",
            "33039": "ExchangeError",
            "33040": "ExchangeError",
            "33041": "ExchangeError",
            "33042": "ExchangeError",
            "33043": "ExchangeError",
            "33044": "ExchangeError",
            "33045": "ExchangeError",
            "33046": "ExchangeError",
            "33047": "ExchangeError",
            "33048": "ExchangeError",
            "33049": "ExchangeError",
            "33050": "ExchangeError",
            "33051": "ExchangeError",
            "33059": "BadRequest",
            "33060": "BadRequest",
            "33061": "ExchangeError",
            "33062": "ExchangeError",
            "33063": "ExchangeError",
            "33064": "ExchangeError",
            "33065": "ExchangeError",
            "21009": "ExchangeError",
            "34001": "PermissionDenied",
            "34002": "InvalidAddress",
            "34003": "ExchangeError",
            "34004": "ExchangeError",
            "34005": "ExchangeError",
            "34006": "ExchangeError",
            "34007": "ExchangeError",
            "34008": "InsufficientFunds",
            "34009": "ExchangeError",
            "34010": "ExchangeError",
            "34011": "ExchangeError",
            "34012": "ExchangeError",
            "34013": "ExchangeError",
            "34014": "ExchangeError",
            "34015": "ExchangeError",
            "34016": "PermissionDenied",
            "34017": "AccountSuspended",
            "34018": "AuthenticationError",
            "34019": "PermissionDenied",
            "34020": "PermissionDenied",
            "34021": "InvalidAddress",
            "34022": "ExchangeError",
            "34023": "PermissionDenied",
            "34026": "ExchangeError",
            "34036": "ExchangeError",
            "34037": "ExchangeError",
            "34038": "ExchangeError",
            "34039": "ExchangeError",
            "35001": "ExchangeError",
            "35002": "ExchangeError",
            "35003": "ExchangeError",
            "35004": "ExchangeError",
            "35005": "AuthenticationError",
            "35008": "InvalidOrder",
            "35010": "InvalidOrder",
            "35012": "InvalidOrder",
            "35014": "InvalidOrder",
            "35015": "InvalidOrder",
            "35017": "ExchangeError",
            "35019": "InvalidOrder",
            "35020": "InvalidOrder",
            "35021": "InvalidOrder",
            "35022": "ExchangeError",
            "35024": "ExchangeError",
            "35025": "InsufficientFunds",
            "35026": "ExchangeError",
            "35029": "OrderNotFound",
            "35030": "InvalidOrder",
            "35031": "InvalidOrder",
            "35032": "ExchangeError",
            "35037": "ExchangeError",
            "35039": "ExchangeError",
            "35040": "InvalidOrder",
            "35044": "ExchangeError",
            "35046": "InsufficientFunds",
            "35047": "InsufficientFunds",
            "35048": "ExchangeError",
            "35049": "InvalidOrder",
            "35050": "InvalidOrder",
            "35052": "InsufficientFunds",
            "35053": "ExchangeError",
            "35055": "InsufficientFunds",
            "35057": "ExchangeError",
            "35058": "ExchangeError",
            "35059": "BadRequest",
            "35060": "BadRequest",
            "35061": "BadRequest",
            "35062": "InvalidOrder",
            "35063": "InvalidOrder",
            "35064": "InvalidOrder",
            "35066": "InvalidOrder",
            "35067": "InvalidOrder",
            "35068": "InvalidOrder",
            "35069": "InvalidOrder",
            "35070": "InvalidOrder",
            "35071": "InvalidOrder",
            "35072": "InvalidOrder",
            "35073": "InvalidOrder",
            "35074": "InvalidOrder",
            "35075": "InvalidOrder",
            "35076": "InvalidOrder",
            "35077": "InvalidOrder",
            "35078": "InvalidOrder",
            "35079": "InvalidOrder",
            "35080": "InvalidOrder",
            "35081": "InvalidOrder",
            "35082": "InvalidOrder",
            "35083": "InvalidOrder",
            "35084": "InvalidOrder",
            "35085": "InvalidOrder",
            "35086": "InvalidOrder",
            "35087": "InvalidOrder",
            "35088": "InvalidOrder",
            "35089": "InvalidOrder",
            "35090": "ExchangeError",
            "35091": "ExchangeError",
            "35092": "ExchangeError",
            "35093": "ExchangeError",
            "35094": "ExchangeError",
            "35095": "BadRequest",
            "35096": "ExchangeError",
            "35097": "ExchangeError",
            "35098": "ExchangeError",
            "35099": "ExchangeError",
            "36001": "BadRequest",
            "36002": "BadRequest",
            "36005": "ExchangeError",
            "36101": "AuthenticationError",
            "36102": "PermissionDenied",
            "36103": "PermissionDenied",
            "36104": "PermissionDenied",
            "36105": "PermissionDenied",
            "36106": "PermissionDenied",
            "36107": "PermissionDenied",
            "36108": "InsufficientFunds",
            "36109": "PermissionDenied",
            "36201": "PermissionDenied",
            "36202": "PermissionDenied",
            "36203": "InvalidOrder",
            "36204": "ExchangeError",
            "36205": "BadRequest",
            "36206": "BadRequest",
            "36207": "InvalidOrder",
            "36208": "InvalidOrder",
            "36209": "InvalidOrder",
            "36210": "InvalidOrder",
            "36211": "InvalidOrder",
            "36212": "InvalidOrder",
            "36213": "InvalidOrder",
            "36214": "ExchangeError",
            "36216": "OrderNotFound",
            "36217": "InvalidOrder",
            "36218": "InvalidOrder",
            "36219": "InvalidOrder",
            "36220": "InvalidOrder",
            "36221": "InvalidOrder",
            "36222": "InvalidOrder",
            "36223": "InvalidOrder",
            "36224": "InvalidOrder",
            "36225": "InvalidOrder",
            "36226": "InvalidOrder",
            "36227": "InvalidOrder",
            "36228": "InvalidOrder",
            "36229": "InvalidOrder",
            "36230": "InvalidOrder"
        },
        "broad": {}
    },
    "precisionMode": "TICK_SIZE",
    "options": {
        "createMarketBuyOrderRequiresPrice": true,
        "fetchMarkets": [
            "spot",
            "futures",
            "swap",
            "option"
        ],
        "defaultType": "spot",
        "auth": {
            "time": "public",
            "currencies": "private",
            "instruments": "public",
            "rate": "public",
            "{instrument_id}/constituents": "public"
        }
    },
    "commonCurrencies": {
        "AE": "AET",
        "HOT": "Hydro Protocol",
        "HSR": "HC",
        "MAG": "Maggie",
        "YOYO": "YOYOW",
        "WIN": "WinToken"
    }
}`)
}

func (self *Okex) FetchMarkets(params map[string]interface{}) []interface{} {
	types := self.SafeValue(self.Options, "fetchMarkets", nil)
	result := []interface{}{}
	for i := 0; i < self.Length(types); i++ {
		typ := self.Member(types, i).(string)
		// TODO, pass option
		if typ == "option" {
			continue
		}
		markets := self.FetchMarketsByType(typ, params)
		result = self.ArrayConcat(result, markets)
	}
	return result
}

func (self *Okex) ParseMarkets(markets []interface{}) []interface{} {
	result := []interface{}{}
	for i := 0; i < self.Length(markets); i++ {
		result = append(result, self.ParseMarket(self.Member(markets, i)))
	}
	return result
}

func (self *Okex) ParseMarket(market interface{}) interface{} {
	id := self.SafeString(market, "instrument_id", "")
	marketType := "spot"
	spot := true
	future := false
	swap := false
	option := false
	baseId := self.SafeString(market, "base_currency", "")
	quoteId := self.SafeString(market, "quote_currency", "")
	contractVal := self.SafeFloat(market, "contract_val", 0)
	if self.ToBool(!self.TestNil(contractVal)) {
		if self.ToBool(self.InMap("option_type", market)) {
			marketType = "option"
			spot = false
			option = true
			underlying := self.SafeString(market, "underlying", "")
			parts := strings.Split(underlying, "-")
			baseId = ""
			quoteId = ""
			if len(parts) >= 2 {
				baseId = parts[0]
				quoteId = parts[1]
			}
		} else {
			marketType = "swap"
			spot = false
			swap = true
			futuresAlias := self.SafeString(market, "alias", "")
			if self.ToBool(!self.TestNil(futuresAlias)) {
				swap = false
				future = true
				marketType = "futures"
				baseId = self.SafeString(market, "underlying_index", "")
			}
		}
	}
	base := self.SafeCurrencyCode(baseId)
	quote := self.SafeCurrencyCode(quoteId)
	symbol := self.IfThenElse(self.ToBool(spot), base+"/"+quote, id)
	lotSize := self.SafeFloat2(market, "lot_size", "trade_increment", 0.0)

	amountPrecision := self.SafeFloat(market, "size_increment", lotSize)
	pricePrecision :=  self.SafeFloat(market, "tick_size", 0)
	precision := map[string]interface{}{
		//"amount": self.SafeFloat(market, "size_increment", lotSize),
		//"price":  self.SafeFloat(market, "tick_size", 0),
		"amount": self.PrecisionFromString(NumberToString(amountPrecision)),
		"price":  self.PrecisionFromString(NumberToString(pricePrecision)),
	}
	minAmount := self.SafeFloat2(market, "min_size", "base_min_size", 0.0)
	active := true
	fees := self.SafeValue2(self.Fees, marketType, "trading", map[string]interface{}{}).(map[string]interface{})
	return self.Extend(fees, map[string]interface{}{
		"id":        id,
		"symbol":    symbol,
		"base":      base,
		"quote":     quote,
		"baseId":    baseId,
		"quoteId":   quoteId,
		"info":      market,
		"type":      marketType,
		"spot":      spot,
		"futures":   future,
		"swap":      swap,
		"option":    option,
		"active":    active,
		"precision": precision,
		"limits": map[string]interface{}{
			"amount": map[string]interface{}{
				"min": minAmount,
				"max": nil,
			},
			"price": map[string]interface{}{
				"min": pricePrecision,
				"max": nil,
			},
			"cost": map[string]interface{}{
				"min": pricePrecision,
				"max": nil,
			},
		},
	})
}

func (self *Okex) FetchMarketsByType(typ string, params map[string]interface{}) []interface{} {
	if typ == "option" {
		underlying := self.ApiFuncReturnList("optionGetUnderlying", params, nil, nil)
		result := []interface{}{}
		for i := 0; i < self.Length(underlying); i++ {
			response := self.ApiFunc("optionGetInstrumentsUnderlying", map[string]interface{}{
				"underlying": self.Member(underlying, i),
			}, nil, nil)
			result = self.ArrayConcat(result, response)
		}
		return self.ParseMarkets(result)
	} else if self.ToBool(typ == "spot" || typ == "futures" || typ == "swap") {
		method := typ + "GetInstruments"
		response := self.ApiFuncReturnList(method, params, nil, nil)
		return self.ParseMarkets(response)
	} else {
		self.RaiseException("NotSupported", self.Id+" fetchMarketsByType does not support market type "+typ)
	}

	return nil
}

func (self *Okex) FetchCurrencies(params map[string]interface{}) map[string]interface{} {
	response := self.ApiFunc("accountGetCurrencies", params, nil, nil)
	result := map[string]interface{}{}
	for i := 0; i < self.Length(response); i++ {
		currency := self.Member(response, i)
		id := self.SafeString(currency, "currency", "")
		code := self.SafeCurrencyCode(id)
		precision := 8
		name := self.SafeString(currency, "name", "")
		canDeposit := self.SafeInteger(currency, "can_deposit", 0)
		canWithdraw := self.SafeInteger(currency, "can_withdraw", 0)
		active := false
		if canDeposit != 0 && canWithdraw != 0 {
			active = true
		}
		self.SetValue(result, code, map[string]interface{}{
			"id":        id,
			"code":      code,
			"info":      currency,
			"type":      nil,
			"name":      name,
			"active":    active,
			"fee":       nil,
			"precision": precision,
			"limits": map[string]interface{}{
				"amount": map[string]interface{}{
					"min": nil,
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
				"withdraw": map[string]interface{}{
					"min": self.SafeFloat(currency, "min_withdrawal", 0),
					"max": nil,
				},
			},
		})
	}
	return result
}

func (self *Okex) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (orderBook *OrderBook, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	method := market.Type + "GetInstrumentsInstrumentId"
	if market.Type == "swap" {
		method += "Depth"
	} else {
		method += "Book"
	}
	request := map[string]interface{}{
		"instrument_id": self.Member(market, "id"),
	}
	if self.ToBool(!self.TestNil(limit)) {
		self.SetValue(request, "size", limit)
	}
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	timestamp := self.Parse8601(self.SafeString(response, "timestamp", ""))
	return self.ParseOrderBook(response, timestamp, "bids", "asks", 0, 1), nil
}

func (self *Okex) ParseAccountBalance(response interface{}) *Account {
	result := map[string]interface{}{
		"info": response,
	}
	for i := 0; i < self.Length(response); i++ {
		balance := self.Member(response, i)
		currencyId := self.SafeString(balance, "currency", "")
		code := self.SafeCurrencyCode(currencyId)
		account := self.Account()
		self.SetValue(account, "total", self.SafeFloat(balance, "balance", 0))
		self.SetValue(account, "used", self.SafeFloat(balance, "hold", 0))
		self.SetValue(account, "free", self.SafeFloat(balance, "available", 0))
		self.SetValue(result, code, account)
	}
	return self.ParseBalance(result)
}

/*
TODO: 杠杆现货暂不支持，因为返回的格式不是Account类型，难以兼容
func (self *Okex) ParseMarginBalance(response interface{}) *Account {
	result := map[string]interface{}{
		"info": response,
	}
	for i := 0; i < self.Length(response); i++ {
		balance := self.Member(response, i).(map[string]interface{})
		marketId := self.SafeString(balance, "instrument_id", "")
		market := self.SafeValue(self.MarketsById, marketId, nil)
		var symbol interface{}
		if self.ToBool(self.TestNil(market)) {
			baseId, quoteId := self.Unpack2(strings.Split(marketId, "-"))
			base := self.SafeCurrencyCode(baseId)
			quote := self.SafeCurrencyCode(quoteId)
			symbol = base + "/" + quote
		} else {
			symbol = self.Member(market, "symbol")
		}
		omittedBalance := self.Omit(balance, []interface{}{"instrument_id", "liquidation_price", "product_id", "risk_rate", "margin_ratio", "maint_margin_ratio", "tiers"})
		keys := reflect.ValueOf(omittedBalance).MapKeys()
		accounts := map[string]interface{}{}
		for k := 0; k < self.Length(keys); k++ {
			key := self.Member(keys, k).(string)
			marketBalance := self.Member(balance, key)
			if strings.Contains(key, ":") {
				parts := strings.Split(key, ":")
				currencyId := self.Member(parts, 1)
				code := self.SafeCurrencyCode(currencyId)
				account := self.Account()
				self.SetValue(account, "total", self.SafeFloat(marketBalance, "balance", 0))
				self.SetValue(account, "used", self.SafeFloat(marketBalance, "hold", 0))
				self.SetValue(account, "free", self.SafeFloat(marketBalance, "available", 0))
				self.SetValue(accounts, code, account)
			} else {
				self.RaiseException("NotSupported", self.Id+" margin balance response format has changed!")
			}
		}
		self.SetValue(result, symbol.(string), self.ParseBalance(accounts))
	}
	return result
}
*/

func (self *Okex) ParseFuturesBalance(response interface{}) *Account {
	result := map[string]interface{}{
		"info": response,
	}
	info := self.SafeValue(response, "info", map[string]interface{}{})
	ids := reflect.ValueOf(info).MapKeys()
	for i := 0; i < self.Length(ids); i++ {
		id := self.Member(ids, i)
		code := self.SafeCurrencyCode(id)
		balance := self.SafeValue(info, id, map[string]interface{}{})
		account := self.Account()
		self.SetValue(account, "total", self.SafeFloat(balance, "equity", 0))
		self.SetValue(account, "free", self.SafeFloat(balance, "total_avail_balance", 0))
		self.SetValue(result, code, account)
	}
	return self.ParseBalance(result)
}

func (self *Okex) ParseSwapBalance(response interface{}) *Account {
	result := map[string]interface{}{
		"info": response,
	}
	info := self.SafeValue(response, "info", []interface{}{})
	for i := 0; i < self.Length(info); i++ {
		balance := self.Member(info, i)
		marketId := self.SafeString(balance, "instrument_id", "")
		symbol := marketId
		if self.ToBool(self.InMap(marketId, self.MarketsById)) {
			symbol = self.Member(self.Member(self.MarketsById, marketId), "symbol").(string)
		}
		account := self.Account()
		self.SetValue(account, "total", self.SafeFloat(balance, "equity", 0))
		self.SetValue(account, "free", self.SafeFloat(balance, "total_avail_balance", 0))
		self.SetValue(result, symbol, account)
	}
	return self.ParseBalance(result)
}

func (self *Okex) ParseBalanceByType(typ string, response interface{}) *Account {
	if self.ToBool(typ == "account" || typ == "spot") {
		return self.ParseAccountBalance(response)
	} else if self.ToBool(typ == "margin") {
		//return self.ParseMarginBalance(response)
	} else if self.ToBool(typ == "futures") {
		return self.ParseFuturesBalance(response)
	} else if self.ToBool(typ == "swap") {
		return self.ParseSwapBalance(response)
	}
	self.RaiseException("NotSupported", self.Id+" fetchBalance does not support the "+typ+" type (the type must be one of account, spot, margin, futures, swap)")
	return nil
}

func (self *Okex) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	defaultType := self.SafeString2(self.Options, "fetchBalance", "defaultType", "")
	typ := self.SafeString(params, "type", defaultType)
	if self.ToBool(self.TestNil(typ)) {
		self.RaiseException("ArgumentsRequired", self.Id+" fetchBalance requires a type parameter (one of account, spot, margin, futures, swap)")
	}
	self.LoadMarkets()
	suffix := "Accounts"
	if typ == "account" {
		suffix = "Wallet"
	}
	method := typ + "Get" + suffix
	query := self.Omit(params, "type")
	response := self.ApiFuncReturnList(method, query, nil, nil)
	return self.ParseBalanceByType(typ, response), nil
}

func (self *Okex) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	self.LoadMarkets()
	market := self.Market(symbol)
	request := map[string]interface{}{
		"instrument_id": market.Id,
	}
	clientOrderId := self.SafeString2(params, "client_oid", "clientOrderId", "")
	if self.ToBool(!self.TestNil(clientOrderId)) {
		self.SetValue(request, "client_oid", clientOrderId)
		params = self.Omit(params, []interface{}{"client_oid", "clientOrderId"})
	}
	var method string
	if market.Future || market.Swap {
		size := self.AmountToPrecision(symbol, amount)
		if market.Future {
			size = self.NumberToString(amount)
		}
		request = self.Extend(request, map[string]interface{}{
			"type":  typ,
			"size":  size,
			"price": self.PriceToPrecision(symbol, price),
		}).(map[string]interface{})
		if market.Future {
			self.SetValue(request, "leverage", "10")
		}
		method = market.Type + "PostOrder"
	} else {
		marginTrading := self.SafeString(params, "margin_trading", "1")
		request = self.Extend(request, map[string]interface{}{
			"side":           side,
			"type":           typ,
			"margin_trading": marginTrading,
		}).(map[string]interface{})
		if self.ToBool(typ == "limit") {
			// TODO:
			self.SetValue(request, "price", self.PriceToPrecision(symbol, price))
			self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
			//self.SetValue(request, "price", price)
			//self.SetValue(request, "size", amount)
		} else if self.ToBool(typ == "market") {
			// TODO, 市价单暂无用到先不做，涉及TICK_SIZE的decimalToPrecision，js的实现太复杂了，要翻译很麻烦
			/*
				if self.ToBool(side == "buy") {
					notional := self.SafeFloat(params, "notional", 0)
					createMarketBuyOrderRequiresPrice := self.SafeValue(self.Options, "createMarketBuyOrderRequiresPrice", True)
					if self.ToBool(createMarketBuyOrderRequiresPrice) {
						if self.ToBool(!self.TestNil(price)) {
							if self.ToBool(self.TestNil(notional)) {
								notional = amount * price
							}
						} else if self.ToBool(self.TestNil(notional)) {
							self.RaiseException("InvalidOrder", self.Id+" createOrder() requires the price argument with market buy orders to calculate total order cost (amount to spend), where cost = amount * price. Supply a price argument to createOrder() call if you want the cost to be calculated for you from price and amount, or, alternatively, add .options[createMarketBuyOrderRequiresPrice] = false and supply the total cost value in the amount argument or in the notional extra parameter (the exchange-specific behaviour)")
						}
					} else {
						notional = self.IfThenElse(self.ToBool(self.TestNil(notional)), amount, notional).(float64)
					}
					precision := market.Precision.Price
					request["notional"], _ = DecimalToPrecision(notional, Truncate, precision, 2, DecimalPlaces, NoPadding))
				} else {
					self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
				}
			*/
		}
		method = self.IfThenElse(self.ToBool(marginTrading == "2"), "marginPostOrders", "spotPostOrders").(string)
	}
	response := self.ApiFunc(method, self.Extend(request, params), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Okex) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	if self.ToBool(self.TestNil(symbol)) {
		self.RaiseException("ArgumentsRequired", self.Id+" cancelOrder() requires a symbol argument")
	}
	self.LoadMarkets()
	market := self.Market(symbol)
	defaultType := self.SafeString2(self.Options, "cancelOrder", "defaultType", market.Type)
	typ := self.SafeString(params, "type", defaultType)
	if self.ToBool(self.TestNil(typ)) {
		self.RaiseException("ArgumentsRequired", self.Id+" cancelOrder requires a type parameter (one of spot, margin, futures, swap).")
	}
	method := typ + "PostCancelOrder"
	request := map[string]interface{}{
		"instrument_id": market.Id,
	}
	if market.Future || market.Swap {
		method += "InstrumentId"
	} else {
		method += "s"
	}
	clientOrderId := self.SafeString2(params, "client_oid", "clientOrderId", "")
	if self.ToBool(!self.TestNil(clientOrderId)) {
		method += "ClientOid"
		self.SetValue(request, "client_oid", clientOrderId)
	} else {
		method += "OrderId"
		self.SetValue(request, "order_id", id)
	}
	query := self.Omit(params, []interface{}{"type", "client_oid", "clientOrderId"})
	response = self.ApiFunc(method, self.Extend(request, query), nil, nil)
	result := self.IfThenElse(self.ToBool(self.InMap("result", response)), response, self.SafeValue(response, market.Id, map[string]interface{}{}))
	return self.ParseOrder(result, market), nil
}

func (self *Okex) ParseOrderStatus(status string) string {
	statuses := map[string]interface{}{
		"-2": "failed",
		"-1": "canceled",
		"0":  "open",
		"1":  "open",
		"2":  "closed",
		"3":  "open",
		"4":  "canceled",
	}
	return self.SafeString(statuses, status, status)
}

func (self *Okex) ParseOrderSide(side string) string {
	sides := map[string]string{
		"1": "buy",  // open long
		"2": "sell", // open short
		"3": "sell", // close long
		"4": "buy",  // close short
	}
	return self.SafeString(sides, side, side)
}

func (self *Okex) ParseOrder(order interface{}, market interface{}) (result map[string]interface{}) {
	id := self.SafeString(order, "order_id", "")
	timestamp := self.Parse8601(self.SafeString(order, "timestamp", ""))
	side := self.SafeString(order, "side", "")
	typ := self.SafeString(order, "type", "")
	if self.ToBool(side != "buy" && side != "sell") {
		side = self.ParseOrderSide(typ)
	}
	if self.ToBool(typ != "limit" && typ != "market") {
		if self.ToBool(self.InMap("pnl", order)) {
			typ = "futures"
		} else {
			typ = "swap"
		}
	}
	var symbol interface{}
	marketId := self.SafeString(order, "instrument_id", "")
	if self.ToBool(self.InMap(marketId, self.MarketsById)) {
		market = self.Member(self.MarketsById, marketId)
		symbol = self.Member(market, "symbol")
	} else {
		symbol = marketId
	}
	if self.ToBool(!self.TestNil(market)) {
		if self.ToBool(self.TestNil(symbol)) {
			symbol = self.Member(market, "symbol")
			symbol = market.(*Market).Symbol
		}
	}
	amount := self.SafeFloat(order, "size", 0)
	filled := self.SafeFloat2(order, "filled_size", "filled_qty", 0.0)
	var remaining interface{}
	if self.ToBool(!self.TestNil(amount)) {
		if self.ToBool(!self.TestNil(filled)) {
			amount = math.Max(amount, filled)
			remaining = math.Max(0, amount-filled)
		}
	}
	if self.ToBool(typ == "market") {
		remaining = 0
	}
	cost := self.SafeFloat2(order, "filled_notional", "funds", 0.0)
	price := self.SafeFloat(order, "price", 0)
	average := self.SafeFloat(order, "price_avg", 0)
	if self.ToBool(self.TestNil(cost)) {
		if self.ToBool(!self.TestNil(filled) && !self.TestNil(average)) {
			cost = average * filled
		}
	} else {
		if self.ToBool(self.TestNil(average) && !self.TestNil(filled) && filled > 0) {
			average = cost / filled
		}
	}
	status := self.ParseOrderStatus(self.SafeString(order, "state", ""))
	feeCost := self.SafeFloat(order, "fee", 0)
	var fee interface{}
	if self.ToBool(!self.TestNil(feeCost)) {
		var feeCurrency interface{}
		fee = map[string]interface{}{
			"cost":     feeCost,
			"currency": feeCurrency,
		}
	}
	clientOrderId := self.SafeString(order, "client_oid", "")
	return map[string]interface{}{
		"info":               order,
		"id":                 id,
		"clientOrderId":      clientOrderId,
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

func (self *Okex) FetchOrder(id string, symbol string, params map[string]interface{}) (result *Order, err error) {
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
	if self.ToBool(self.TestNil(typ)) {
		self.RaiseException("ArgumentsRequired", self.Id+" fetchOrder requires a type parameter (one of spot, margin, futures, swap).")
	}
	instrumentId := ""
	if market.Future || market.Swap {
		instrumentId = "InstrumentId"
	}
	method := typ + "GetOrders" + instrumentId
	request := map[string]interface{}{
		"instrument_id": self.Member(market, "id"),
	}
	clientOid := self.SafeString(params, "client_oid", "")
	if self.ToBool(!self.TestNil(clientOid)) {
		method += "ClientOid"
		self.SetValue(request, "client_oid", clientOid)
	} else {
		method += "OrderId"
		self.SetValue(request, "order_id", id)
	}
	query := self.Omit(params, "type")
	response := self.ApiFunc(method, self.Extend(request, query), nil, nil)
	return self.ToOrder(self.ParseOrder(response, market)), nil
}

func (self *Okex) FetchOrdersByState(state string, symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}) {
	if self.ToBool(self.TestNil(symbol)) {
		self.RaiseException("ArgumentsRequired", self.Id+" fetchOrdersByState requires a symbol argument")
	}
	self.LoadMarkets()
	market := self.Market(symbol)
	defaultType := self.SafeString2(self.Options, "fetchOrder", "defaultType", market.Type)
	typ := self.SafeString(params, "type", defaultType)
	if self.ToBool(self.TestNil(typ)) {
		self.RaiseException("ArgumentsRequired", self.Id+" fetchOrder requires a type parameter (one of spot, margin, futures, swap).")
	}
	request := map[string]interface{}{
		"instrument_id": self.Member(market, "id"),
		"state":         state,
	}
	method := typ + "GetOrders"
	if market.Future || market.Swap {
		method += "InstrumentId"
	}
	query := self.Omit(params, "type")
	response := self.ApiFuncReturnList(method, self.Extend(request, query), nil, nil)
	if self.ToBool(self.Member(market, "type") == "swap" || self.Member(market, "type") == "futures") {
		orders = self.SafeValue(response, "order_info", []interface{}{})
	} else {
		orders = response
		responseLength := self.Length(response)
		if self.ToBool(responseLength < 1) {
			return []interface{}{}
		}
		if self.ToBool(responseLength > 1) {
			before := self.SafeValue(self.Member(response, 1), "before", nil)
			if self.ToBool(!self.TestNil(before)) {
				orders = self.Member(response, 0)
			}
		}
	}
	return self.ParseOrders(orders, market, since, limit)
}

func (self *Okex) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (result []*Order, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = self.PanicToError(e)
		}
	}()
	return self.ToOrders(self.FetchOrdersByState("6", symbol, since, limit, params)), nil
}

func (self *Okex) GetPathAuthenticationType(path string) string {
	// https://github.com/ccxt/ccxt/issues/6651
	// a special case to handle the optionGetUnderlying interefering with
	// other endpoints containing this keyword
	if path == "underlying" {
		return "public"
	}
	auth := self.SafeValue(self.Options, "auth", "")
	key := self.FindBroadlyMatchedKey(auth, path)
	return self.SafeString(auth, key, "private")
}

func (self *Okex) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}) {
	// TODO: only support map params
	request := "/api/" + api + "/" + self.Version + "/"
	request += self.ImplodeParams(path, params)
	query := self.Omit(params, self.ExtractParams(path))
	url := self.ImplodeParams(self.Member(self.Urls["api"], "rest").(string), map[string]interface{}{
		"hostname": self.Hostname,
	}) + request
	typ := self.GetPathAuthenticationType(path)
	if self.ToBool(typ == "public") {
		if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
			url += "?" + self.Urlencode(query)
		}
	} else if self.ToBool(typ == "private") {
		self.CheckRequiredCredentials()
		timestamp := self.Iso8601Okex(self.Milliseconds())
		headers = map[string]interface{}{
			"OK-ACCESS-KEY":        self.ApiKey,
			"OK-ACCESS-PASSPHRASE": self.Password,
			"OK-ACCESS-TIMESTAMP":  timestamp,
		}
		auth := timestamp + method + request
		if self.ToBool(method == "GET") {
			if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
				urlencodedQuery := "?" + self.Urlencode(query)
				url += urlencodedQuery
				auth += urlencodedQuery
			}
		} else {
			if self.Length(query) > 0 {
				body = self.Json(query)
				auth += body.(string)
			}
			self.SetValue(headers, "Content-Type", "application/json")
		}
		signature := self.Hmac(self.Encode(auth), self.Encode(self.Secret), "sha256", "base64")
		self.SetValue(headers, "OK-ACCESS-SIGN", self.Decode(signature))
	}
	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    body,
		"headers": headers,
	}
}

func (self *Okex) HandleErrors(httpCode int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
	fmt.Println(response)
	feedback := self.Id + " " + body
	if self.ToBool(httpCode == 503) {
		self.RaiseException("ExchangeNotAvailable", feedback)
	}
	if self.ToBool(!self.ToBool(response)) {
		return
	}
	message := self.SafeString(response, "message", "")
	errorCode := self.SafeString2(response, "code", "error_code", "")
	if self.ToBool(!self.TestNil(message)) {
		self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), message, feedback)
		self.ThrowBroadlyMatchedException(self.Member(self.Exceptions, "broad"), message, feedback)
		self.ThrowExactlyMatchedException(self.Member(self.Exceptions, "exact"), errorCode, feedback)
		nonEmptyMessage := message != ""
		nonZeroErrorCode := !self.TestNil(errorCode) && errorCode != "0"
		if self.ToBool(nonZeroErrorCode || nonEmptyMessage) {
			self.RaiseException("ExchangeError", feedback)
		}
	}
}
