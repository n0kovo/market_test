package marketplace

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/n0kovo/market_test/modules/apis"
)

/*
	Models
*/

type PaymentReceipt struct {
	Uuid           string `json:"uuid" gorm:"primary_key" sql:"size:1024"`
	Type           string
	SerializedData string

	BTCPaymentResultItem  apis.BTCPaymentResult
	BCHPaymentResultItem  apis.BCHPaymentResult
	ETHPaymentResultItems []apis.ETHPaymentResult
}

/*
	Model Interface Implementation
*/

func (r PaymentReceipt) Validate() error {
	if r.Uuid == "" {
		return errors.New("Uuid can't be empty")
	}
	if r.Type == "" {
		return errors.New("Type can not be empty")
	}
	if r.SerializedData == "" {
		return errors.New("Serialized data can not be empty")
	}
	return nil
}

/*
	Database Methods
*/

func (i PaymentReceipt) Remove() error {
	return database.Delete(&i).Error
}

func (itm PaymentReceipt) Save() error {
	if err := itm.Validate(); err != nil {
		return err
	}
	return itm.SaveToDatabase()
}

func (itm PaymentReceipt) SaveToDatabase() error {
	return database.Create(&itm).Error
}

/*
	Cryptocurrency-Specific Methods
*/

func (i PaymentReceipt) BTCPaymentResult() (apis.BTCPaymentResult, error) {
	pr := apis.BTCPaymentResult{}
	err := json.Unmarshal([]byte(i.SerializedData), &pr)
	return pr, err
}

func (i PaymentReceipt) BCHPaymentResult() (apis.BCHPaymentResult, error) {
	pr := apis.BCHPaymentResult{}
	err := json.Unmarshal([]byte(i.SerializedData), &pr)
	return pr, err
}

func (i PaymentReceipt) ETHPaymentResults() ([]apis.ETHPaymentResult, error) {
	pr := []apis.ETHPaymentResult{}
	err := json.Unmarshal([]byte(i.SerializedData), &pr)
	return pr, err
}

/*
	Factory Methods
*/

func CreateBTCPaymentReceipt(transactionResult apis.BTCPaymentResult) (PaymentReceipt, error) {
	txResultJSON, err := json.Marshal(transactionResult)
	if err != nil {
		return PaymentReceipt{}, err
	}
	receipt := PaymentReceipt{
		Uuid:           transactionResult.Hash,
		Type:           "bitcoin",
		SerializedData: string(txResultJSON),
	}
	return receipt, receipt.Save()
}

func CreateBCHPaymentReceipt(transactionResult apis.BCHPaymentResult) (PaymentReceipt, error) {
	txResultJSON, err := json.Marshal(transactionResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt := PaymentReceipt{
		Uuid:           transactionResult.Hash,
		Type:           "bitcoin_cash",
		SerializedData: string(txResultJSON),
	}
	return receipt, receipt.Save()
}

func CreateETHPaymentReceipt(transactionResult []apis.ETHPaymentResult) (PaymentReceipt, error) {
	txResultJSON, err := json.Marshal(transactionResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	h := md5.New()
	for _, pr := range transactionResult {
		io.WriteString(h, pr.Hash)
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	receipt := PaymentReceipt{
		Uuid:           hash,
		Type:           "ethereum",
		SerializedData: string(txResultJSON),
	}
	return receipt, receipt.Save()
}

/*
	Queries
*/

func FindPaymentReceiptByUuid(uuid string) []ReferralPayment {
	var payments []ReferralPayment

	database.
		Where("uuid=?", uuid).
		Find(&payments)

	return payments
}
