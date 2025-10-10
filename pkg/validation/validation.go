package validation

var (
	databaseRepository DatabaseValidationRepository
)

func Init(db DatabaseValidationRepository) {
	databaseRepository = db
}
