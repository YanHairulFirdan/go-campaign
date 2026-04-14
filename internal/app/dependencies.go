package app

import (
	"database/sql"

	"go-campaign.com/internal/config"
	"go-campaign.com/internal/shared/services/payment"
	"go-campaign.com/pkg/filesystem"
)

type Dependencies struct {
	Config         *config.Config
	DB             *sql.DB
	FileSystem     filesystem.Filesystem
	PaymentGateway payment.PaymentGateway
}

func NewDependencies(
	config *config.Config,
	db *sql.DB,
	fileSystem filesystem.Filesystem,
	paymentGateway payment.PaymentGateway,
) *Dependencies {
	return &Dependencies{
		Config:         config,
		DB:             db,
		FileSystem:     fileSystem,
		PaymentGateway: paymentGateway,
	}
}

func (d *Dependencies) CloseDatabaseConnection() error {
	if d.DB != nil {
		return d.DB.Close()
	}

	return nil
}
