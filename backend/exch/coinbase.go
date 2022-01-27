package exch

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
	"os"
	"strings"
	"time"
)

var AllowedCurrencies map[string]bool = map[string]bool{"EUR": true}

type Coinbase struct {
	client *coinbasepro.Client
	hedges map[string]decimal.Decimal
}

var _ HedgeExchange = &Coinbase{}

func NewCoinbase() (*Coinbase, error) {
	godotenv.Load()

	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		Key:        os.Getenv("CB_API_KEY"),
		Passphrase: os.Getenv("CB_API_PHRASE"),
		Secret:     os.Getenv("CB_API_SECRET"),
	})

	_, err := client.GetBook("BTC-USD", 1)
	if err != nil {
		return nil, err
	}
	return &Coinbase{
		client: client,
		hedges: map[string]decimal.Decimal{},
	}, nil
}

func (c *Coinbase) PricePerBTC(currency string) (decimal.Decimal, error) {
	bal, err := c.client.GetBook(fmt.Sprintf("BTC-%s", strings.ToUpper(currency)), 1)
	if err != nil {
		return decimal.Decimal{}, err
	}

	// Use asking price, add a fee on top
	price, err := decimal.NewFromString(bal.Asks[0].Price)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return price, nil
}

func (c *Coinbase) Available(currency string) (decimal.Decimal, error) {
	acc, err := c.client.GetAccounts()
	if err != nil {
		return decimal.Decimal{}, err
	}
	for _, acc := range acc {
		if acc.Currency == currency {
			return decimal.NewFromString(acc.Available)
		}
	}
	return decimal.Decimal{}, err
}

func (c *Coinbase) AddHedgeAmount(msat decimal.Decimal, currency string) error {
	if AllowedCurrencies[currency] {
		return errors.New("Currency not supported")
	}
	c.hedges[currency] = c.hedges[currency].Add(msat)
	return nil
}

func (c *Coinbase) FinishHedge() error {
	for curr, amt := range c.hedges {
		order, err := c.client.CreateOrder(&coinbasepro.Order{
			Type:      "market",
			Side:      "buy",
			Size:      amt.Div(MsatToBtc).String(),
			ProductID: fmt.Sprintf("BTC-%s", strings.ToUpper(curr)),
		})

		if err != nil {
			return errors.Wrap(err, "Could not create order")
		}
		for i := 1; i < 5; i++ {
			order, err = c.client.GetOrder(order.ID)
			if err != nil {
				if err2 := c.client.CancelOrder(order.ID); err2 != nil {
					return err2
				}
				return err
			}
			if order.Settled {
				println("Order settled", &order)
				return nil
			}
			time.Sleep(time.Millisecond * 500)
		}
		if err2 := c.client.CancelOrder(order.ID); err2 != nil {
			return err2
		}
		return errors.New("Could not complete tx in time")
	}
	return nil
}
