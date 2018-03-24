package marketplace

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type ViewTransaction struct {
	Transaction
	Seller                User
	Buyer                 User
	RatingReview          RatingReview
	Amount                string
	AmountToPay           string
	CurrentAmountPaid     string
	CompletedAt           string
	ConfirmedAt           string
	CreatedAt             string
	CurrentPaymentStatus  string
	CurrentShippingStatus string
	NextStatusTime        string
	NextStatusPercent     float64

	FEAllowed     bool
	CancelAllowed bool
	IsCompleted   bool
	IsReleased    bool
	IsFrozen      bool
	IsCancelled   bool
	IsPending     bool
	IsFailed      bool
	IsDispatched  bool
	IsShipped     bool
	IsReviewed    bool
	IsDisputed    bool

	NumberOfMessages int

	ViewTransactionStatus []ViewTransactionStatus
	ViewShippingStatus    []ViewShippingStatus
}

func (t Transaction) ViewTransaction() ViewTransaction {

	buyer := t.Buyer
	seller := t.Seller

	completedAt := t.CompletedAt()
	completedAtStr := ""
	if completedAt != nil {
		completedAtStr = completedAt.Format("02.01.2006 15:04")
	}

	vt := ViewTransaction{
		Seller:                seller,
		Transaction:           t,
		Buyer:                 buyer,
		CompletedAt:           completedAtStr,
		CreatedAt:             t.CreatedAt().Format("02.01.2006 15:04"),
		CurrentAmountPaid:     humanize.Ftoa(t.CurrentAmountPaid()),
		CurrentPaymentStatus:  t.CurrentPaymentStatus(),
		CurrentShippingStatus: t.CurrentShippingStatus(),
		FEAllowed:             seller.Premium,
		CancelAllowed:         t.IsCompleted() && !t.IsDispatched() && !t.IsShipped() && t.Package.Type != "digital" && t.Package.Type != "drop",

		NumberOfMessages: t.NumberOfMessages(),

		IsCancelled:  t.IsCancelled(),
		IsCompleted:  t.IsCompleted(),
		IsDispatched: t.IsDispatched(),
		IsDisputed:   t.DisputeUuid != "",
		IsFailed:     t.IsFailed(),
		IsFrozen:     t.IsFrozen(),
		IsPending:    t.IsPending(),
		IsReleased:   t.IsReleased(),
		IsShipped:    t.IsShipped(),
	}

	switch t.Type {
	case "bitcoin":
		vt.Amount = humanize.Ftoa(t.BitcoinTransaction.Amount)
		vt.AmountToPay = humanize.Ftoa(t.BitcoinTransaction.Amount - t.CurrentAmountPaid())
	case "bitcoin_cash":
		vt.Amount = humanize.Ftoa(t.BitcoinCashTransaction.Amount)
		vt.AmountToPay = humanize.Ftoa(t.BitcoinCashTransaction.Amount - t.CurrentAmountPaid())
	case "ethereum":
		vt.Amount = humanize.Ftoa(t.EthereumTransaction.Amount)
		vt.AmountToPay = humanize.Ftoa(t.EthereumTransaction.Amount - t.CurrentAmountPaid())
	}

	review, _ := FindRatingReviewByTransactionUuid(t.Uuid)
	if review != nil {
		vt.RatingReview = *review
		vt.IsReviewed = true
	}

	vtss := []ViewTransactionStatus{}
	for _, ts := range t.Status {
		vts := ViewTransactionStatus{
			Amount:  ts.Amount,
			Time:    humanize.Time(ts.Time),
			Comment: ts.Comment,
			Status:  ts.Status,
		}
		if ts.PaymentReceipt.Uuid != "" {
			switch t.Type {
			case "bitcoin":
				pr, err := ts.PaymentReceipt.BTCPaymentResult()
				if err == nil {
					vts.BTCPaymentResult = pr
				}
			case "bitcoin_cash":
				pr, err := ts.PaymentReceipt.BCHPaymentResult()
				if err == nil {
					vts.BCHPaymentResult = pr
				}
			case "ethereum":
				pr, err := ts.PaymentReceipt.ETHPaymentResults()
				if err == nil {
					vts.ETHPaymentResults = pr
				}
			}
		}
		vtss = append(vtss, vts)
	}
	vt.ViewTransactionStatus = vtss

	vsss := []ViewShippingStatus{}
	for _, ts := range t.ShippingStatus {
		vts := ViewShippingStatus{
			Time:    humanize.Time(ts.Time),
			Comment: ts.Comment,
			Status:  ts.Status,
		}
		vsss = append(vsss, vts)
	}
	vt.ViewShippingStatus = vsss

	if vt.IsPending {
		now := time.Now()
		minutesLeft := int(t.CreatedAt().Add(pendingDuration).Sub(now).Minutes())
		vt.NextStatusTime = fmt.Sprintf("%d minutes", minutesLeft)
		vt.NextStatusPercent = float64(int(float64(minutesLeft) / pendingDuration.Minutes() * 100))
	}

	if vt.IsCompleted {
		now := time.Now()
		minutesLeft := int(t.CreatedAt().Add(completedDuration).Sub(now).Hours())
		vt.NextStatusTime = fmt.Sprintf("%d hours", int(minutesLeft))
		vt.NextStatusPercent = float64(int(float64(minutesLeft) / completedDuration.Minutes() * 100))
	}

	return vt
}

func (ts Transactions) ViewTransactions() []ViewTransaction {
	vts := []ViewTransaction{}
	for _, t := range ts {
		vts = append(vts, t.ViewTransaction())
	}
	return vts
}
