package handlers

import (
	er "errors"
	"net/http"

	"github.com/Creative-genius001/Stacklo/pkg/paystack"
	"github.com/Creative-genius001/Stacklo/services/payment/api/services"
	"github.com/Creative-genius001/Stacklo/services/payment/pkg/binance"
	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PaymentService struct {
	payment services.PaymentService
}

func NewpaymentService(ps services.PaymentService) *PaymentService {
	return &PaymentService{payment: ps}
}

func (s *PaymentService) GetBankList(c *gin.Context) {
	banks, err := s.payment.GetBankList()
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from PaystackAPIWrapper", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": banks})
}

func (s *PaymentService) ResolveAccountNumber(c *gin.Context) {
	accountNumber := c.Query("account_number")
	bankCode := c.Query("bank_code")
	if accountNumber == "" && bankCode == "" {
		logger.Logger.Error("Account number and Bank code are required")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	resp, err := s.payment.ResolveAccountNumber(accountNumber, bankCode)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from PaystackAPIWrapper")
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (s *PaymentService) GetOTP(c *gin.Context) {
	var inputData paystack.StartTransferData
	if err := c.ShouldBindJSON(&inputData); err != nil {
		logger.Logger.Error("Invalid input data", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	trfRecipientData := paystack.CreateTransferRecipientRequest{
		Type:          "nuban",
		Name:          inputData.Name,
		AccountNumber: inputData.AccountNumber,
		BankCode:      inputData.BankCode,
		Currency:      "NGN",
	}

	recipient, err := s.payment.CreateTransferRecipient(&trfRecipientData)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from PaystackAPIWrapper")
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	referenceNumber := uuid.New().String()

	otpRequestData := paystack.TransferOtpRequest{
		Source:    "balance",
		Reason:    inputData.Reason,
		Amount:    inputData.Amount,
		Recipeint: recipient.Data.RecipientCode,
		Reference: referenceNumber,
	}

	otp, err := s.payment.GetOTP(&otpRequestData)
	// if err != nil {
	// 	signedStr, tokenErr := utils.CreateRetryTokenOtpRequest(referenceNumber, recipient.Data.RecipientCode)
	// 	if tokenErr != nil {
	// 		res := utils.NewError(http.StatusInternalServerError, tokenErr.Error())
	// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{
	// 			"error": "Generate Otp failed and retry token creation failed",
	// 		})
	// 		return
	// 	}

	// 	//hmacSignature := utils.GenerateHMAC(signedStr)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":       err.Error(),
	// 		"retry_token": signedStr,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"transfer_code": otp.Data.TransferCode})
}

// func RetryOtp(c *gin.Context) {

// 	type RetryOtp struct {
// 		Recipient string `json:"recipient"`
// 		Reference string `json:"reference"`
// 		Reason    string `json:"reason"`
// 		Amount    string `json:"amount"`
// 	}

// 	type RetryToken struct {
// 		RetryToken string `json:"retry_token"`
// 		Reason     string `json:"reason"`
// 		Amount     string `json:"amount"`
// 	}

// 	var inputData RetryToken
// 	if err := c.ShouldBindJSON(&inputData); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid retry payload"})
// 		return
// 	}

// 	parsed, err := utils.VerifyRetryToken(inputData.RetryToken)
// 	if err != nil {
// 		res := utils.NewError(http.StatusBadRequest, err.Error())
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{
// 			"error": res.Error,
// 		})
// 		return
// 	}

// 	otpRequestData := paystack.TransferOtpRequest{
// 		Source:    "balance",
// 		Reason:    inputData.Reason,
// 		Amount:    inputData.Amount,
// 		Recipeint: parsed.Recipient,
// 		Reference: parsed.Reference,
// 	}

// 	otp, err := services.RequestOTP(&otpRequestData)
// 	if err != nil {
// 		res := utils.NewError(http.StatusInternalServerError, err.Error())
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{
// 			"error": "Generate Otp failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"transfer_code": otp.Data.TransferCode})
// }

func (s *PaymentService) Transfer(c *gin.Context) {
	var data paystack.FianlTransferRequest

	if err := c.BindJSON(&data); err != nil {
		logger.Logger.Error("Invalid input data", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	resp, err := s.payment.Transfer(data)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from PaystackAPIWrapper")
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	json := paystack.FinalTransferJson{
		Amount:    resp.Data.Amount,
		Recipient: resp.Data.Recipient,
		Reference: resp.Data.Reference,
		Reason:    resp.Data.Reason,
		Status:    resp.Data.Status,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": json})
}

func (s *PaymentService) Ping(c *gin.Context) {

	err := s.payment.Ping()
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from BINANCE WRAPPER")
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "type": "ping"})
}

func (s *PaymentService) Convert(c *gin.Context) {
	var data binance.ConvertAssetRequest
	if err := c.BindJSON(&data); err != nil {
		logger.Logger.Error("Invalid input data", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	resp, err := s.payment.Convert(data)

	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from BINANCE WRAPPER")
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			logger.Logger.Error("Unexpected error fuck", zap.Error(appErr))
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeExternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeExternal})
			return
		default:
			logger.Logger.Error("Unexpected error fuck")
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
