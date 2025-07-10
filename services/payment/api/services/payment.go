package services

import (
	"github.com/Creative-genius001/Stacklo/services/payment/pkg/binance"
	"github.com/Creative-genius001/Stacklo/services/payment/pkg/paystack"
	"github.com/Creative-genius001/Stacklo/services/payment/types"
)

type PaymentService interface {
	GetBankList() (*types.Banks, error)
	GetOTP(payload *types.TransferOtpRequest) (*types.TransferOtpResponse, error)
	ResolveAccountNumber(accountNumber string, bankCode string) (*types.AccountResolutionResponse, error)
	CreateTransferRecipient(payload *types.CreateTransferRecipientRequest) (*types.CreateTransferRecipientResponse, error)
	Transfer(payload types.FianlTransferRequest) (*types.FinalTransferResponse, error)
}

type paymentServiceImpl struct {
	paystackClient paystack.Paystack
	binanceClient  binance.Binance
}

func NewPaymentService(b binance.Binance, p paystack.Paystack) PaymentService {
	return &paymentServiceImpl{p, b}
}

func (p *paymentServiceImpl) GetBankList() (*types.Banks, error) {

	res, err := p.paystackClient.GetBankList()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) ResolveAccountNumber(accountNumber string, bankCode string) (*types.AccountResolutionResponse, error) {
	res, err := p.paystackClient.ResolveAccountNumber(accountNumber, bankCode)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) CreateTransferRecipient(transferRecipient *types.CreateTransferRecipientRequest) (*types.CreateTransferRecipientResponse, error) {
	res, err := p.paystackClient.CreateTransferRecipient(transferRecipient)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) GetOTP(payload *types.TransferOtpRequest) (*types.TransferOtpResponse, error) {
	res, err := p.paystackClient.GetOTP(payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *paymentServiceImpl) Transfer(payload types.FianlTransferRequest) (*types.FinalTransferResponse, error) {
	res, err := p.paystackClient.Transfer(payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}
