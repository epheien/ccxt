package margin_kucoin

import (
	. "github.com/georgexdz/ccxt/go/base"
	"github.com/georgexdz/ccxt/go/kucoin"
)

type MarginKucoin struct {
	kucoin.Kucoin
}

func New(config *ExchangeConfig) (ex *MarginKucoin, err error) {
	ex = new(MarginKucoin)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *MarginKucoin) Describe() []byte {
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
        "tradeType": "MARGIN_TRADE",
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
