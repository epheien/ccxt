package ccxt

import (
	"fmt"
	"github.com/georgexdz/ccxt/go/ascendex"
	"github.com/georgexdz/ccxt/go/base"
	"github.com/georgexdz/ccxt/go/bitmax"
	"github.com/georgexdz/ccxt/go/bitmax2"
	"github.com/georgexdz/ccxt/go/kucoin"
	"github.com/georgexdz/ccxt/go/margin_bitmax"
	"github.com/georgexdz/ccxt/go/margin_kucoin"
)

type IExchange = base.ExchangeInterface
type ExchangeConfig = base.ExchangeConfig
type Order = base.Order

func New(exchange string, config *base.ExchangeConfig) (ex IExchange, err error) {
	switch exchange {
	case "kucoin":
		ex, err = kucoin.New(config)
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
	default:
		err = fmt.Errorf("exchange %s is not supported", exchange)
	}
	return
}
