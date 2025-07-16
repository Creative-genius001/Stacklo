package types

type LoginType struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterType struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Phone     string `form:"phone" json:"phone" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Country   string `form:"country" json:"country" binding:"required"`
}
