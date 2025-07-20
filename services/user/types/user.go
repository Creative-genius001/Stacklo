package types

type UpdateUser struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Country   *string `json:"country"`
	KycStatus *string `json:"kyc_status"`
}
