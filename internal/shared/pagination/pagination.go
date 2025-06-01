package pagination

type Builder struct {
	DataRetriever  func(page, perPage int) ([]any, error)
	CountRetriever func() (int, error)
}
