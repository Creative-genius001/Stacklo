package services

import (
	"errors"

	"github.com/Creative-genius001/user/db"
	"github.com/Creative-genius001/user/models"
	"github.com/Creative-genius001/user/types"
	"github.com/google/uuid"
)

func LoginService() {

}

func RegisterService(RegForm types.RegisterType) error {
	user := models.User{
		Id:           uuid.New().String(),
		Phone:        RegForm.Phone,
		PasswordHash: RegForm.Password,
		Email:        RegForm.Email,
		FirstName:    RegForm.FirstName,
		LastName:     RegForm.LastName,
	}

	res := db.DB.Create(&user)
	if res.Error != nil {
		return errors.New("signup failed! server error")
	}
	return nil
}
