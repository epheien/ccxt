package ccxt

import (
	"fmt"
	"github.com/epheien/ccxt/go/ascendex"
	"github.com/epheien/ccxt/go/base"
	"github.com/epheien/ccxt/go/binance"
	"github.com/epheien/ccxt/go/bitmax"
	"github.com/epheien/ccxt/go/bitmax2"
	"github.com/epheien/ccxt/go/bybit"
	"github.com/epheien/ccxt/go/futures_binance"
	"github.com/epheien/ccxt/go/futures_gateio"
	"github.com/epheien/ccxt/go/futures_kucoin"
	"github.com/epheien/ccxt/go/gateio"
	"github.com/epheien/ccxt/go/kucoin"
	"github.com/epheien/ccxt/go/kucoin_hf"
	"github.com/epheien/ccxt/go/margin_bitmax"
	"github.com/epheien/ccxt/go/margin_kucoin"
	"github.com/epheien/ccxt/go/mexc"
)

type IExchange = base.ExchangeInterface
type ExchangeConfig = base.ExchangeConfig
type Order = base.Order
type Position = base.Position

func New(exchange string, config *base.ExchangeConfig) (ex IExchange, err error) {
	switch exchange {
	case "binance":
		ex, err = binance.New(config)
	case "bybit":
		ex, err = bybit.New(config)
	case "kucoin":
		ex, err = kucoin.New(config)
	case "kucoin_hf":
		ex, err = kucoin_hf.New(config)
	case "bitmax":
		ex, err = bitmax.New(config)
	case "bitmax2":
		ex, err = bitmax2.New(config)
	case "ascendex":
		ex, err = ascendex.New(config)
	case "margin_bitmax":
		ex, err = margin_bitmax.New(config)
	case "margin_kucoin":
		ex, err = margin_kucoin.New(config)
	case "gateio", "gateio4":
		ex, err = gateio.New(config)
	case "futures_kucoin", "futures_kumex":
		ex, err = futures_kucoin.New(config)
	case "futures_binance":
		ex, err = futures_binance.New(config)
	case "futures_gateio":
		ex, err = futures_gateio.New(config)
	case "mexc":
		ex, err = mexc.New(config)
	default:
		err = fmt.Errorf("exchange %s is not supported", exchange)
	}
	return
}
