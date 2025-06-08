package infrastucture

import (
	"database/sql"
	"fmt"

	validationRepo "go-campaign.com/internal/shared/repository"
	"go-campaign.com/pkg/validation"
)

func InitValidation(db *sql.DB) {
	validationRepository := validationRepo.NewDatabaseValidationRepository(db)

	err := validation.Init(validationRepository)

	if err != nil {
		panic(fmt.Sprintf("Error initializing validation: %v", err))
	}
}
