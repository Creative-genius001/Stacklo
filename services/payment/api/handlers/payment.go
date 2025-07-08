package handlers

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/payment/api/services"
	errors "github.com/Creative-genius001/Stacklo/services/payment/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/gin-gonic/gin"
)

func GetBankList(c *gin.Context) {

	banks, err := services.GetBankList(c.Request.Context())

	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			logger.Logger.Error("Unable to assert err with CustomError type")
			c.JSON(errors.GetHTTPStatus(errors.TypeConflict), gin.H{"status": "error", "message": errors.TypeConflict})
			return
		} else {
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
	}

	c.JSON(http.StatusOK, gin.H{"data": banks})
}

// func ResolveAccountNumber(c *gin.Context) {

// 	accountNumber := c.Query("account_number")
// 	bankCode := c.Query("bank_code")
// 	if accountNumber == "" && bankCode == "" {
// 		res := utils.NewError(http.StatusInternalServerError, "account_number and bank_code parameters are required")
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
// 		return
// 	}

// 	resp, err := services.ResolveAccountNumber(accountNumber, bankCode)
// 	if err != nil {
// 		res := utils.NewError(http.StatusInternalServerError, err.Error())
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": resp})
// }

// func RequestOTP(c *gin.Context) {
// 	var inputData types.StartTransferData
// 	if err := c.ShouldBindJSON(&inputData); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
// 		return
// 	}

// 	tranferRecipientInputData := types.CreateTransferRecipientRequest{
// 		Type:          "nuban",
// 		Name:          inputData.Name,
// 		AccountNumber: inputData.AccountNumber,
// 		BankCode:      inputData.BankCode,
// 		Currency:      "NGN",
// 	}

// 	recipient, err := services.CreateTransferRecipient(&tranferRecipientInputData)
// 	if err != nil {
// 		res := utils.NewError(http.StatusInternalServerError, err.Error())
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
// 		return
// 	}

// 	referenceNumber := uuid.New().String()

// 	otpRequestData := types.TransferOtpRequest{
// 		Source:    "balance",
// 		Reason:    inputData.Reason,
// 		Amount:    inputData.Amount,
// 		Recipeint: recipient.Data.RecipientCode,
// 		Reference: referenceNumber,
// 	}

// 	otp, err := services.RequestOTP(&otpRequestData)
// 	if err != nil {
// 		signedStr, tokenErr := utils.CreateRetryTokenOtpRequest(referenceNumber, recipient.Data.RecipientCode)
// 		if tokenErr != nil {
// 			res := utils.NewError(http.StatusInternalServerError, tokenErr.Error())
// 			c.AbortWithStatusJSON(res.StatusCode, gin.H{
// 				"error": "Generate Otp failed and retry token creation failed",
// 			})
// 			return
// 		}

// 		//hmacSignature := utils.GenerateHMAC(signedStr)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error":       err.Error(),
// 			"retry_token": signedStr,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"transfer_code": otp.Data.TransferCode})
// }

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

// 	otpRequestData := types.TransferOtpRequest{
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

// func Transfer(c *gin.Context) {
// 	var tranferData types.QueuedTransferRequest

// 	if err := c.BindJSON(&tranferData); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid data payload"})
// 		return
// 	}

// 	resp, err := services.FinalTransfer(&tranferData)
// 	if err != nil {
// 		res := utils.NewError(http.StatusInternalServerError, err.Error())
// 		c.AbortWithStatusJSON(res.StatusCode, gin.H{
// 			"error": "Transfer failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": resp})
// }
