package usecase

import (
	"fmt"

	"mnc-techtest/entity/dto"
	"mnc-techtest/repository"
	"mnc-techtest/shared/service"

	"github.com/sirupsen/logrus"
)

type CustomerUsecase interface {
	CheckUserByEmail(email string) (bool, error)
	CheckCustomerBalance(userId string) (float64, error)
	CheckMerchantBalance(userId string) (float64, error)
	Payment(token string, payment dto.Payment) (dto.Payment, error)
}

type customerUsecase struct {
	repo       repository.CustomerRepository
	jwtService service.JwtService
}

func (c *customerUsecase) CheckUserByEmail(email string) (bool, error) {
	return c.repo.CheckCustomerByEmail(email)
}

func (c *customerUsecase) CheckCustomerBalance(userId string) (float64, error) {
	return c.repo.CheckCustomerBalance(userId)
}

func (c *customerUsecase) CheckMerchantBalance(userId string) (float64, error) {
	return c.repo.CheckMerchantBalance(userId)
}

func (c *customerUsecase) Payment(token string, payment dto.Payment) (dto.Payment, error) {
	if c.jwtService == nil {
		logrus.Error("jwtService is nil")
		return dto.Payment{}, fmt.Errorf("internal server error")
	}

	if token == "" {
		logrus.Error("Token is empty")
		return dto.Payment{}, fmt.Errorf("invalid token")
	}

	logrus.Infof("Parsing token: %s", token)
	claims, err := c.jwtService.ParseToken(token)
	if err != nil {
		logrus.Errorf("Invalid token: %v", err)
		return dto.Payment{}, fmt.Errorf("invalid token: %v", err)
	}

	claimUserID, ok := claims["user_id"].(string)
	if !ok {
		logrus.Error("User ID claim missing or invalid")
		return dto.Payment{}, fmt.Errorf("user_id from token does not match customer_id in payment")
	}
	logrus.Infof("Claimed user_id: %s", claimUserID)
	if claimUserID != payment.CustomerId {
		logrus.Errorf("Token user_id %s does not match payment customer_id %s", claimUserID, payment.CustomerId)
		return dto.Payment{}, fmt.Errorf("user_id from token does not match customer_id in payment")
	}

	customerBalance, _ := c.repo.CheckCustomerBalance(payment.CustomerId)
	merchantBalance, _ := c.repo.CheckMerchantBalance(payment.MerchantId)
	if customerBalance < payment.Amount {
		return dto.Payment{}, fmt.Errorf("customer balance is not enough")
	}
	customerBalance = customerBalance - payment.Amount
	merchantBalance = merchantBalance + payment.Amount

	payment, err = c.repo.Payment(payment, customerBalance, merchantBalance)
	if err != nil {
		return dto.Payment{}, err
	}

	return payment, nil
}

func NewCustomerUsecase(repo repository.CustomerRepository, jwtService service.JwtService) CustomerUsecase {
	return &customerUsecase{repo: repo, jwtService: jwtService}
}
