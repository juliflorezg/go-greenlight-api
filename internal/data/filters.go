package data

import (
	"strings"

	"github.com/juliflorezg/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page <= 10_000_000, "page", "must be less than or equal to 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be less than or equal to 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "must be a valid sort value")
}

// Check that the client-provided Sort field matches one of the entries in our safe-list
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
