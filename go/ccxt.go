package ccxt

import (
	"fmt"
	"github.com/georgexdz/ccxt/go/base"
	"github.com/georgexdz/ccxt/go/kucoin"
	"github.com/georgexdz/ccxt/go/bitmax"
	"github.com/georgexdz/ccxt/go/margin_bitmax"
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
	case "margin_bitmax":
		ex, err = margin_bitmax.New(config)
	default:
		err = fmt.Errorf("exchange %s is not supported", exchange)
	}
	return
}
