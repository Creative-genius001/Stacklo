package services

import (
	"context"

	"github.com/Creative-genius001/Stacklo/services/payment/pkg/binance"
	"github.com/Creative-genius001/Stacklo/services/payment/pkg/paystack"
)

type PaymentService interface {
	GetBankList() (*paystack.Banks, error)
	GetOTP(payload *paystack.TransferOtpRequest) (*paystack.TransferOtpResponse, error)
	ResolveAccountNumber(accountNumber string, bankCode string) (*paystack.AccountResolutionResponse, error)
	CreateTransferRecipient(payload *paystack.CreateTransferRecipientRequest) (*paystack.CreateTransferRecipientResponse, error)
	Transfer(payload paystack.FianlTransferRequest) (*paystack.FinalTransferResponse, error)
	Ping() error
	PlaceBuyOrder(ctx context.Context, req binance.BinanceOrderRequest) (*binance.BinanceOrderResponse, error)
	Convert(creq binance.ConvertAssetRequest) ([]*binance.ConvertAssetResponse, error)
}

type paymentServiceImpl struct {
	paystackClient paystack.Paystack
	binanceClient  binance.Binance
}

func NewPaymentService(b binance.Binance, p paystack.Paystack) PaymentService {
	return &paymentServiceImpl{p, b}
}

func (p *paymentServiceImpl) GetBankList() (*paystack.Banks, error) {

	res, err := p.paystackClient.GetBankList()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) ResolveAccountNumber(accountNumber string, bankCode string) (*paystack.AccountResolutionResponse, error) {
	res, err := p.paystackClient.ResolveAccountNumber(accountNumber, bankCode)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) CreateTransferRecipient(transferRecipient *paystack.CreateTransferRecipientRequest) (*paystack.CreateTransferRecipientResponse, error) {
	res, err := p.paystackClient.CreateTransferRecipient(transferRecipient)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) GetOTP(payload *paystack.TransferOtpRequest) (*paystack.TransferOtpResponse, error) {
	res, err := p.paystackClient.GetOTP(payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) Transfer(payload paystack.FianlTransferRequest) (*paystack.FinalTransferResponse, error) {
	res, err := p.paystackClient.Transfer(payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) Ping() error {
	err := p.binanceClient.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (p *paymentServiceImpl) PlaceBuyOrder(ctx context.Context, req binance.BinanceOrderRequest) (*binance.BinanceOrderResponse, error) {
	panic("unimplemented")
}

func (p *paymentServiceImpl) Convert(creq binance.ConvertAssetRequest) ([]*binance.ConvertAssetResponse, error) {
	res, err := p.binanceClient.Convert(creq)
	if err != nil {
		return nil, err
	}
	return res, nil
}
