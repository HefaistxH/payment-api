package usecase

import (
	"mnc-techtest/entity/dto"
	"mnc-techtest/repository"
)

type CustomerUsecase interface {
	CheckUserByEmail(email string) (bool, error)
	CheckCustomerBalance(userId string) (float64, error)
	CheckMerchantBalance(userId string) (float64, error)
	Payment(payment dto.Payment) (dto.Payment, error)
}

type customerUsecase struct {
	repo repository.CustomerRepository
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

func (c *customerUsecase) Payment(payment dto.Payment) (dto.Payment, error) {
	return c.repo.Payment(payment)
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{repo: repo}
}
