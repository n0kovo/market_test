package apis

import (
	"encoding/json"
	"fmt"
	"net/url"
	"github.com/n0kovo/market_test/modules/util"
)

type ETHPayment struct {
	Address string  `json:"address"`
	Percent float64 `json:"percent,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

type ETHWalletBalance struct {
	Balance float64 `json:"balance"`
}

type ETHPaymentResult struct {
	Hash       string  `json:"hash"`
	RawTx      string  `json:"rawTransaction"`
	Cost       int     `json:"cost"`
	WalletFrom string  `json:"wallet_from"`
	WalletTo   string  `json:"wallet_to"`
	Amount     float64 `json:"amount"`
}

func GenerateETHAddress(walletType string) (string, error) {
	apiEndpoint := fmt.Sprintf("%s/ethereum/wallets/new", APPLICATION_SETTINGS.PaymentGate)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"type": {walletType}})
	if err != nil {
		return "", err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		return "", err
	}

	address := dat["address"].(string)
	return address, nil
}

func GetAmountOnETHAddress(address string) (ETHWalletBalance, error) {
	apiEndpoint := fmt.Sprintf("%s/ethereum/wallets/%s", APPLICATION_SETTINGS.PaymentGate, address)

	walletBalance := ETHWalletBalance{}

	body, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return walletBalance, err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		return walletBalance, err
	}

	var balance float64
	if dat["balance"] != nil {
		balance = dat["balance"].(float64)
	} else {
		balance = 0.0
	}

	walletBalance.Balance = balance

	return walletBalance, nil
}

func SendETHFromSingleWallet(addressFrom string, btcPayments []ETHPayment) ([]ETHPaymentResult, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/ethereum/wallets/%s/send", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result      = []ETHPaymentResult{}
	)

	paymentsJSON, err := json.Marshal(btcPayments)
	if err != nil {
		return result, err
	}

	response, err := util.DirectPOST(
		apiEndpoint,
		url.Values{
			"payments": {string(paymentsJSON)},
		},
	)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}
