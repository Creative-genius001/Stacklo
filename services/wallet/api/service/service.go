package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/wallet/config"
	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Service interface {
	GetFiatWallet(ctx context.Context, id string) (*types.Wallet, error)
	CreateFiatWallet(ctx context.Context, wt types.Wallet) (*types.Wallet, error)
	CreateCryptoWallet(ctx context.Context, wt types.Wallet) error
	GetAllWallets(ctx context.Context, id string) ([]*types.Wallet, error)
}

type walletService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &walletService{r}
}

func CreateCustomer(customerReq types.CreateCustomerRequest) (*types.CreateCustomerResponse, error) {
	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	customerReqJSON, err := json.Marshal(customerReq)
	if err != nil {
		logger.Logger.Warn("Failed to marshal customer request")
		return nil, errors.New(errors.TypeInternal, "Failed to marshal customer request")
	}

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, _ := http.NewRequest("POST", PAYSTACK_BASE_URL+"/customer", bytes.NewBuffer(customerReqJSON))

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Logger.Error("Paystack API error: Failed to establish a connection", zap.String("API URL", PAYSTACK_BASE_URL+"/dedicated_account"))
		return nil, errors.Wrap(errors.TypeExternal, "Paystack API error: Failed to establish a connection", err)
	}

	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		logger.Logger.Warn("Failed to read response data", zap.String("errorMsg", string(errorBody)))
		return nil, errors.New(errors.TypeInternal, "Failed to read response data")
	}

	var customer types.CreateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		logger.Logger.Warn("failed to decode response")
		return nil, errors.New(errors.TypeInternal, "Failed to decode response")
	}

	return &customer, nil
}

func CreateDVAWallet(createWalletReq *types.CreateDVAWalletRequest) (*types.CreateDVAWalletResponse, error) {
	var wallet types.CreateDVAWalletResponse

	PAYSTACK_BASE_URL := config.Cfg.PaystackBaseUrl
	PAYSTACK_API_KEY := config.Cfg.PaystackTestKey

	createWalletReqJSON, err := json.Marshal(createWalletReq)
	if err != nil {
		logger.Logger.Warn("Failed to marshal customer request")
		return nil, errors.New(errors.TypeInternal, "Failed to marshal customer request")
	}

	// Create client with timeout and retry
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	//create request
	req, _ := http.NewRequest("POST", PAYSTACK_BASE_URL+"/dedicated_account", bytes.NewBuffer(createWalletReqJSON))

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+PAYSTACK_API_KEY)

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Logger.Error("Paystack API error: Failed to establish a connection", zap.String("API URL", PAYSTACK_BASE_URL+"/dedicated_account"))
		return nil, errors.Wrap(errors.TypeExternal, "Paystack API error: Failed to establish a connection", err)
	}

	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		logger.Logger.Warn("Failed to read response data", zap.String("errorMsg", string(errorBody)))
		return nil, errors.New(errors.TypeInternal, "Failed to read response data")
	}

	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		logger.Logger.Warn("failed to decode response")
		return nil, errors.New(errors.TypeInternal, "Failed to decode response")
	}

	return &wallet, nil
}

func CreateWalletPaystack(c types.CreateCustomerRequest) (*types.Wallet, error) {

	// customer, err := services.CreateCustomer(c)
	// if err != nil {
	// 	return nil, err
	// }

	// createWalletReq := types.CreateDVAWalletRequest{
	// 	FirstName:     c.FirstName,
	// 	LastName:      c.LastName,
	// 	Phone:         c.Phone,
	// 	CustomerCode:      3356, //customer.Data.Code
	// 	PreferredBank: "wema-bank",
	// }

	// wallet, err := services.CreateDVAWallet(&createWalletReq)
	// if err != nil {
	// 	return nil, err
	// }

	accountName := c.FirstName + c.LastName + "/ PAYSTACK"

	wPayStack := types.Wallet{
		Currency:             "NGN",
		Active:               true,
		VirtualAccountName:   utils.StringPtr(accountName),
		VirtualAccountNumber: utils.StringPtr("0091728654"),
		VirtualBankName:      utils.StringPtr("Providus Bank"),
		WalletType:           "FIAT",
	}

	// wPayStack := types.Wallet{
	// 	Currency:             wallet.Data.Currency,
	// 	Active:               wallet.Data.Active,
	// 	VirtualAccountName:   wallet.Data.AccountName,
	// 	VirtualAccountNumber: wallet.Data.AccountNumber,
	// 	VirtualBankName:      wallet.Data.Bank.Name,
	// 	CreatedAt:            wallet.Data.CreatedAt,
	// 	UpdatedAt:            wallet.Data.UpdatedAt,
	// }

	return &wPayStack, nil
}

func (w walletService) GetFiatWallet(ctx context.Context, id string) (*types.Wallet, error) {
	wallet, err := w.repository.GetFiatWallet(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(errors.TypeNotFound, "wallet not found", err)
		}
		return nil, err
	}

	return wallet, nil
}

func (w walletService) GetAllWallets(ctx context.Context, id string) ([]*types.Wallet, error) {
	wallets, err := w.repository.GetAllWallets(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(errors.TypeNotFound, "wallet not found", err)
		}
		return nil, err
	}

	return wallets, nil
}

func (w walletService) CreateFiatWallet(ctx context.Context, wt types.Wallet) (*types.Wallet, error) {
	wt.UserId = "a1b2c3d4-e5f6-4789-90ab-cdef01234567"
	wallet, err := w.repository.CreateFiatWallet(ctx, wt)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w walletService) CreateCryptoWallet(ctx context.Context, wt types.Wallet) error {
	wt.ID = uuid.New().String()
	wt.WalletType = "CRYPTO"
	err := w.repository.CreateCryptoWallet(ctx, wt)
	if err != nil {
		return err
	}
	return nil
}
