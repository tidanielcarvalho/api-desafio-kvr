package repositories

// Params to sort query
type SortParams struct {
	Field string
	Asc   bool
}

func SortDefault() SortParams {
	return SortParams{"name", true}
}
