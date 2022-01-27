package exch

import (
	"github.com/shopspring/decimal"
)

var BtcToMsat, _ = decimal.NewFromString("0.00000000001")
var MsatToBtc = decimal.NewFromInt(1).Div(BtcToMsat)

type HedgeExchange interface {
	PricePerBTC(currency string) (decimal.Decimal, error)

	Available(currency string) (decimal.Decimal, error)
	AddHedgeAmount(msat decimal.Decimal, currency string) error

	FinishHedge() error
}
