package paystack

import (
	"encoding/json"
	"net/http"

	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	"go.uber.org/zap"
)

func (p *paystackClient) GetBankList() (*Banks, error) {

	resp, err := p.PaystackAPIWrapper(http.MethodGet, string(UBankList), nil, nil)
	if err != nil {
		return nil, err
	}

	var result Banks
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) ResolveAccountNumber(accountNumber string, bankCode string) (*AccountResolutionResponse, error) {
	url := string(UResolveAccNum) + accountNumber + "&bank_code=" + bankCode

	resp, err := p.PaystackAPIWrapper(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	var result AccountResolutionResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) CreateTransferRecipient(payload *CreateTransferRecipientRequest) (*CreateTransferRecipientResponse, error) {

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UCreateTrfRecpt), nil, bytes)
	if err != nil {
		return nil, err
	}

	var result CreateTransferRecipientResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) GetOTP(payload *TransferOtpRequest) (*TransferOtpResponse, error) {

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UTransfer), nil, bytes)
	if err != nil {
		return nil, err
	}

	var result TransferOtpResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) Transfer(payload FianlTransferRequest) (*FinalTransferResponse, error) {

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UFTransfer), nil, bytes)
	if err != nil {
		return nil, err
	}

	var result FinalTransferResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) CreateCustomer(payload CreateCustomerRequest) (*CreateCustomerResponse, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UTransfer), nil, bytes)
	if err != nil {
		return nil, err
	}

	var result CreateCustomerResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) CreateDVAWallet(payload *CreateDVAWalletRequest) (*CreateDVAWalletResponse, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.PaystackAPIWrapper(http.MethodPost, string(UTransfer), nil, bytes)
	if err != nil {
		return nil, err
	}

	var result CreateDVAWalletResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		p.logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}
