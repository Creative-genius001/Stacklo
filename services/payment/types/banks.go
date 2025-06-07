package types

type Banks struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		Name        string      `json:"name"`
		Slug        string      `json:"slug"`
		Code        string      `json:"code"`
		Longcode    string      `json:"longcode"`
		Gateway     interface{} `json:"gateway"`
		PayWithBank bool        `json:"pay_with_bank"`
		Active      bool        `json:"active"`
		IsDeleted   bool        `json:"is_deleted,omitempty"`
		Country     string      `json:"country,omitempty"`
		Currency    string      `json:"currency,omitempty"`
		Type        string      `json:"type,omitempty"`
		CreatedAt   string      `json:"createdAt,omitempty"`
		UpdatedAt   string      `json:"updatedAt,omitempty"`
	} `json:"data"`
}
