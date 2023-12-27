package postgres

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func BuildFilter(query string, filters map[string][]string, spQuery ...SpecialQuery) (string, []interface{}, error) {
	var (
		listFilters []string
		inputArgs   []interface{}
	)

	fkeys := sortFilterKeys(filters)

	var i Iteration
	for _, fkey := range fkeys {
		fval := filters[fkey]

		if !validateValue(fval) {
			continue
		}

		if isNullFilter(fval) {
			listFilters = append(listFilters, fmt.Sprintf("%s %s", fkey, strings.Join(fval, "")))
			i = i.Increment()
			continue
		}

		listFilters = append(listFilters, fmt.Sprintf("%s IN (?)", fkey))
		i, inputArgs = i.Increment(), append(inputArgs, fval)
	}

	flt := strings.Join(listFilters, " AND ")

	for _, spq := range spQuery {
		if !validateValue(spq.Args) {
			continue
		}

		switch {
		case len(inputArgs) == 0:
			flt = spq.Query
		default:
			flt = strings.Join([]string{flt, spq.Type.String(), spq.Query}, " ")
		}

		for _, v := range spq.Args {
			inputArgs = append(inputArgs, v)
		}
	}

	query, args, err := sqlx.In(fmt.Sprintf(query, flt), inputArgs...)
	if err != nil {
		return "", nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	if flt == "" {
		query = strings.Replace(query, "WHERE", "", -1)
	}

	return query, args, nil
}

func sortFilterKeys(filters map[string][]string) []string {
	fkeys := make([]string, 0)
	for k := range filters {
		fkeys = append(fkeys, k)
	}
	sort.Strings(fkeys)

	return fkeys
}

func validateValue(fval []string) bool {
	switch {
	case len(fval) == 0:
		return false
	case isNullFilter(fval):
		return true
	default:
		return !cmp.Equal(fval, []string{""})
	}
}

type SpecialQuery struct {
	Query string
	Args  []string
	Type  QueryType
}

type QueryType int

const (
	And QueryType = iota
	Or
)

func (q QueryType) String() string {
	switch q {
	case And:
		return "AND"
	case Or:
		return "OR"
	default:
		return "AND"
	}
}

type Iteration int

func (i Iteration) IsFirstIteration() bool {
	return i == 0
}

func (i Iteration) Increment() Iteration {
	return i + 1
}

func isNullFilter(req []string) bool {
	for i := range req {
		if req[i] == "IS NULL" {
			return true
		}
		if req[i] == "IS NOT NULL" {
			return true
		}
	}
	return false
}
