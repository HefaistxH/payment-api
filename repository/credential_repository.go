package repository

import (
	"database/sql"
	"mnc-techtest/config"
	"mnc-techtest/entity"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type CredentialRepository interface {
	GetCredByEmail(email string) (entity.Credential, error)
}

type credentialRepository struct {
	db *sql.DB
}

func (r *credentialRepository) GetCredByEmail(email string) (entity.Credential, error) {
	var credential entity.Credential
	err := r.db.QueryRow(config.GetCredByEmailQuery, email).Scan(&credential.Id, &credential.UserId, &credential.Email, &credential.Password, &credential.Role)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("Error getting credential by email: %v", err)
		return entity.Credential{}, err
	}
	return credential, nil
}

func NewCredentialRepository(db *sql.DB) CredentialRepository {
	return &credentialRepository{db: db}
}
