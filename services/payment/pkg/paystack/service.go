package paystack

import (
	"encoding/json"
	er "errors"
	"net/http"

	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"go.uber.org/zap"
)

func (p *paystackClient) GetBankList() (*Banks, error) {

	resp, err := p.PaystackAPIWrapper(http.MethodGet, string(UBankList), nil, nil)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	var result Banks
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}

func (p *paystackClient) ResolveAccountNumber(accountNumber string, bankCode string) (*AccountResolutionResponse, error) {
	url := string(UResolveAccNum) + accountNumber + "&bank_code=" + bankCode

	resp, err := p.PaystackAPIWrapper(http.MethodGet, url, nil, nil)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	var result AccountResolutionResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
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
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	var result CreateTransferRecipientResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
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
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	var result TransferOtpResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
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
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			return nil, errors.Wrap(errors.TypeExternal, "External API error", err)
		}

		logger.Logger.Error("It aserted at the service level", zap.String("messg", appErr.Message), zap.String("type", string(appErr.Type)), zap.Error(appErr.Err))
		return nil, appErr
	}

	var result FinalTransferResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logger.Logger.Warn("Failed to unmarshal request body", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Failed to unmarshal request body", err)
	}

	return &result, nil
}
