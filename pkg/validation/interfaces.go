package validation

type DatabaseValidationRepository interface {
	// IsUnique
	IsUnique(table, column string, value any) (bool, error)
	IsUniqueWithCondition(table, column, conditionColumn string, value, conditionValue any) (bool, error)
}

type IsUniqueFunc func(table, column string, value any) (bool, error)
type IsUniqueWithConditionFunc func(table, column, conditionColumn string, value, conditionValue any) (bool, error)
