package service

import (
	"context"

	"github.com/Creative-genius001/Stacklo/pkg/paystack"
	"github.com/Creative-genius001/Stacklo/services/wallet/model"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	GetFiatWallet(ctx context.Context, id string) (*model.Wallet, error)
	CreateFiatWallet(ctx context.Context, payload paystack.CreateCustomerRequest) (*model.Wallet, error)
	CreateCryptoWallet(ctx context.Context, wt model.Wallet) error
	GetAllWallets(ctx context.Context, id string) ([]*model.Wallet, error)
	CreateWalletPaystack(c paystack.CreateCustomerRequest) (*model.Wallet, error)
}

type walletService struct {
	repository Repository
	paystack   paystack.Paystack
}

func NewService(r Repository, p paystack.Paystack) Service {
	return &walletService{r, p}
}

func (s *walletService) CreateWalletPaystack(c paystack.CreateCustomerRequest) (*model.Wallet, error) {

	// res, err := s.paystack.CreateCustomer(c)
	// if err != nil {
	// 	return nil, err
	// }

	// crtWalletRq := paystack.CreateDVAWalletRequest{
	// 	FirstName:     c.FirstName,
	// 	LastName:      c.LastName,
	// 	Phone:         c.Phone,
	// 	CustomerCode:   res.Data.ID,
	// 	PreferredBank: "wema-bank",
	// }

	// wallet, err := s.paystack.CreateDVAWallet(&crtWalletRq)
	// if err != nil {
	// 	return nil, err
	// }

	accountName := c.FirstName + c.LastName + "/ PAYSTACK"

	wPayStack := model.Wallet{
		Currency:             "NGN",
		Active:               true,
		UserId:               c.ID,
		VirtualAccountName:   utils.StringPtr(accountName),
		VirtualAccountNumber: utils.StringPtr("0091728654"),
		VirtualBankName:      utils.StringPtr("Providus Bank"),
		WalletType:           "FIAT",
	}

	// wPayStack := model.Wallet{
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

func (w walletService) GetFiatWallet(ctx context.Context, id string) (*model.Wallet, error) {
	wallet, err := w.repository.GetFiatWallet(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(errors.TypeNotFound, "wallet not found", err)
		}
		return nil, err
	}

	return wallet, nil
}

func (w walletService) GetAllWallets(ctx context.Context, id string) ([]*model.Wallet, error) {
	wallets, err := w.repository.GetAllWallets(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(errors.TypeNotFound, "wallet not found", err)
		}
		return nil, err
	}

	return wallets, nil
}

func (w walletService) CreateFiatWallet(ctx context.Context, payload paystack.CreateCustomerRequest) (*model.Wallet, error) {

	cs, err := w.CreateWalletPaystack(payload)
	if err != nil {
		return nil, err
	}
	walletPayload := *cs
	wallet, err := w.repository.CreateFiatWallet(ctx, walletPayload)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w walletService) CreateCryptoWallet(ctx context.Context, wt model.Wallet) error {
	wt.WalletType = "CRYPTO"
	err := w.repository.CreateCryptoWallet(ctx, wt)
	if err != nil {
		return err
	}
	return nil
}
