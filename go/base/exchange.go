package base

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/imdario/mergo"
	"github.com/thoas/go-funk"
	"os"

	//"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	urllib "net/url"
	"reflect"
	"regexp"
	"sync"
	"syscall"

	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"

	"github.com/georgexdz/ccxt/go/encoding/jsonx"
	"github.com/tidwall/gjson"
)

//type JSONTime int64

type SignInfo struct {
	Url     string
	Method  string
	Body    interface{} // []byte or string
	Headers map[string]interface{}
}

// Market struct
type Market struct {
	Id             string  `json:"id"`     // exchange specific
	Symbol         string  `json:"symbol"` // ccxt unified
	Base           string  `json:"base"`
	BaseNumericId  string  `json:"baseNumericId"`
	Quote          string  `json:"quote"`
	QuoteNumericId string  `json:"quoteNumericId"`
	BaseId         string  `json:"baseId"`  // from bitmex
	QuoteId        string  `json:"quoteId"` // from bitmex
	Active         bool    `json:"active"`  // from bitmex
	Taker          float64 `json:"taker"`   // from bitmex
	Maker          float64 `json:"maker"`   // from bitmex
	Type           string  `json:"type"`    // from bitmex
	Spot           bool    `json:"spot"`    // from bitmex
	Swap           bool    `json:"swap"`    // from bitmex
	Future         bool    `json:"future"`  // from bitmex
	Option         bool
	Prediction     bool        `json:"prediction"` // from bitmex
	Precision      Precision   `json:"precision"`
	Limits         Limits      `json:"limits"`
	Lot            float64     `json:"lot"`
	Info           interface{} `json:"info"`
}

// Precision struct
type Precision struct {
	Price       int `json:"price"` // 精度, 即小数位数, 例如 0.01 即为2, 1 即为 0
	PriceSacle  int `priceScale`   // 缩放, 10 的整数倍, 其他均为非法值. > 1 时才使用. 例如 0.25 => amount: 25, scale 100
	Amount      int `json:"amount"`
	AmountScale int `amountScale`
	Base        int `json:"base"`
	Quote       int `json:"quote"`
}

// Limits struct
type Limits struct {
	Amount MinMax `json:"amount"` // 最小最大下单数量
	Price  MinMax `json:"price"`  // 最小最大下单价格
	Cost   MinMax `json:"cost"`   // 最小最大下单交易额
}

// MinMax struct
type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ExchangeConfig for main configuration
// Timeout takes json value in milliseconds
type ExchangeConfig struct {
	ApiKey          string        `json:"apiKey"`
	Password        string        `json:"password"`
	Uid             string        `json:"uid"`
	Secret          string        `json:"secret"`
	Timeout         time.Duration `json:"timeout"`
	EnableRateLimit bool          `json:"enableRateLimit"`
	Test            bool          `json:"test"`
	Verbose         bool          `json:"verbose"`
}

// ExchangeInfo for the exchange
type ExchangeInfo struct {
	Id                               string         `json:"id"`
	Name                             string         `json:"name"`
	Countries                        StringSlice    `json:"countries"`
	Version                          string         `json:"version"`
	RateLimit                        int            `json:"rateLimit"`
	Has                              HasDescription `json:"has"`
	Urls                             map[string]interface{}
	Api                              Apis              `json:"api"`
	Timeframes                       map[string]string `json:"timeframes"`
	Fees                             map[string]interface{}
	UserAgents                       map[string]string `json:"userAgents"`
	Header                           http.Header       `json:"header"`
	Proxy                            string            `json:"proxy"`
	Origin                           string            `json:"origin"`
	Limits                           Limits            `json:"limits"`
	Precision                        Precision         `json:"precision"`
	Exceptions                       map[string]interface{}
	DontGetUsedBalanceFromStaleCache bool `json:"dontGetUsedBalanceFromStaleCache"`
	CommonCurrencies                 map[string]string
}

// Apis public and private methods
type Apis struct {
	Public  ApiMethods `json:"public"`
	Private ApiMethods `json:"private"`
}

// ApiMethods get/post/put/delete urls
type ApiMethods struct {
	Get    map[string]string `json:"get"`
	Post   map[string]string `json:"post"`
	Put    map[string]string `json:"put"`
	Delete map[string]string `json:"delete"`
}

// Urls for exchange
type Urls struct {
	WWW  string      `json:"www"`
	Test string      `json:"test"`
	Api  string      `json:"api"`
	Logo StringSlice `json:"logo"`
	Doc  StringSlice `json:"doc"`
	Fees StringSlice `json:"fees"`
}

// status = {
//     'status': 'ok',
//     'updated': None,
//     'eta': None,
//     'url': None,
// }
type ExchangeStatus struct {
	Status  string // "ok", "maintenance"
	Updated int64  // 更新时间戳?
}

// Exception takes the string and applies the error method
type Exception map[string]error

// UnmarshalJSON accepts strings and links to the appropaite error method:
func (e Exception) UnmarshalJSON(b []byte) error {
	var errTypes map[string]string
	err := json.Unmarshal(b, &errTypes)
	if err != nil {
		return err
	}
	for msg, errType := range errTypes {
		if e == nil {
			e = make(map[string]error)
		}
		e[msg] = TypedError(errType, msg)
	}
	return nil
}

// HasDescription for exchange functionality
type HasDescription struct {
	CancelAllOrders      bool `json:"cancelAllOrders"`
	CancelOrder          bool `json:"cancelOrder"`
	CancelOrders         bool `json:"cancelOrders"`
	CORS                 bool `json:"CORS"`
	CreateDepositAddress bool `json:"createDepositAddress"`
	CreateLimitOrder     bool `json:"createLimitOrder"`
	CreateMarketOrder    bool `json:"createMarketOrder"`
	CreateOrder          bool `json:"createOrder"`
	Deposit              bool `json:"deposit"`
	EditOrder            bool `json:"editOrder"`
	FetchBalance         bool `json:"fetchBalance"`
	FetchBidsAsks        bool `json:"fetchBidsAsks"`
	FetchClosedOrders    bool `json:"fetchClosedOrders"`
	FetchCurrencies      bool `json:"fetchCurrencies"`
	FetchDepositAddress  bool `json:"fetchDepositAddress"`
	FetchDeposits        bool `json:"fetchDeposits"`
	FetchFundingFees     bool `json:"fetchFundingFees"`
	FetchL2OrderBook     bool `json:"fetchL2OrderBook"`
	FetchLedger          bool `json:"fetchLedger"`
	FetchMarkets         bool `json:"fetchMarkets"`
	FetchMyTrades        bool `json:"fetchMyTrades"`
	FetchOHLCV           bool `json:"fetchOHLCV"`
	FetchOpenOrders      bool `json:"fetchOpenOrders"`
	FetchOrder           bool `json:"fetchOrder"`
	FetchOrderBook       bool `json:"fetchOrderBook"`
	FetchOrderBooks      bool `json:"fetchOrderBooks"`
	FetchOrders          bool `json:"fetchOrders"`
	FetchTicker          bool `json:"fetchTicker"`
	FetchTickers         bool `json:"fetchTickers"`
	FetchTrades          bool `json:"fetchTrades"`
	FetchTradingFee      bool `json:"fetchTradingFee"`
	FetchTradingFees     bool `json:"fetchTradingFees"`
	FetchTradingLimits   bool `json:"fetchTradingLimits"`
	FetchTransactions    bool `json:"fetchTransactions"`
	FetchWithdrawals     bool `json:"fetchWithdrawals"`
	PrivateApi           bool `json:"privateApi"`
	PublicApi            bool `json:"publicApi"`
	Withdraw             bool `json:"withdraw"`
}

// StringSlice a custom type for handling variable JSON
type StringSlice []string

// UnmarshalJSON accepts both forms for StringSlice:
//   - ["s1", "s2"...]
//   - "s"
// For the latter, ss will hold a slice of one element "s"
// todo: unify to array form ?
func (ss *StringSlice) UnmarshalJSON(b []byte) (err error) {
	// try slice unmarshal
	var slice []string
	err = json.Unmarshal(b, &slice)
	if err == nil {
		*ss = slice
		return nil
	}
	// try string unmarshal
	var s string
	err2 := json.Unmarshal(b, &s)
	if err2 == nil {
		*ss = append(*ss, s)
		return nil
	}
	// return original error
	return err
}

// ApiUrls for different types of urls
type ApiUrls struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

// UnmarshalJSON accepts both forms for ApiUrls:
//   - {"public": "urlpub", "private": "urlpriv"} or
//   - "url"
// For the latter, "url" is assigned to both a.Private and a.Public
// todo: unify to struct form ?
func (a *ApiUrls) UnmarshalJSON(b []byte) (err error) {
	// default struct unmarshal to compatible type
	type t ApiUrls
	err = json.Unmarshal(b, (*t)(a))
	if err == nil {
		return nil
	}
	// string unmarshal
	var s string
	err2 := json.Unmarshal(b, &s)
	if err2 == nil {
		a.Private = s
		a.Public = s
		return nil
	}
	// return original error
	return err
}

// Fees for using the exchange
type Fees struct {
	Trading TradingFees `json:"trading"`
	Funding FundingFees `json:"funding"`
}

// TradingFees on the exchange
type TradingFees struct {
	TierBased  bool             `json:"tierBased"`
	Percentage bool             `json:"percentage"`
	Maker      float64          `json:"maker"`
	Taker      float64          `json:"taker"`
	Tiers      TradingFeesTiers `json:"tiers"`
}

