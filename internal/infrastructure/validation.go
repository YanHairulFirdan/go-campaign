package infrastructure

import (
	"database/sql"

	validationRepo "go-campaign.com/internal/shared/repository"
	"go-campaign.com/pkg/validation"
)

func InitValidation(db *sql.DB) {
	validationRepository := validationRepo.NewDatabaseValidationRepository(db)

	validation.Init(validationRepository)
}
