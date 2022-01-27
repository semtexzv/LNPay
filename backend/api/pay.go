package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"lnpay/db"
	"lnpay/exch"
	"lnpay/lnet"
	"net/http"
	"os"
	"strings"
	"time"

	decodepay "github.com/fiatjaf/ln-decodepay"
)

const MinAmount = 5
const MaxAmount = 500

var FeeRatio = decimal.NewFromFloat(1.05)
var FiatCurrency string
var Exchange exch.HedgeExchange

func init() {
	godotenv.Load()
	FiatCurrency = os.Getenv("FIAT_CURRENCY")
	if FiatCurrency == "" {
		FiatCurrency = "EUR"
	}
	var err error
	Exchange, err = exch.NewCoinbase()
	if err != nil {
		panic(err)
	}
}

func CreateCharge(amount decimal.Decimal, token string, currency string) (string, error) {
	c, err := charge.New(&stripe.ChargeParams{
		Amount:   stripe.Int64(amount.Mul(decimal.NewFromInt(100)).IntPart()),
		Currency: stripe.String(strings.ToLower(currency)),
		Source:   &stripe.SourceParams{Token: stripe.String(token)},
	})

	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func RefundCharge(chargeId string) (string, error) {
	r, err := refund.New(&stripe.RefundParams{
		Charge: stripe.String(chargeId),
	})
	if err != nil {
		return "", err
	}
	return r.ID, nil
}

func CreatePayment(c *gin.Context) {
	invoiceStr, hasInvoice := c.GetQuery("invoice")
	if !hasInvoice {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"missing invoice param"})
		return
	}

	bolt11, err := decodepay.Decodepay(invoiceStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Invalid invoice "})
		return
	}

	if os.Getenv("LN_CURRENCY") != bolt11.Currency {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Invalid LN currency"})
		return
	}

	if time.Now().Unix() > int64(bolt11.CreatedAt)+int64(bolt11.Expiry) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"LN Invoice timed out"})
		return
	}

	price, err := Exchange.PricePerBTC(FiatCurrency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not get price"})
		return
	}
	effectivePrice := price.Mul(FeeRatio)

	amount := effectivePrice.Mul(decimal.NewFromInt(bolt11.MSatoshi)).Div(exch.MsatToBtc)
	println(amount.String())
	if amount.LessThan(decimal.NewFromInt(MinAmount)) || amount.GreaterThan(decimal.NewFromInt(MaxAmount)) {
		resp := ErrorResp{fmt.Sprintf("We only support payments between %v and %v %v", MinAmount, MaxAmount, FiatCurrency)}
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	invoice := db.Invoice{
		ID:             uuid.New(),
		Invoice:        invoiceStr,
		ChargeAmount:   amount,
		ChargeCurrency: FiatCurrency,
	}
	db.DB.Create(&invoice)

	c.JSON(http.StatusOK, &struct {
		ID       string          `json:"id"`
		Invoice  string          `json:"invoice"`
		Amount   decimal.Decimal `json:"amount"`
		Currency string          `json:"currency"`
	}{
		ID:       invoice.ID.String(),
		Invoice:  invoiceStr,
		Amount:   amount,
		Currency: FiatCurrency,
	})
}

type FinishPaymentIn struct {
	Token string `json:"token"`
}

type FinishPaymentRes struct {
	PaymentID string `json:"payment_id"`
	ChargeID  string `json:"charge_id"`
	RefundID  string `json:"refund_id,omitempty"`
	Preimage  string `json:"preimage"`
	Error     string `json:"error,omitempty"`
}

func FinishPayment(c *gin.Context) {
	var err error
	paymentID := c.Param("payment_id")

	if paymentID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"missing payment id param"})
		return
	}

	var input FinishPaymentIn
	if err := c.ShouldBindJSON(&input); err != nil || input.Token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Invalid body provided"})
		return
	}

	var invoice db.Invoice
	if err = db.DB.Find(&invoice, "id = ?", paymentID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResp{"Did not find the invoice"})
		return
	}

	var bolt11 decodepay.Bolt11
	if bolt11, err = decodepay.Decodepay(invoice.Invoice); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Invalid invoice "})
		return
	}

	price, err := Exchange.PricePerBTC(FiatCurrency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not get price"})
		return
	}
	originalAmount := price.Mul(decimal.NewFromInt(bolt11.MSatoshi))
	amount := originalAmount.Mul(FeeRatio)

	var avail decimal.Decimal
	if avail, err = Exchange.Available(FiatCurrency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not query the exchange"})
		return
	}

	if avail.LessThan(originalAmount.Mul(decimal.NewFromInt(5))) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Not enough liquidity"})
		return
	}

	chargeID, err := CreateCharge(amount, input.Token, FiatCurrency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not create stripe charge"})
		return
	}

	invoice.ChargeID = &chargeID
	if err = db.DB.Save(&invoice).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResp{"Could not save stripe charge into the database"})
		return
	}

	result := FinishPaymentRes{
		PaymentID: invoice.ID.String(),
		ChargeID:  chargeID,
	}

	var refundID string
	if err = Exchange.AddHedgeAmount(originalAmount, FiatCurrency); err != nil {
		result.Error = err.Error()

		if refundID, err = RefundCharge(chargeID); err != nil {
			panic(err)
		}

		if err = Exchange.AddHedgeAmount(originalAmount.Neg(), FiatCurrency); err != nil {
			panic(err)
		}

		result.RefundID = refundID
		invoice.RefundID = &refundID
		db.DB.Save(&invoice)

		c.AbortWithStatusJSON(http.StatusFailedDependency, &result)
		return
	}

	preimage, err := lnet.SendPayment(invoice)
	if err != nil {
		result.Error = fmt.Sprintf("LN Error: %s", err.Error())

		if refundID, err = RefundCharge(chargeID); err != nil {
			panic(err)
		}
		if err = Exchange.AddHedgeAmount(originalAmount.Neg(), FiatCurrency); err != nil {
			panic(err)
		}

		result.RefundID = refundID
		invoice.RefundID = &refundID
		db.DB.Save(&invoice)

		c.AbortWithStatusJSON(http.StatusFailedDependency, &result)
		return
	}

	if err = Exchange.FinishHedge(); err != nil {
		panic("Failed to finish Hedge")
	}

	invoice.PreImage = &preimage
	result.Preimage = preimage
	db.DB.Save(&invoice)
	c.JSON(http.StatusOK, &result)
}
