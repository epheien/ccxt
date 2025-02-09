package margin_bitmax

import (
	. "github.com/epheien/ccxt/go/base"
	"github.com/epheien/ccxt/go/bitmax"
)

type MarginBitmax struct {
	bitmax.Bitmax
}

func New(config *ExchangeConfig) (ex *MarginBitmax, err error) {
	ex = new(MarginBitmax)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		ex = nil
		return
	}

	return
}

func (self *MarginBitmax) Describe() []byte {
	return []byte(`{
    "id": "margin_bitmax",
    "name": "BitMax",
    "countries": [
        "CN"
    ],
    "rateLimit": 500,
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
        "logo": "https://user-images.githubusercontent.com/1294454/66820319-19710880-ef49-11e9-8fbe-16be62a11992.jpg",
        "api": "https://bitmax.io",
        "test": "https://bitmax-test.io",
        "www": "https://bitmax.io",
        "doc": [
            "https://bitmax-exchange.github.io/bitmax-pro-api/#bitmax-pro-api-documentation"
        ],
        "fees": "https://bitmax.io/#/feeRate/tradeRate",
        "referral": "https://bitmax.io/#/register?inviteCode=EL6BXBQM"
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
        "accountGroup": {
            "get": [
                "{account-category}/balance",
                "cash/balance",
                "margin/balance",
                "margin/risk",
                "transfer",
                "futures/collateral-balance",
                "futures/position",
                "futures/risk",
                "futures/funding-payments",
                "{account-category}/order/open",
                "{account-category}/order/status",
                "{account-category}/order/hist/current",
                "order/hist"
            ],
            "post": [
                "futures/transfer/deposit",
                "futures/transfer/withdraw",
                "{account-category}/order",
                "{account-category}/order/batch"
            ],
            "delete": [
                "{account-category}/order",
                "{account-category}/order/all",
                "{account-category}/order/batch"
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
            "tierBased": true,
            "percentage": true,
            "taker": 0.001,
            "maker": 0.001
        }
    },
    "precisionMode": "TICK_SIZE",
    "options": {
        "account-category": "margin",
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
        "BTCBEAR": "BEAR",
        "BTCBULL": "BULL"
    }
}`)
}
