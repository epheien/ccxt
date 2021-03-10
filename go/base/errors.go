package base

import (
	"errors"
	"fmt"
)

// NOTE: go 无法实现类似 ccxt 的层次式的错误架构
// var (
// 	BaseError = errors.New("")
// 
// 	InternalError = fmt.Errorf("%w", BaseError)
// 
// 	ExchangeError = fmt.Errorf("%w", BaseError)
// 
// 	AuthenticationError = fmt.Errorf("%w", ExchangeError)
// 	PermissionDenied = fmt.Errorf("%w", AuthenticationError)
// 	AccountSuspended = fmt.Errorf("%w", AuthenticationError)
// 
// 	ArgumentsRequired = fmt.Errorf("%w", ExchangeError)
// 	BadRequest = fmt.Errorf("%w", ExchangeError)
// 	BadSymbol = fmt.Errorf("%w", BadRequest)
// 
// 	BadResponse = fmt.Errorf("%w", ExchangeError)
// 	NullResponse = fmt.Errorf("%w", BadResponse)
// 
// 	InsufficientFunds = fmt.Errorf("%w", ExchangeError)
// 	InvalidAddress = fmt.Errorf("%w", ExchangeError)
// 	AddressPending = fmt.Errorf("%w", InvalidAddress)
// 
// 	InvalidOrder = fmt.Errorf("%w", ExchangeError)
// 	OrderNotFound = fmt.Errorf("%w", InvalidOrder)
// 	OrderNotCached = fmt.Errorf("%w", InvalidOrder)
// 	CancelPending = fmt.Errorf("%w", InvalidOrder)
// 	OrderImmediatelyFillable = fmt.Errorf("%w", InvalidOrder)
// 	OrderNotFillable = fmt.Errorf("%w", InvalidOrder)
// 	DuplicateOrderId = fmt.Errorf("%w", InvalidOrder)
// 
// 	NotSupported = fmt.Errorf("%w", ExchangeError)
// 
// 	NetworkError = fmt.Errorf("%w", BaseError)
// 	DDoSProtection = fmt.Errorf("%w", NetworkError)
// 	RateLimitExceeded = fmt.Errorf("%w", DDoSProtection)
// 	ExchangeNotAvailable = fmt.Errorf("%w", NetworkError)
// 	OnMaintenance = fmt.Errorf("%w", ExchangeNotAvailable)
// 	InvalidNonce = fmt.Errorf("%w", NetworkError)
// 	RequestTimeout = fmt.Errorf("%w", NetworkError)
// )

// 用最简单的方式实现即可
func TypedError(t string, msg string) error {
	var err error
	err = errors.New(t)
	return fmt.Errorf("%w: %v", err, msg)
}