// TradingFeesTiers on the exchange
type TradingFeesTiers struct {
	Taker [][2]float64 `json:"taker"`
	Maker [][2]float64 `json:"maker"`
}

// FundingFees for using the exchange
type FundingFees struct {
	TierBased  bool               `json:"tierBased"`
	Percentage bool               `json:"percentage"`
	Deposit    map[string]float64 `json:"deposit"`
	Withdraw   map[string]float64 `json:"withdraw"`
}

// Balance details
type Balance struct {
	Free      float64 `json:"free"`
	Used      float64 `json:"used"`
	RealPnl   float64 `json:"realPnl"`
	UnrealPnl float64 `json:"unrealPnl"`
	Total     float64 `json:"total"`
}

// Account details
type Account struct {
	Free    map[string]float64 `json:"free"`
	Used    map[string]float64 `json:"used"`
	Total   map[string]float64 `json:"total"`
	Account map[string]*Balance
	Info    interface{}
}

type MarkPrice struct {
	Symbol     string
	MarkPrice  float64
	IndexPrice float64
	Timestamp  int64
	Info       interface{}
}

type Position struct {
	Symbol     string
	Side       string
	Leverage   float64
	Price      float64
	Amount     float64
	UsedAmount float64
	UnrealPnl  float64
	RealPnl    float64
	Info       interface{}
}

