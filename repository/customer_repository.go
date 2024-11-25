package repository

import (
	"database/sql"
	"mnc-techtest/config"
	"mnc-techtest/entity/dto"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type CustomerRepository interface {
	CheckCustomerByEmail(email string) (bool, error)
	CheckCustomerBalance(userId string) (float64, error)
	CheckMerchantBalance(userId string) (float64, error)
	Payment(payment dto.Payment) (dto.Payment, error)
}

type customerRepository struct {
	db *sql.DB
}

func (r *customerRepository) CheckCustomerByEmail(email string) (bool, error) {
	isExist := false
	err := r.db.QueryRow(config.CheckUserByEmailQuery, email).Scan(&isExist)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("Error checking user by email: %v", err)
		return false, err
	}
	return isExist, nil
}

func (r *customerRepository) CheckCustomerBalance(userId string) (float64, error) {
	balance := 0.00
	err := r.db.QueryRow(config.CheckCustomerBalanceQuery, userId).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("Error checking customer balance: %v", err)
		return balance, err
	}
	return balance, nil
}

func (r *customerRepository) CheckMerchantBalance(userId string) (float64, error) {
	balance := 0.00
	err := r.db.QueryRow(config.CheckMerchantBalanceQuery, userId).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("Error checking merchant balance: %v", err)
		return balance, err
	}
	return balance, nil
}

// func (r *userRepository) UpdateCustomerBalanceQuery(userId string, balance float64) error {
// 	_, err := r.db.Exec(config.UpdateCustomerBalanceQuery, balance, userId)
// 	if err != nil {
// 		logrus.Errorf("Error updating customer balance: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func (r *userRepository) UpdateMerchantBalanceQuery(userId string, balance float64) error {
// 	_, err := r.db.Exec(config.UpdateMerchantBalanceQuery, balance, userId)
// 	if err != nil {
// 		logrus.Errorf("Error updating merchant balance: %v", err)
// 		return err
// 	}
// 	return nil
// }

func (r *customerRepository) Payment(payment dto.Payment) (dto.Payment, error) {
	tx, err := r.db.Begin()
	if err != nil {
		logrus.Errorf("Error starting transaction: %v", err)
		return dto.Payment{}, err
	}

	_, err = tx.Exec(config.UpdateCustomerBalanceQuery, payment.Amount, payment.CustomerId)
	if err != nil {
		logrus.Errorf("Error updating customer balance: %v", err)
		tx.Rollback()
		return dto.Payment{}, err
	}

	_, err = tx.Exec(config.UpdateMerchantBalanceQuery, payment.Amount, payment.MerchantId)
	if err != nil {
		logrus.Errorf("Error updating merchant balance: %v", err)
		tx.Rollback()
		return dto.Payment{}, err
	}

	_, err = tx.Exec(config.AddHistoryQuery, payment.CustomerId, payment.MerchantId, "Payment", payment.Amount, payment.Message, time.Now())
	if err != nil {
		logrus.Errorf("Error adding history: %v", err)
		tx.Rollback()
		return dto.Payment{}, err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("Error committing transaction: %v", err)
		tx.Rollback()
		return dto.Payment{}, err
	}

	return payment, nil

}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}
