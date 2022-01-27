package db

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&Invoice{})
	if err != nil {
		panic(err)
	}
}

type Invoice struct {
	ID      uuid.UUID `json:"id" sql:"uuid;primary_key" gorm:"primary_key"`
	Invoice string    `json:"invoice"`

	ChargeAmount   decimal.Decimal `json:"charge_amount"`
	ChargeCurrency string          `json:"charge_currency"`

	ChargeID *string `json:"charge_id"`
	RefundID *string `json:"refund_id"`
	PreImage *string `json:"preimage"`
}