// Order structure
type Order struct {
	Id            string      `json:"id"`
	ClientOrderId string      `json:"clientOrderId"`
	Timestamp     int64       `json:"timestamp"`
	Datetime      string      `json:"datetime"`
	Symbol        string      `json:"symbol"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	Side          string      `json:"side"`
	Price         float64     `json:"price"`
	Cost          float64     `json:"cost"`
	Amount        float64     `json:"amount"`
	Filled        float64     `json:"filled"`
	Remaining     float64     `json:"remaining"`
	Average       float64     `json:"average"`
	Fee           float64     `json:"fee"`
	Info          interface{} `json:"info"`
}

func (o *Order) InitFromMap(m map[string]interface{}) (result *Order) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: 需要提取出来具体是什么错误
			panic(r)
		}
	}()

	for k, v := range m {
		if v == nil {
			continue
		}

		switch k {
		case "id":
			o.Id = v.(string)
		case "symbol":
			var ok bool
			// 有时候这个是 reflect.Value
			o.Symbol, ok = v.(string)
			if !ok {
				o.Symbol = fmt.Sprint(v)
			}
		case "type":
			o.Type = v.(string)
		case "side":
			o.Side = v.(string)
		case "price":
			o.Price = v.(float64)
		case "amount":
			o.Amount = v.(float64)
		case "cost":
			o.Cost = v.(float64)
		case "filled":
			o.Filled = v.(float64)
		case "remaining":
			o.Remaining = v.(float64)
		case "timestamp":
			o.Timestamp = v.(int64)
		case "datetime":
			o.Datetime = v.(string)
		case "fee":
			// NOTE: fee 有可能是字典也可能是浮点, 暂时无视
			//o.Fee = v.(float64)
		case "average":
			o.Average = v.(float64)
		case "status":
			o.Status = v.(string)
		case "clientOrderId":
			o.ClientOrderId = v.(string)
		case "info":
			o.Info = v
		default:
			// ignore
		}
	}
	result = o
	return
}

// OrderBook struct
type OrderBook struct {
	Asks      [][2]float64
	Bids      [][2]float64
	Timestamp int64
	Datetime  string
	Nonce     int64
}

// BookEntry struct
type BookEntry struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// UnmarshalJSON accepts both forms for BookEntry:
//   - []float64 of size 2 or
//   - default BookEntry struct
func (o *BookEntry) UnmarshalJSON(b []byte) (err error) {
	// []float64 unmarshal
	var f []float64
	err = json.Unmarshal(b, &f)
	if err == nil {
		if len(f) != 2 {
			return fmt.Errorf("UnmarshalJSON: expected 2 floats for BookEntry, got %d", len(f))
		}
		o.Price, o.Amount = f[0], f[1]
		return nil
	}
	// default struct unmarshal to compatible type
	type t BookEntry
	err2 := json.Unmarshal(b, (*t)(o))
	if err2 == nil {
		return nil
	}
	return err
}

// Trade struct
type Trade struct {
	Id        string      `json:"id"`
	Symbol    string      `json:"symbol"`
	Amount    float64     `json:"amount"`
	Price     float64     `json:"price"`
	Timestamp int64       `json:"timestamp"`
	Datetime  string      `json:"datetime"`
	Order     string      `json:"order"` // ignore
	Type      string      `json:"type"`  // ignore
	Side      string      `json:"side"`
	Info      interface{} `json:"info"`
}

// Ticker struct
type Ticker struct {
	Symbol      string      `json:"symbol"`
	Timestamp   int64       `json:"timestamp"`
	Datetime    string      `json:"datetime"`
	Last        float64     `json:"last"`
	Ask         float64     `json:"ask"`
	Bid         float64     `json:"bid"`
	AskQty      float64     `json:"askQty"`
	BidQty      float64     `json:"bidQty"`
	Open        float64     `json:"open"`
	High        float64     `json:"high"`
	Low         float64     `json:"low"`
	Close       float64     `json:"close"`
	BaseVolume  float64     `json:"baseVolume"`
	QuoteVolume float64     `json:"quoteVolume"`
	Vwap        float64     `json:"vwap"`
	Change      float64     `json:"change"`
	Percentage  float64     `json:"percentage"`
	Info        interface{} `json:"info"`
}

// Currency struct
type Currency struct {
	Id        string `json:"id"`
	Code      string `json:"code"`
	NumericId string `json:"numericId"`
	Precision int    `json:"precision"`
}

// DepositAddress struct
type DepositAddress struct {
	Currency string      `json:"currency"`
	Address  string      `json:"address"`
	Status   string      `json:"status"`
	Info     interface{} `json:"info"`
}

type ApiDecode struct {
	Api    string // public, private, ...
	Method string // GET, POST, ...
	Path   string // margin/asset, ...
}

// OHLCV open, high, low, close, volume
type OHLCV struct {
	Timestamp int64       `json:"timestamp"`
	Open      float64     `json:"open"`
	High      float64     `json:"high"`
	Low       float64     `json:"low"`
	Close     float64     `json:"close"`
	Volume    float64     `json:"volume"`
	Info      interface{} `json:"info"`
}

// UnmarshalJSON accepts both forms for OHLCV:
//   - default struct or
//   - []float64 of size 6
func (o *OHLCV) UnmarshalJSON(b []byte) (err error) {
	// default struct unmarshal to compatible type
	type t OHLCV
	err = json.Unmarshal(b, (*t)(o))
	if err == nil {
		return nil
	}
	// []float64 unmarshal
	var f []float64
	err2 := json.Unmarshal(b, &f)
	if err2 != nil {
		return err2 // return float64 unmarshal error as it's the expected form
	}
	if len(f) != 6 {
		return fmt.Errorf("UnmarshalJSON: expected 6 floats for OHLCV, got %d", len(f))
	}

	err2 = json.Unmarshal([]byte(fmt.Sprintf("%f", f[0])), &o.Timestamp)
	if err2 != nil {
		return fmt.Errorf("UnmarshalJSON: couldn't unmarshal timestamp: %s", err2)
	}
	o.Open, o.High, o.Low, o.Close, o.Volume = f[1], f[2], f[3], f[4], f[5]
	return nil
}

// Exchange is a common interface of methods
type ExchangeInterface interface {
	FetchTicker(symbol string, params map[string]interface{}) (*Ticker, error)
	FetchTickers(symbols []string, params map[string]interface{}) ([]*Ticker, error)
	FetchOHLCV(symbol, timeframe string, since int64, limit int64, params map[string]interface{}) ([]*OHLCV, error)
	FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (*OrderBook, error)
	FetchStatus(params map[string]interface{}) (*ExchangeStatus, error) // 默认实现为返回 ok 状态
	// FetchL2OrderBook(symbol string, limit *int, params map[string]interface{}) (OrderBook, error)
	FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) ([]*Trade, error)
	FetchOrder(id string, symbol string, params map[string]interface{}) (*Order, error)
	// FetchOrders(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Order, error)
	FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) ([]*Order, error)
	// FetchClosedOrders(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Order, error)
	// FetchMyTrades(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Trade, error)
	FetchBalance(params map[string]interface{}) (*Account, error)
	FetchPositions(symbol string, params map[string]interface{}) ([]*Position, error)
	FetchMarkPrice(symbol string, params map[string]interface{}) (*MarkPrice, error)
	//FetchCurrencies() (map[string]*Currency, error)
	FetchMarkets(params map[string]interface{}) ([]*Market, error)
	FetchAccounts(params map[string]interface{}) []interface{}

	CreateOrder(symbol, otype, side string, amount float64, price float64, params map[string]interface{}) (*Order, error)
	LimitBuy(symbol string, price, amount float64, params map[string]interface{}) (*Order, error)
	LimitSell(symbol string, price, amount float64, params map[string]interface{}) (*Order, error)
	CancelOrder(id string, symbol string, params map[string]interface{}) (interface{}, error)

	// Describe() []byte
	//GetMarkets() map[string]*Market
	SetMarkets([]*Market, map[string]interface{}) map[string]*Market
	//GetMarketsById() map[string]Market
	//SetMarketsById(map[string]Market)
	//GetCurrencies() map[string]Currency
	//SetCurrencies(map[string]Currency)
	//GetCurrenciesById() map[string]Currency
	//SetCurrenciesById(map[string]Currency)
	//SetSymbols([]string)
	//SetIds([]string)
	// GetOrders() []Order
	LoadMarkets() map[string]*Market
	// LoadMarkets(reload bool, params map[string]interface{}) (map[string]*Market, error)
	// GetMarket(symbol string) (Market, error)
	// CreateLimitBuyOrder(symbol string, amount float64, price *float64, params map[string]interface{}) (Order, error)
	// CreateLimitSellOrder(symbol string, amount float64, price *float64, params map[string]interface{}) (Order, error)
	// CreateMarketBuyOrder(symbol string, amount float64, params map[string]interface{}) (Order, error)
	// CreateMarketSellOrder(symbol string, amount float64, params map[string]interface{}) (Order, error)

	SetApiKey(string)
	SetSecret(string)
	SetPassword(string)
	SetUid(string)
	SetBaseUrl(string)
	BaseUrl(key string) string
	SetVerbose(verbose bool)

	FetchCurrencies(params map[string]interface{}) map[string]interface{}
	ApiFunc(function string, params interface{}, headers map[string]interface{}, body interface{}) (response map[string]interface{})
	ApiFuncRaw(function string, params map[string]interface{}, headers map[string]interface{}, body interface{}) (response []byte)
	SetHttpLib(lib string) // fasthttp, net/http
}

type ExchangeInterfaceInternal interface {
	ExchangeInterface
	// 返回值类型一般来说是 map[string]interface{}
	Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) interface{}
	ApiFuncDecode(function string) (path string, api string, method string)
	ApiFunc(function string, params interface{}, headers map[string]interface{}, body interface{}) (response map[string]interface{})
	ApiFuncReturnList(function string, params interface{}, headers map[string]interface{}, body interface{}) (response []interface{})
	// 返回原始 []byte, 一般用于给外部使用
	ApiFuncRaw(function string, params map[string]interface{}, headers map[string]interface{}, body interface{}) (response []byte)
	// 底层 http 调用, 状态号非 200 会抛出错误, 除非是不被关注的错误 (参考 Exchange.httpExceptions)
	Fetch(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, jsonResponse interface{})
	FetchViaFastHttp(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, jsonResponse interface{})
	Request(path string, api string, method string, params map[string]interface{}, headers map[string]interface{}, body interface{}) (response interface{})
	Describe() []byte
	ParseTrade(interface{}, *Market) *Trade
	ParseOrder(interface{}, interface{}) map[string]interface{}
	HandleErrors(code int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{})
	Market(string) *Market
}

// Exchange struct
type Exchange struct {
	sync.RWMutex
	ExchangeInfo
	ExchangeConfig

	EnableFasthttp bool
	Client         *http.Client
	FastHttpClient *fasthttp.Client

	Markets        map[string]*Market
	MarketsById    map[string]*Market
	Ids            []string
	Symbols        []string
	Currencies     map[string]*Currency
	CurrenciesById map[string]*Currency
	Accounts       []interface{}
	AccountsById   map[string]interface{}

	Child         ExchangeInterfaceInternal
	ApiDecodeInfo map[string]*ApiDecode
	//ApiUrls        map[string]string
	DescribeMap    map[string]interface{}
	Options        map[string]interface{}
	httpExceptions map[string]string
	Hostname       string
	requestTimeout time.Duration

	DescribeJson gjson.Result
}

func (self *Exchange) SetHttpLib(lib string) {
	switch strings.ToLower(lib) {
	case "fasthttp":
		self.EnableFasthttp = true
	default:
		self.EnableFasthttp = false
	}
}

func (self *Exchange) SetProxy(proxy string) {
	self.FastHttpClient.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
	// 以下用于支持 net/http 库
	os.Setenv("HTTP_PROXY", proxy)
	os.Setenv("HTTPS_PROXY", proxy)
}

func (self *Exchange) Init(config *ExchangeConfig) (err error) {
	self.Child = self

	if config != nil {
		self.ExchangeConfig = *config
	}

	// 默认超时时间 10 秒
	self.requestTimeout = 10 * time.Second
	if self.ExchangeConfig.Timeout > 0 {
		self.requestTimeout = self.ExchangeConfig.Timeout
	}

	// @ 初始化 net/http 客户端
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	self.Client = &http.Client{
		Transport: tr,
		Timeout:   self.requestTimeout,
	}

	// @ 初始化 fasthttp 客户端
	readTimeout := self.requestTimeout
	writeTimeout := self.requestTimeout
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	self.FastHttpClient = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	self.httpExceptions = map[string]string{
		"422": "ExchangeError",
		"418": "DDoSProtection",
		"429": "RateLimitExceeded",
		"404": "ExchangeNotAvailable",
		"409": "ExchangeNotAvailable",
		"410": "ExchangeNotAvailable",
		"500": "ExchangeNotAvailable",
		"501": "ExchangeNotAvailable",
		"502": "ExchangeNotAvailable",
		"520": "ExchangeNotAvailable",
		"521": "ExchangeNotAvailable",
		"522": "ExchangeNotAvailable",
		"525": "ExchangeNotAvailable",
		"526": "ExchangeNotAvailable",
		"400": "ExchangeNotAvailable",
		"403": "ExchangeNotAvailable",
		"405": "ExchangeNotAvailable",
		"503": "ExchangeNotAvailable",
		"530": "ExchangeNotAvailable",
		"408": "RequestTimeout",
		"504": "RequestTimeout",
		"401": "AuthenticationError",
		"511": "AuthenticationError",
	}

	return
}

func (self *Exchange) Describe() []byte {
	return nil
}

func (self *Exchange) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	return nil, errors.New("FetchMarkets not supported yet")
}

func (self *Exchange) ToMarket(market map[string]interface{}) *Market {
	return self.MarketFromMap(market)
}

func (self *Exchange) ToMarkets(markets []interface{}) []*Market {
	result := []*Market{}
	for _, market := range markets {
		result = append(result, self.ToMarket(market.(map[string]interface{})))
	}
	return result
}

func (self *Exchange) FetchTicker(symbol string, params map[string]interface{}) (*Ticker, error) {
	return nil, errors.New("FetchTicker not supported yet")
}

func (self *Exchange) FetchTickers(symbols []string, params map[string]interface{}) ([]*Ticker, error) {
	return nil, errors.New("FetchTickers not supported yet")
}

func (self *Exchange) FetchOHLCV(symbol, timeframe string, since int64, limit int64, params map[string]interface{}) ([]*OHLCV, error) {
	return nil, errors.New("FetchOHLCV not supported yet")
}

func (self *Exchange) FetchOrderBook(symbol string, limit int64, params map[string]interface{}) (*OrderBook, error) {
	return nil, errors.New("FetchOrderBook not supported yet")
}

func (self *Exchange) FetchStatus(params map[string]interface{}) (*ExchangeStatus, error) {
	return &ExchangeStatus{Status: "ok", Updated: self.Milliseconds()}, nil
}

func (self *Exchange) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) interface{} {
	return nil
}

func (self *Exchange) MarketId(symbol string) string {
	market := self.Child.Market(symbol)
	if market != nil {
		return market.Id
	} else {
		return symbol
	}
}

func (self *Exchange) MarketFromMap(o interface{}) *Market {
	p := &Market{}

	if m, ok := o.(map[string]interface{}); ok {
		p.Id = m["id"].(string)
		p.Symbol = m["symbol"].(string)
		p.Base = m["base"].(string)
		//p.BaseNumericId
		p.Quote = m["quote"].(string)
		//p.QuoteNumericId
		p.BaseId = m["baseId"].(string)
		p.QuoteId = m["quoteId"].(string)
		if m["active"] != nil {
			p.Active = m["active"].(bool)
		}
		if m["taker"] != nil {
			p.Taker = m["taker"].(float64)
		}
		if m["maker"] != nil {
			p.Maker = m["maker"].(float64)
		}
		if m["type"] != nil {
			p.Type = m["type"].(string)
		}
		if m["spot"] != nil {
			p.Spot = m["spot"].(bool)
		}
		if m["swap"] != nil {
			p.Swap = m["swap"].(bool)
		}
		if m["futures"] != nil {
			p.Future = m["futures"].(bool)
		}
		if m["option"] != nil {
			p.Option = m["option"].(bool)
		}
		//p.Prediction
		if m["precision"] != nil {
			precisionMap := m["precision"].(map[string]interface{})
			if precisionMap["amount"] != nil {
				p.Precision.Amount = precisionMap["amount"].(int)
			}
			if precisionMap["price"] != nil {
				p.Precision.Price = precisionMap["price"].(int)
			}
		}
		//p.Limits
		if m["limits"] != nil {
			limitsMap := m["limits"].(map[string]interface{})
			if limitsMap["amount"] != nil {
				p.Limits.Amount.Min = self.SafeFloat(limitsMap["amount"], "min")
				p.Limits.Amount.Max = self.SafeFloat(limitsMap["amount"], "max")
			}
			if limitsMap["price"] != nil {
				p.Limits.Price.Min = self.SafeFloat(limitsMap["price"], "min")
				p.Limits.Price.Max = self.SafeFloat(limitsMap["price"], "max")
			}
			if limitsMap["cost"] != nil {
				p.Limits.Cost.Min = self.SafeFloat(limitsMap["cost"], "min")
				p.Limits.Cost.Max = self.SafeFloat(limitsMap["cost"], "max")
			}
		}
		//p.Lot
		p.Info = m["info"]
	}

	return p
}

func (self *Exchange) SetMarkets(markets []*Market, currencies map[string]interface{}) map[string]*Market {
	symbols := make([]string, len(markets))
	Ids := make([]string, len(markets))
	marketsBySymbol := make(map[string]*Market, len(markets))
	marketsById := make(map[string]*Market, len(markets))
	baseCurrencies := make([]*Currency, 0)
	quoteCurrencies := make([]*Currency, 0)

	for i, market := range markets {
		marketsBySymbol[market.Symbol] = market
		marketsById[market.Id] = market
		symbols[i] = market.Symbol
		Ids[i] = market.Id
		// currency logic
		baseCurrency := new(Currency)
		if market.Base != "" {
			baseCurrency.Id = market.BaseId
			if baseCurrency.Id == "" {
				baseCurrency.Id = market.Base
			}
			baseCurrency.NumericId = market.BaseNumericId
			baseCurrency.Code = market.Base
			switch {
			case market.Precision.Base != 0:
				baseCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				baseCurrency.Precision = market.Precision.Amount
			default:
				baseCurrency.Precision = 8
			}
		}
		baseCurrencies = append(baseCurrencies, baseCurrency)
		quoteCurrency := new(Currency)
		if market.Quote != "" {
			quoteCurrency.Id = market.QuoteId
			if baseCurrency.Id == "" {
				quoteCurrency.Id = market.Quote
			}
			quoteCurrency.NumericId = market.QuoteNumericId
			quoteCurrency.Code = market.Quote
			switch {
			case market.Precision.Base != 0:
				quoteCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				quoteCurrency.Precision = market.Precision.Amount
			default:
				quoteCurrency.Precision = 8
			}
		}
		quoteCurrencies = append(quoteCurrencies, quoteCurrency)
	}
	allCurrencies := append(baseCurrencies, quoteCurrencies...)
	groupedCurrencies := make(map[string][]*Currency)
	for _, currency := range allCurrencies {
		groupedCurrencies[currency.Code] = append(groupedCurrencies[currency.Code], currency)
	}
	sortedCurrencies := make(map[string]*Currency)
	for code, currencies := range groupedCurrencies {
		for _, currency := range currencies {
			if sortedCurrencies[code] == nil {
				continue
			}
			if sortedCurrencies[code].Id == "" {
				sortedCurrencies[code] = currency
			}
			if sortedCurrencies[code].Precision < currency.Precision {
				sortedCurrencies[code] = currency
			}
		}
	}

	sort.Strings(symbols)
	sort.Strings(Ids)

	self.Symbols = symbols
	self.Ids = Ids
	self.MarketsById = marketsById
	self.Markets = marketsBySymbol

	if len(currencies) == 0 {
		xCurrencies := self.Currencies
		if xCurrencies == nil {
			xCurrencies = make(map[string]*Currency)
		}
		for code, currency := range sortedCurrencies {
			xCurrencies[code] = currency
		}
		self.Currencies = xCurrencies
	} else {
		self.Currencies = sortedCurrencies
	}
	currenciesById := self.CurrenciesById
	if len(currenciesById) == 0 {
		currenciesById = make(map[string]*Currency, len(currencies))
	}
	for _, currency := range sortedCurrencies {
		currenciesById[currency.Id] = currency
	}
	self.CurrenciesById = currenciesById
	return self.Markets
}

// func (self *Exchange) LoadMarkets(reload bool, params map[string]interface{}) (map[string]*Market, error) {
func (self *Exchange) LoadMarkets() map[string]*Market {
	if self.Markets != nil {
		return self.Markets
	}

	var currencies map[string]interface{}
	hasfetchCurrencies := self.DescribeMap["has"].(map[string]interface{})["fetchCurrencies"]
	if hasfetchCurrencies != nil && hasfetchCurrencies.(bool) {
		currencies = self.Child.FetchCurrencies(map[string]interface{}{})
	}

	markets, err := self.Child.FetchMarkets(nil)
	if err != nil {
		if err.Error() == "FetchMarkets not supported yet" {
			return map[string]*Market{}
		}
		self.RaiseException("ExchangeError", fmt.Sprintf("failed to FetchMarkets(): %s", err))
	}
	return self.Child.SetMarkets(markets, currencies)
}

func (self *Exchange) LoadAccounts() []interface{} {
	//self.Lock()
	//defer self.Unlock()
	if len(self.Accounts) > 0 {
		return self.Accounts
	}
	accounts := self.Child.FetchAccounts(nil)
	for _, account := range accounts {
		one := map[string]interface{}{
			"id":    account.(map[string]interface{})["id"],
			"state": account.(map[string]interface{})["state"],
			"type":  account.(map[string]interface{})["type"],
		}
		self.Accounts = append(self.Accounts, one)
	}
	self.AccountsById = self.IndexBy(self.Accounts, "id")
	return self.Accounts
}

func (self *Exchange) Request(
	path string,
	api string,
	method string,
	params map[string]interface{},
	headers map[string]interface{},
	body interface{},
) (response interface{}) {
	signInfo := self.Child.Sign(path, api, method, params, headers, body)

	url := self.Member(signInfo, "url").(string)
	// 使用新值覆盖
	method = self.Member(signInfo, "method").(string)
	var rawResp []byte

	if self.EnableFasthttp {
		rawResp, response = self.Child.FetchViaFastHttp(
			url,
			method,
			self.Member(signInfo, "headers").(map[string]interface{}),
			self.Member(signInfo, "body"),
		)
	} else {
		rawResp, response = self.Child.Fetch(
			url,
			method,
			self.Member(signInfo, "headers").(map[string]interface{}),
			self.Member(signInfo, "body"),
		)
	}

	self.HandleRestResponse(string(rawResp), response, url, method)

	return
}

func (self *Exchange) PrepareRequestHeaders(req *http.Request, headers map[string]interface{}) {
	//req.Header.Set("Accept-Encoding", "gzip, deflate")

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
}

func (self *Exchange) FetchViaFastHttp(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, jsonResponse interface{}) {
	var rbody []byte
	if body != nil {
		switch body.(type) {
		case string:
			rbody = []byte(body.(string))
		case []byte:
			rbody = body.([]byte)
		default:
			self.RaiseException("InternalError", fmt.Sprintf("Invalid Argument body: %v", body))
			return
		}
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.SetBodyRaw(rbody)

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}

	if self.Verbose {
		log.Println("Request:", method, url, headers, body)
	}

	resp := fasthttp.AcquireResponse()
	err := self.FastHttpClient.DoTimeout(req, resp, self.requestTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		msg := fmt.Sprintf("%v %v %v", method, url, err)
		if err == fasthttp.ErrTimeout {
			self.RaiseException("RequestTimeout", msg)
		} else if err == fasthttp.ErrNoFreeConns {
			self.RaiseException("NetworkError", msg)
		} else if err == fasthttp.ErrConnectionClosed {
			self.RaiseException("NetworkError", msg)
		} else {
			errName := reflect.TypeOf(err).String()
			if errName == "*net.OpError" {
				// Write and Read errors are not so often and in fact they just mean timeout problems
				self.RaiseException("RequestTimeout", msg)
			}
		}
		self.RaiseException("ExchangeError", msg)
	}

	response = resp.Body()
	strRawResp := string(response)

	if self.Verbose {
		log.Println("Response:", method, url, resp.StatusCode(), resp.Header.String(), strRawResp)
	}

	status := fasthttp.StatusMessage(resp.StatusCode())
	// 这里忽略错误, 上层做类型断言的时候会产生错误信息
	json.Unmarshal(response, &jsonResponse)
	// FIXME: header 需要转为 map[string]interface{} 或 map[string][]string
	self.Child.HandleErrors(int64(resp.StatusCode()), status, url, method, resp.Header, strRawResp, jsonResponse, headers, body)
	if resp.StatusCode() != fasthttp.StatusOK {
		self.HandleRestErrors(resp.StatusCode(), status, strRawResp, url, method)
	}

	return
}

func (self *Exchange) Fetch(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, jsonResponse interface{}) {
	var rbody []byte
	if body != nil {
		switch body.(type) {
		case string:
			rbody = []byte(body.(string))
		case []byte:
			rbody = body.([]byte)
		default:
			self.RaiseException("InternalError", fmt.Sprintf("Invalid Argument body: %v", body))
			return
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(rbody))
	if err != nil {
		self.RaiseException("InternalError", fmt.Sprintf("NewRequest err: %v", err))
		return
	}

	self.PrepareRequestHeaders(req, headers)

	if self.Verbose {
		log.Println("Request:", method, url, headers, body)
	}

	resp, err := self.Client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			self.RaiseException("RequestTimeout", fmt.Sprintf("%v %v %v", method, url, err))
		}
		switch t := err.(type) {
		case syscall.Errno:
			if t == syscall.ECONNREFUSED {
				self.RaiseException("NetworkError", fmt.Sprintf("%v %v %v", method, url, err))
			}
		default:
			self.RaiseException("ExchangeError", fmt.Sprintf("%v %v %v", method, url, err))
		}
	}

	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		self.RaiseException("InternalError", fmt.Sprintf("read response err: %v", err))
	}

	strRawResp := string(response)

	if self.Verbose {
		log.Println("Response:", method, url, resp.StatusCode, resp.Header, strRawResp)
	}

	// 这里忽略错误, 上层做类型断言的时候会产生错误信息
	json.Unmarshal(response, &jsonResponse)
	self.Child.HandleErrors(int64(resp.StatusCode), resp.Status, url, method, resp.Header, strRawResp, jsonResponse, headers, body)
	if resp.StatusCode != http.StatusOK {
		self.HandleRestErrors(resp.StatusCode, resp.Status, strRawResp, url, method)
	}

	return
}

func (self *Exchange) RegSplit(text string, delimeter string) (result []string) {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result = make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func (self *Exchange) DefineRestApi() (err error) {
	// NOTE: HACK! 初始化一些必要的成员变量
	self.Id = self.DescribeMap["id"].(string)
	self.Name = self.DescribeMap["name"].(string)

	self.ApiDecodeInfo = make(map[string]*ApiDecode)

	if jsonApiInfo, ok := self.DescribeMap["api"].(map[string]interface{}); ok {
		for strApi, apiInfo := range jsonApiInfo {
			if methodInfo, ok := apiInfo.(map[string]interface{}); ok {
				for strMethod, methodInfo := range methodInfo {
					if pathList, ok := methodInfo.([]interface{}); ok {
						for _, path := range pathList {
							if strPath, ok := path.(string); ok {
								var strDealPath string
								splitParts := self.RegSplit(strPath, "[^a-zA-Z0-9]")
								for _, part := range splitParts {
									strDealPath += strings.Title(part)
								}
								self.ApiDecodeInfo[strApi+strings.Title(strMethod)+strDealPath] = &ApiDecode{Api: strApi, Method: strings.ToUpper(strMethod), Path: strPath}

								if self.Verbose {
									//log.Println("\napiDecodeInfo:", strApi, strPath, strMethod, strDealPath)
								}
							}
						}
					}
				}
			}
		}
	}

	return
}

func (self *Exchange) ApiFuncDecode(function string) (path string, api string, method string) {
	if info, ok := self.ApiDecodeInfo[function]; ok {
		return info.Path, info.Api, info.Method
	} else {
		self.RaiseException("InternalError", fmt.Sprintf("func %v not found!", function))
	}
	return
}

func (self *Exchange) ApiFunc(function string, params interface{}, headers map[string]interface{}, body interface{}) (result map[string]interface{}) {
	path, api, method := self.Child.ApiFuncDecode(function)
	if params == nil {
		params = map[string]interface{}{}
	}
	return self.Child.Request(path, api, method, params.(map[string]interface{}), headers, body).(map[string]interface{})
}

func (self *Exchange) ApiFuncReturnList(function string, params interface{}, headers map[string]interface{}, body interface{}) (result []interface{}) {
	path, api, method := self.Child.ApiFuncDecode(function)
	if params == nil {
		params = map[string]interface{}{}
	}
	return self.Child.Request(path, api, method, params.(map[string]interface{}), headers, body).([]interface{})
}

func (self *Exchange) ApiFuncRaw(function string, params map[string]interface{}, headers map[string]interface{}, body interface{}) (response []byte) {
	path, api, method := self.Child.ApiFuncDecode(function)

	signInfo := self.Child.Sign(path, api, method, params, headers, body)
	url := self.Member(signInfo, "url").(string)
	method = self.Member(signInfo, "method").(string)
	response, _ = self.Child.Fetch(
		url,
		method,
		self.Member(signInfo, "headers").(map[string]interface{}),
		self.Member(signInfo, "body"),
	)

	var resp interface{}
	// 这里忽略错误, 上层做类型断言的时候会产生错误信息
	json.Unmarshal(response, &resp)
	self.HandleRestResponse(string(response), resp, url, method)

	return
}

func (self *Exchange) Parse8601(x string) int64 {
	t, err := time.Parse(time.RFC3339, x)
	if err != nil {
		self.RaiseInternalException("Parse8601 " + x + " err!")
	}
	return t.Unix()
}

func (self *Exchange) Iso8601Okex(milliseconds int64) string {
	var seconds int64
	seconds = milliseconds / 1000
	return time.Unix(seconds, 0).In(time.UTC).Format("2006-01-02T15:04:05.070Z")
}

func (self *Exchange) Iso8601(milliseconds int64) string {
	var seconds int64
	// 0 表示无效值, 而不是 1970-01-01 这种迷惑性的时间
	if milliseconds == 0 {
		return ""
	}
	seconds = milliseconds / 1000
	return time.Unix(seconds, 0).Format("2006-01-02T15:04:05-0700")
}

func (self *Exchange) Milliseconds() int64 {
	return time.Now().UnixNano() / 1000000
}

// Exchanges returns the available exchanges
func Exchanges() []string {
	available := []string{"bitmex"}
	return available
}

func MapValues(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	keys := v.MapKeys()
	output := []interface{}{}
	for i, l := 0, v.Len(); i < l; i++ {
		output = append(output, v.MapIndex(keys[i]))
	}
	return output
}

func getCurrencyUsedOnOpenOrders(currency string) float64 {
	// TODO: implement getCurrencyUsedOnOpenOrders()
	return 0.0
}

func SortSliceByIndex(s [][2]float64, idx int, descending bool) {
	if !descending {
		sort.Slice(s, func(i, j int) bool {
			// edge cases
			if len(s[i]) == 0 && len(s[j]) == 0 {
				return false // two empty slices - so one is not less than other i.e. false
			}
			if len(s[i]) == 0 || len(s[j]) == 0 {
				return len(s[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
			}

			// both slices len() > 0, so can test this now:
			return s[i][idx] < s[j][idx]
		})
	} else {
		sort.Slice(s, func(i, j int) bool {
			// edge cases
			if len(s[i]) == 0 && len(s[j]) == 0 {
				return false // two empty slices - so one is not less than other i.e. false
			}
			if len(s[i]) == 0 || len(s[j]) == 0 {
				return len(s[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
			}

			// both slices len() > 0, so can test this now:
			return s[i][idx] >= s[j][idx]
		})
	}
}

func ToInteger(v interface{}) int64 {
	switch v.(type) {
	case int:
		return int64(v.(int))
	case int64:
		return v.(int64)
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case string:
		vStr := v.(string)
		vv, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("ToInteger error (%s): %v", err.Error(), v))
		}
		return vv
	default:
		panic(fmt.Sprintf("ToInteger error: %v", v))
	}
}

func ToFloat(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, err := strconv.ParseFloat(vStr, 64)
		if err != nil {
			panic(fmt.Sprintf("ToFloat error (%s): %v", err.Error(), v))
		}
		return vF
	default:
		panic(fmt.Sprintf("ToFloat error: %v", v))
	}
}

func (self *Exchange) ParseBidsAsks(bidsAsks []interface{}, priceKey int64, amountKey int64) (out [][2]float64) {
	if len(bidsAsks) == 0 {
		return
	}

	if _, ok := bidsAsks[0].([]interface{}); ok {
		for _, one := range bidsAsks {
			if bidAsk, ok := one.([]interface{}); ok {
				price := bidAsk[priceKey]
				amount := bidAsk[amountKey]
				if price != "" && amount != "" {
					priceF := ToFloat(price)
					amountF := ToFloat(amount)
					out = append(out, [2]float64{priceF, amountF})
				}
			}
		}
	} else {
		self.RaiseException("ExchangeError", "unrecognized bidask format: "+fmt.Sprint(bidsAsks[0]))
	}

	return
}

func (self *Exchange) Extend(maps ...interface{}) interface{} {
	if len(maps) == 0 {
		return make(map[string]interface{})
	}
	output := make(map[string]interface{})
	for _, m := range maps {
		if oneMap, ok := m.(map[string]interface{}); ok {
			for k, v := range oneMap {
				output[k] = v
			}
		}
	}
	return output
}

func (self *Exchange) InMap(k interface{}, m interface{}) bool {
	if strMap, ok := m.(map[string]interface{}); ok {
		if _, ok := strMap[k.(string)]; ok {
			return true
		}
	}
	return false
}

func (self *Exchange) ToBool(v interface{}) bool {
	if v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
		return !self.TestNil(v)
	} else {
		return false
	}
}

func (self *Exchange) SafeList(m map[string]interface{}, key string, defaultVal []interface{}) (val []interface{}) {
	if val, ok := m[key]; ok {
		if listVal, ok := val.([]interface{}); ok {
			return listVal
		}
	}
	return defaultVal
}

func (self *Exchange) SafeValue(m interface{}, key interface{}, defaultVal ...interface{}) (val interface{}) {
	var def interface{}
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}

	switch key.(type) {
	case string:
		if mm, ok := m.(map[string]interface{}); ok {
			if val, ok := mm[key.(string)]; ok {
				return val
			}
		}
	case int, int64, int8, int32:
		if li, ok := m.([]interface{}); ok {
			idx := int(ToInteger(key))
			if idx >= 0 && idx < len(li) {
				return li[idx]
			}
		}
	}

	return def
}

func NestedMapLookup(m map[string]interface{}, ks ...string) (rval interface{}, err error) {
	var ok bool

	if len(ks) == 0 { // degenerate input
		return nil, fmt.Errorf("NestedMapLookup needs at least one key")
	}
	if rval, ok = m[ks[0]]; !ok {
		return nil, fmt.Errorf("key not found; remaining keys: %v", ks)
	} else if len(ks) == 1 { // we've reached the final key
		return rval, nil
	} else if m, ok = rval.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return NestedMapLookup(m, ks[1:]...)
	}
}

func (self *Exchange) ParseOrderBook(orderBook interface{}, timestamp int64, bidsKey string, asksKey string, priceKey int64, amountKey int64) *OrderBook {
	var result OrderBook

	if orderBookMap, ok := orderBook.(map[string]interface{}); ok {
		if bids, ok := orderBookMap[bidsKey]; ok {
			if bidsList, ok := bids.([]interface{}); ok {
				result.Bids = self.ParseBidsAsks(bidsList, priceKey, amountKey)
				SortSliceByIndex(result.Bids, 0, true)
			}
		}
		if asks, ok := orderBookMap[asksKey]; ok {
			if asksList, ok := asks.([]interface{}); ok {
				result.Asks = self.ParseBidsAsks(asksList, priceKey, amountKey)
				SortSliceByIndex(result.Asks, 0, false)
			}
		}
		result.Timestamp = timestamp
		if result.Timestamp > 0 {
			result.Datetime = self.Iso8601(result.Timestamp)
		}

		return &result
	}
	return nil
}

func (self *Exchange) SafeInteger(d interface{}, key string, def ...int64) (ret int64) {
	var defaultVal int64
	if len(def) > 0 {
		defaultVal = def[0]
	}
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key]; ok {
			if intVal, ok := val.(int); ok {
				return int64(intVal)
			} else if intVal, ok := val.(int64); ok {
				return intVal
			} else if floatVal, ok := val.(float64); ok {
				return int64(floatVal)
			} else if v, ok := val.(string); ok {
				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("SafeInteger error (%s): %s", err.Error(), v))
				}
				return i
			}
		}
	}
	return defaultVal
}

func (self *Exchange) SafeInteger2(d interface{}, key1 string, key2 string, defaultVal int64) int64 {
	return ToInteger(self.SafeEither(d, key1, key2, defaultVal))
}

func (self *Exchange) SafeBool(d interface{}, key string, def ...bool) bool {
	var defaultVal bool
	if len(def) > 0 {
		defaultVal = def[0]
	}
	return self.ToBool(self.SafeValue(d, key, defaultVal))
}

func (self *Exchange) SafeFloat2(d interface{}, key1 string, key2 string, defaultVal float64) float64 {
	return ToFloat(self.SafeEither(d, key1, key2, defaultVal))
}

func (self *Exchange) SafeString2(d interface{}, key1 string, key2 string, defaultVal string) string {
	return self.SafeEither(d, key1, key2, defaultVal).(string)
}

func (self *Exchange) SafeValue2(d interface{}, key1 string, key2 string, defaultVal interface{}) interface{} {
	return self.SafeEither(d, key1, key2, defaultVal)
}

func (self *Exchange) SafeEither(d interface{}, key1 string, key2 string, defaultVal interface{}) interface{} {
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key1]; ok {
			return val
		}
		if val, ok := d[key2]; ok {
			return val
		}
	}
	return defaultVal
}

func (self *Exchange) NumberToString(v interface{}) string {
	return NumberToString(v.(float64))
}

func (self *Exchange) SafeString(d interface{}, key string, def ...string) string {
	defaultVal := ""
	if len(def) > 0 {
		defaultVal = def[0]
	}
	if d, ok := d.(map[string]interface{}); ok {
		val := d[key]
		if val != nil {
			switch val.(type) {
			case int:
				return strconv.Itoa(val.(int))
			case int64:
				return strconv.FormatInt(val.(int64), 10)
			}
			return fmt.Sprintf("%v", val)
		}
	}
	if d, ok := d.([]string); ok {
		if idx, err := strconv.Atoi(key); err != nil {
			val := d[idx]
			return fmt.Sprintf("%v", val)
		}
	}
	return defaultVal
}

func (self *Exchange) SafeStringLower(d interface{}, key string, defaultVal string) string {
	return strings.ToLower(self.SafeString(d, key, defaultVal))
}

func (self *Exchange) SafeFloat(d interface{}, key string, def ...float64) (result float64) {
	defaultVal := 0.0
	if len(def) > 0 {
		defaultVal = def[0]
	}
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key]; ok {
			switch val.(type) {
			case string:
				fVal, err := strconv.ParseFloat(val.(string), 64)
				if err == nil {
					return fVal
				}
			case int:
				return float64(val.(int))
			case int64:
				return float64(val.(int64))
			case float32:
				return float64(val.(float32))
			case float64:
				return val.(float64)
			case nil:
				return defaultVal
			}
		}
	}
	return defaultVal
}

func (self *Exchange) SafeNumber(d interface{}, key string, def_ ...float64) float64 {
	def := 0.0
	if len(def_) > 0 {
		def = def_[0]
	}
	return self.SafeFloat(d, key, def)
}

func (self *Exchange) SafeNumber2(d interface{}, key1, key2 string, def_ ...float64) float64 {
	def := 0.0
	if len(def_) > 0 {
		def = def_[0]
	}
	return self.SafeFloat2(d, key1, key2, def)
}

func (self *Exchange) Omit(d map[string]interface{}, args interface{}) (result map[string]interface{}) {
	if argList, ok := args.([]string); ok {
		for _, arg := range argList {
			if _, ok := d[arg]; ok {
				delete(d, arg)
			}
		}
		return d
	}

	if arg, ok := args.(string); ok {
		delete(d, arg)
		return d
	}

	return d
}

func (self *Exchange) ExtractParams(s string) (result []string) {
	r := regexp.MustCompile(`{([^{}]*)}`)
	matches := r.FindAllStringSubmatch(s, -1)
	for _, v := range matches {
		result = append(result, v[1])
	}
	return result
}

func (self *Exchange) ImplodeParams(s string, params interface{}) string {
	if paramsMap, ok := params.(map[string]interface{}); ok {
		for k, v := range paramsMap {
			s = strings.ReplaceAll(s, "{"+k+"}", fmt.Sprintf("%v", v))
		}
	}
	return s
}

func (self *Exchange) ImplodeHostname(s string) string {
	return self.ImplodeParams(s, map[string]interface{}{"hostname": self.Hostname})
}

var hashers = map[string]func() hash.Hash{
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
	"sha384": sha512.New384,
	"md5":    md5.New,
}

// Hash encodes the payload based on the available hashing algo
func Hash(payload, algo, encoding string) (string, error) {
	if hashers[algo] == nil {
		return "", fmt.Errorf("hash: unsupported algo \"%s\"", algo)
	}
	h := hashers[algo]()
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", fmt.Errorf("hash: %s", err)
	}
	return string(encode(h.Sum(nil), encoding)), nil
}

// HMAC encodes the payload based on the available hashing algo
// algo: "md5", "sha512", "sha384", "sha256", "sha1"
// encoding: "hex", "base64"
func (self *Exchange) Hmac(payload, key, algo, encoding string) string {
	if hashers[algo] == nil {
		self.RaiseException("InternalError", fmt.Sprintf("HMAC: unsupported hashing algo \"%s\"", algo))
	}
	h := hmac.New(hashers[algo], []byte(key))
	_, err := h.Write([]byte(payload))
	if err != nil {
		self.RaiseException("InternalError", fmt.Sprintf("hmac: %s", err))
	}
	return string(encode(h.Sum(nil), encoding))
}

// JWT creates a new signed token
func JWT(claims map[string]interface{}, secret string) string {
	var signer jwt.SigningMethod = jwt.SigningMethodHS256
	token := jwt.New(signer)
	token.Claims = jwt.MapClaims(claims)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		RaiseException("InternalError", fmt.Sprintf("JWT error: %v", err))
	}
	return result
}

func encode(payload []byte, encoding string) []byte {
	var result []byte
	switch encoding {
	case "hex":
		result = []byte(hex.EncodeToString(payload))
	case "base64":
		buf := make([]byte, base64.StdEncoding.EncodedLen(len(payload)))
		base64.StdEncoding.Encode(buf, payload)
		result = buf
	default:
		result = payload
	}
	return result
}

func (self *Exchange) ParseBalance(balances map[string]interface{}) (pAccount *Account) {
	var account Account
	account.Free = make(map[string]float64)
	account.Used = make(map[string]float64)
	account.Total = make(map[string]float64)
	account.Info = balances["info"]

	account.Account = map[string]*Balance{}
	for currency, balance := range self.Omit(balances, []string{"info", "free", "used", "total"}) {
		if balance, ok := balance.(map[string]interface{}); ok {
			free := self.SafeFloat(balance, "free", 0)
			used := self.SafeFloat(balance, "used", 0)
			total := self.SafeFloat(balance, "total", 0)
			realPnl := self.SafeFloat(balance, "realPnl", 0)
			unrealPnl := self.SafeFloat(balance, "unrealPnl", 0)
			account.Free[currency] = free
			account.Used[currency] = used
			account.Total[currency] = total
			account.Account[currency] = &Balance{
				Free:      free,
				Used:      used,
				Total:     total,
				RealPnl:   realPnl,
				UnrealPnl: unrealPnl,
			}
		}
	}

	return &account
}

func (self *Exchange) Uuid() string {
	return uuid.NewV4().String()
}

func (self *Exchange) PriceToPrecision(symbol string, price float64) string {
	if self.Markets[symbol] == nil {
		return self.Float64ToString(price)
	}
	ret, _ := DecimalToPrecision(price, Round, self.Markets[symbol].Precision.Price, DecimalPlaces, NoPadding)
	return ret
}

func (self *Exchange) AmountToPrecision(symbol string, amount float64) string {
	if self.Markets[symbol] == nil {
		return self.Float64ToString(amount)
	}
	ret, _ := DecimalToPrecision(amount, Truncate, self.Markets[symbol].Precision.Amount, DecimalPlaces, NoPadding)
	return ret
}

func (self *Exchange) Account() map[string]interface{} {
	return map[string]interface{}{
		"free":  nil,
		"used":  nil,
		"total": nil,
	}
}

func (self *Exchange) SafeCurrencyCode(x interface{}) string {
	code := ""

	if !self.TestNil(x) {
		currencyId := x.(string)
		if self.CurrenciesById != nil && self.CurrenciesById[currencyId] != nil {
			code = self.CurrenciesById[currencyId].Code
		} else {
			code = self.CommonCurrencyCode(strings.ToUpper(currencyId))
		}
	}

	return code
}

func (self *Exchange) Length(o interface{}) int {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Slice:
		return reflect.ValueOf(o).Len()
	case reflect.Map:
		return reflect.ValueOf(o).Len()
	default:
		return 0
	}
}

func (self *Exchange) Member(o interface{}, idx interface{}) interface{} {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Slice:
		return reflect.ValueOf(o).Index(idx.(int)).Interface()
	case reflect.Map:
		return reflect.ValueOf(o).MapIndex(reflect.ValueOf(idx)).Interface()
	case reflect.Struct:
		return reflect.ValueOf(o).FieldByName(self.Capitalize(idx.(string))).Interface()
	case reflect.Ptr:
		return reflect.Indirect(reflect.ValueOf(o)).FieldByName(self.Capitalize(idx.(string)))
	}

	return nil
}

func (self *Exchange) Market(symbol string) *Market {
	if self.Markets == nil {
		self.RaiseException("ExchangeError", self.Id+" markets not loaded")
	}

	m := self.Markets[symbol]
	if m == nil {
		self.RaiseException("BadSymbol", self.Id+" does not have market symbol "+symbol)
	}
	return m
}

func (self *Exchange) Unpack2(l interface{}) (interface{}, interface{}) {
	switch l.(type) {
	case []string:
		if ll, ok := l.([]string); ok {
			return ll[0], ll[1]
		}
	case []int64:
		if ll, ok := l.([]int64); ok {
			return ll[0], ll[1]
		}
	case []int:
		if ll, ok := l.([]int); ok {
			return ll[0], ll[1]
		}
	default:
		return nil, nil
	}
	return nil, nil
}

func (self *Exchange) IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// x == undefined
func (self *Exchange) TestNil(x interface{}) bool {
	if x == nil {
		return true
	}

	switch reflect.TypeOf(x).Kind() {
	case reflect.Map:
		if reflect.ValueOf(x).Len() == 0 {
			return true
		}
	case reflect.Slice:
		if reflect.ValueOf(x).Len() == 0 {
			return true
		}
	}

	return reflect.ValueOf(x).IsZero()
}

func (self *Exchange) SetValue(x interface{}, k string, v interface{}) {
	if m, ok := x.(map[string]interface{}); ok {
		m[k] = v
	}
}

func (self *Exchange) CheckRequiredCredentials() {
}

func (self *Exchange) UrlencodeWithArrayRepeat(i interface{}) string {
	re := regexp.MustCompile(`%5B\d*%5D`)
	return re.ReplaceAllString(self.Urlencode(i), "")
}

// NOTE: 编码后的结果都是根据key排序后的
func (self *Exchange) Urlencode(i interface{}) string {
	if m, ok := i.(map[string]interface{}); ok {
		v := urllib.Values{}

		keys := make([]string, 0, len(m))
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			val := m[k]
			v.Add(k, fmt.Sprintf("%v", val))
		}
		return v.Encode()
	}
	return ""
}

func (self *Exchange) Json(i interface{}) string {
	ret, err := json.Marshal(i)
	if err == nil {
		return string(ret)
	}
	return ""
}

func (self *Exchange) JsonIndent(i interface{}) string {
	ret, err := json.MarshalIndent(i, "", "    ")
	if err == nil {
		return string(ret)
	}
	return ""
}

func (self *Exchange) Encode(s interface{}) string {
	return s.(string)
}
func (self *Exchange) Decode(s interface{}) interface{} {
	return s
}

func (self *Exchange) AddTwoInterface(a interface{}, b interface{}) interface{} {
	if a == nil || b == nil {
		return nil
	}

	switch a.(type) {
	case string:
		return a.(string) + b.(string)
	case int:
		return a.(int) + b.(int)
	case int64:
		return a.(int64) + b.(int64)
	case float64:
		return a.(float64) + b.(float64)
	case float32:
		return a.(float32) + b.(float32)
	default:
		return nil
	}
}

func (self *Exchange) FetchMarkPrice(symbol string, params map[string]interface{}) (*MarkPrice, error) {
	return nil, fmt.Errorf("%s FetchMarkPrice not supported yet", self.Id)
}

func (self *Exchange) FetchPositions(symbol string, params map[string]interface{}) ([]*Position, error) {
	return nil, fmt.Errorf("%s FetchPositions not supported yet", self.Id)
}

func (self *Exchange) FetchBalance(params map[string]interface{}) (*Account, error) {
	return nil, fmt.Errorf("%s FetchBalance not supported yet", self.Id)
}

func (self *Exchange) CreateOrder(symbol string, otype string, side string, amount float64, price float64, params map[string]interface{}) (*Order, error) {
	return nil, fmt.Errorf("%s CreateOrder not supported yet", self.Id)
}

func (self *Exchange) LimitBuy(symbol string, price, amount float64, params map[string]interface{}) (*Order, error) {
	return self.Child.CreateOrder(symbol, "limit", "buy", amount, price, params)
}

func (self *Exchange) LimitSell(symbol string, price, amount float64, params map[string]interface{}) (*Order, error) {
	return self.Child.CreateOrder(symbol, "limit", "sell", amount, price, params)
}

func (self *Exchange) FetchCurrencies(params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{}
}

func (self *Exchange) CancelOrder(id string, symbol string, params map[string]interface{}) (interface{}, error) {
	return nil, fmt.Errorf("%s CancelOrder not supported yet", self.Id)
}

func (self *Exchange) FetchTrades(symbol string, since int64, limit int64, params map[string]interface{}) ([]*Trade, error) {
	return nil, fmt.Errorf("%s FetchTrades not supported yet", self.Id)
}

func (self *Exchange) FetchOrder(id string, symbol string, params map[string]interface{}) (*Order, error) {
	return nil, fmt.Errorf("%s FetchOrder not supported yet", self.Id)
}

func (self *Exchange) HandleErrors(code int64, reason string, url string, method string, headers interface{}, body string, response interface{}, requestHeaders interface{}, requestBody interface{}) {
}

func (self *Exchange) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) ([]*Order, error) {
	return nil, fmt.Errorf("%s FetchOpenOrders not supported yet", self.Id)
}

func (self *Exchange) SetApiKey(s string) {
	self.ApiKey = s
}
func (self *Exchange) SetSecret(s string) {
	self.Secret = s
}
func (self *Exchange) SetPassword(s string) {
	self.Password = s
}
func (self *Exchange) SetUid(s string) {
	self.Uid = s
}

func (self *Exchange) ParseTrades(trades []interface{}, market *Market, since int64, limit int64) (result []*Trade) {
	for _, trade := range trades {
		result = append(result, self.Child.ParseTrade(trade, market))
	}
	return result
}

func (self *Exchange) ParseOrders(orders interface{}, market interface{}, since int64, limit int64) (result []interface{}) {
	for _, order := range orders.([]interface{}) {
		result = append(result, self.Child.ParseOrder(order, market))
	}
	return result
}

func (self *Exchange) ParseTrade(trade interface{}, market *Market) *Trade {
	return &Trade{}
}

func (self *Exchange) ParseOrder(order interface{}, market interface{}) map[string]interface{} {
	return order.(map[string]interface{})
}

func (self *Exchange) ToOrder(order interface{}) (result *Order) {
	result = &Order{}
	return result.InitFromMap(order.(map[string]interface{}))
}

func (self *Exchange) ToOrders(orders interface{}) (result []*Order) {
	result = []*Order{}
	for _, one := range orders.([]interface{}) {
		order := (&Order{}).InitFromMap(one.(map[string]interface{}))
		result = append(result, order)
	}
	return
}

// first character only, rest characters unchanged
func (self *Exchange) Capitalize(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}

func (self *Exchange) Nonce() int64 {
	return self.Milliseconds()
}

// 0.00010000 => 4
func (self *Exchange) PrecisionFromString(s string) int {
	re := regexp.MustCompile(`0+$`)
	s = re.ReplaceAllString(s, "")
	sp := strings.Split(s, ".")
	if len(sp) > 1 {
		return len(sp[1])
	} else {
		return 0
	}
}

func RaiseException(errCls interface{}, msg interface{}) {
	panic([]string{errCls.(string), msg.(string)})
}

func (self *Exchange) RaiseInternalException(msg interface{}) {
	self.RaiseException("InternalError", msg)
}

func (self *Exchange) RaiseException(errCls interface{}, msg interface{}) {
	RaiseException(errCls, msg)
}

// key 必须为 string
func (self *Exchange) ThrowExactlyMatchedException(exact interface{}, key interface{}, message interface{}) {
	if strMap, ok := exact.(map[string]interface{}); ok {
		if val, ok := strMap[key.(string)]; ok {
			self.RaiseException(val, message)
		}
	}
}

func (self *Exchange) FindBroadlyMatchedKey(broad interface{}, key interface{}) string {
	for k, _ := range broad.(map[string]interface{}) {
		if strings.Contains(key.(string), k) {
			return k
		}
	}
	return ""
}

func (self *Exchange) ThrowBroadlyMatchedException(broad interface{}, s interface{}, message interface{}) {
	broadKey := self.FindBroadlyMatchedKey(broad, s)
	if broadKey != "" {
		self.RaiseException(broad.(map[string]string)[broadKey], message)
	}
}

func (self *Exchange) PanicToError(e interface{}) (err error) {
	switch e.(type) {
	case []string:
		args := e.([]string)
		if len(args) == 2 {
			errCls := args[0]
			message := args[1]
			//err = errors.New(errCls + ": " + message)
			err = TypedError(errCls, message)
		} else {
			if self.Verbose {
				log.Println(string(debug.Stack()))
			}
			err = fmt.Errorf("Catch unknown panic: %v", e)
		}
	default:
		if self.Verbose {
			log.Println(string(debug.Stack()))
		}
		err = fmt.Errorf("Catch unknown panic: %v", e)
	}
	return
}

func (self *Exchange) HandleRestErrors(httpStatusCode int, httpStatusText string, body string, url string, method string) {
	errCls := ""
	strCode := strconv.Itoa(httpStatusCode)
	if _, ok := self.httpExceptions[strCode]; ok {
		errCls = self.httpExceptions[strCode]
		if errCls == "ExchangeNotAvailable" {
			matched, err := regexp.MatchString("(?i)(cloudflare|incapsula|overload|ddos)", body)
			if matched && err == nil {
				errCls = "DDoSProtection"
			}
		}
	}
	if errCls != "" {
		self.RaiseException(errCls, strings.Join([]string{method, url, strCode, httpStatusText, body}, " "))
	}
}

func (self *Exchange) IsJsonEncodedObject(input interface{}) bool {
	strInput, ok := input.(string)
	if ok {
		if len(strInput) >= 2 {
			if strInput[0] == '{' || strInput[0] == '[' {
				return true
			}
		}
	}
	return false
}

func (self *Exchange) HandleRestResponse(response string, jsonResponse interface{}, url string, method string) {
	if self.IsJsonEncodedObject(response) && jsonResponse == nil {
		dDoSProtectionMatched, _ := regexp.MatchString("(?i)(cloudflare|incapsula|overload|ddos)", response)
		if dDoSProtectionMatched {
			self.RaiseException("DDoSProtection", strings.Join([]string{method, url, response}, " "))
		}
		exchangeNotAvailableMatched, _ := regexp.MatchString("(?i)(offline|busy|retry|wait|unavailable|maintain|maintenance|maintenancing)", response)
		if exchangeNotAvailableMatched {
			message := response + " exchange downtime, exchange closed for maintenance or offline, DDoS protection or rate-limiting in effect"
			self.RaiseException("ExchangeNotAvailable", strings.Join([]string{method, url, response, message}, " "))
		}
		self.RaiseException("ExchangeError", strings.Join([]string{method, url, response}, " "))
	}
}

func (self *Exchange) Float64ToString(f float64) string {
	//return fmt.Sprintf("%v", f)
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// 如果是 map 就使用值转为 slice
func (self *Exchange) Values(x interface{}) []interface{} {
	if v, ok := x.([]interface{}); ok {
		return v
	} else if v, ok := x.(map[string]interface{}); ok {
		out := make([]interface{}, 0, len(v))
		keys := make([]string, 0, len(v))

		for k, _ := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			out = append(out, v[k])
		}
		return out
	}
	self.RaiseException("InternalError", "Values error: "+fmt.Sprint(x))
	return nil
}

/*
   Accepts a map/array of objects and a key name to be used as an index:
   array = [
      { someKey: 'value1', anotherKey: 'anotherValue1' },
      { someKey: 'value2', anotherKey: 'anotherValue2' },
      { someKey: 'value3', anotherKey: 'anotherValue3' },
   ]
   key = 'someKey'

   Returns a map:
  {
      value1: { someKey: 'value1', anotherKey: 'anotherValue1' },
      value2: { someKey: 'value2', anotherKey: 'anotherValue2' },
      value3: { someKey: 'value3', anotherKey: 'anotherValue3' },
  }
*/
func (self *Exchange) IndexBy(x interface{}, k string) map[string]interface{} {
	out := map[string]interface{}{}
	for _, v := range self.Values(x) {
		m := v.(map[string]interface{})
		if _, ok := m[k]; ok {
			out[fmt.Sprintf("%v", m[k])] = v
		}
	}
	return out
}

func (self *Exchange) InArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (self *Exchange) FetchAccounts(params map[string]interface{}) []interface{} {
	return nil
}

func (self *Exchange) ToArray(o interface{}) (result []interface{}) {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Map:
		for k, v := range o.([]interface{}) {
			result = append(result, []interface{}{k, v})
		}
	case reflect.Slice:
		result = o.([]interface{})
	default:
		self.RaiseInternalException("unsupport type for ToArray!")
	}
	return
}

func (self *Exchange) ArrayConcat(a interface{}, b interface{}) (result []interface{}) {
	return append(a.([]interface{}), b.([]interface{})...)
}

func (self *Exchange) FilterByValueSinceLimit(arr []interface{}, field string, value interface{}, since interface{}, limit interface{}, key string, tail bool) (result []interface{}) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("filter_by_symbol_since_limit err: %v, arr: %v\n", e, arr)
		}
	}()

	result = self.ToArray(arr)

	if value != nil {
		result = funk.Filter(result, func(x interface{}) bool {
			return x.(map[string]interface{})[field] == value
		}).([]interface{})
	}

	if since != nil {
		result = funk.Filter(result, func(x interface{}) bool {
			return x.(map[string]interface{})[key].(int64) >= since.(int64)
		}).([]interface{})
	}

	if limit != nil {
		limitNum := limit.(int64)
		lenNum := int64(len(result))
		if limitNum > lenNum {
			limitNum = lenNum
		}
		if tail && since != nil {
			result = result[lenNum-limitNum:]
		} else {
			result = result[:limitNum]
		}
	}
	return
}
func (self *Exchange) FilterBySymbolSinceLimit(arr []interface{}, symbol interface{}, since interface{}, limit interface{}) (result []interface{}) {
	return self.FilterByValueSinceLimit(arr, "symbol", symbol, since, limit, "timestamp", false)
}

func (self *Exchange) DeepExtend(args ...interface{}) (result map[string]interface{}) {
	for _, arg := range args {
		err := mergo.Merge(&result, arg, mergo.WithOverride)
		if err != nil {
			self.RaiseInternalException(fmt.Sprintf("deepExtend err:%v, args:%v", err, args))
		}
	}
	return
}

func (self *Exchange) InitDescribe() (err error) {
	//err = json.Unmarshal(self.Child.Describe(), &self.DescribeMap)
	err = jsonx.Unmarshal(self.Child.Describe(), &self.DescribeMap, jsonx.WithExtraComma(), jsonx.WithComment())
	if err != nil {
		return
	}

	self.DescribeJson = gjson.ParseBytes(self.Child.Describe())

	err = self.DefineRestApi()
	if err != nil {
		return
	}

	self.Options = self.DescribeMap["options"].(map[string]interface{})
	self.Timeframes = map[string]string{}
	for key, val := range self.DescribeMap["timeframes"].(map[string]interface{}) {
		self.Timeframes[key] = val.(string)
	}
	self.Urls = self.DescribeMap["urls"].(map[string]interface{})
	if self.DescribeMap["version"] != nil {
		self.Version = self.DescribeMap["version"].(string)
	}
	self.Exceptions = self.DescribeMap["exceptions"].(map[string]interface{})
	if hostName, ok := self.DescribeMap["hostname"]; ok {
		self.Hostname = hostName.(string)
	}
	if fees, ok := self.DescribeMap["fees"]; ok {
		self.Fees = fees.(map[string]interface{})
	}
	self.CommonCurrencies = map[string]string{
		"XBT":    "BTC",
		"BCC":    "BCH",
		"DRK":    "DASH",
		"BCHABC": "BCH",
		"BCHSV":  "BSV",
	}

	return
}

func (self *Exchange) SetBaseUrl(u string) {
	if m, ok := self.Urls["api"].(map[string]interface{}); ok {
		for key := range m {
			m[key] = u
		}
	} else {
		self.Urls["api"] = u
	}
}

// key: public / private
func (self *Exchange) BaseUrl(key string) string {
	if self.Urls["api"] == nil {
		return ""
	}
	if m, ok := self.Urls["api"].(map[string]interface{}); ok {
		if m[key] == nil {
			return ""
		} else if u, ok := m[key].(string); ok {
			return u
		}
	}
	if u, ok := self.Urls["api"].(string); ok {
		return u
	}
	return ""
}

func (self *Exchange) SetVerbose(verbose bool) {
	self.Verbose = verbose
}

func (self *Exchange) Ymdhms(m int64, t string) string {
	unixTimeUTC := time.Unix(m/1000, 0).In(time.UTC) // gives unix time stamp in utc
	return unixTimeUTC.Format(fmt.Sprintf("2006-01-02%v15:04:05", t))
}

func (self *Exchange) CommonCurrencyCode(currency string) string {
	//return self.SafeString(self.CommonCurrencies, currency, currency)
	// NOTE: 我们不需要 CommonCurrencyCode 功能
	return currency
}

func (self *Exchange) ParseNumber(v interface{}) float64 {
	def := 0.0
	if v == nil {
		return def
	}
	f, ok := v.(float64)
	if ok {
		return f
	}
	return def
}

func (self *Exchange) SafeSymbol(marketId string, market interface{}, delimiter string) (symbol string) {
	if marketId != "" {
		if self.ToBool(self.InMap(marketId, self.MarketsById)) {
			market = self.Member(self.MarketsById, marketId)
		} else {
			baseId, quoteId := self.Unpack2(strings.Split(marketId, "/"))
			base := self.SafeCurrencyCode(baseId)
			quote := self.SafeCurrencyCode(quoteId)
			symbol = base + delimiter + quote
		}
	}
	if self.ToBool(self.TestNil(symbol) && !self.TestNil(market)) {
		symbol = self.Member(market, "symbol").(string)
	}
	return
}
