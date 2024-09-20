package connor

import (
	"github.com/sourcenetwork/defradb/client"

	"github.com/sourcenetwork/immutable"
)

// anyOp is an operator which allows the evaluation of
// a number of conditions over a list of values
// matching if any of them match.
func anyOp(condition, data any) (bool, error) {
	switch t := data.(type) {
	case []string:
		return anySlice(condition, t)

	case []immutable.Option[string]:
		return anySlice(condition, t)

	case []int64:
		return anySlice(condition, t)

	case []immutable.Option[int64]:
		return anySlice(condition, t)

	case []bool:
		return anySlice(condition, t)

	case []immutable.Option[bool]:
		return anySlice(condition, t)

	case []float64:
		return anySlice(condition, t)

	case []immutable.Option[float64]:
		return anySlice(condition, t)

	default:
		return false, client.NewErrUnhandledType("data", data)
	}
}

func anySlice[T any](condition any, data []T) (bool, error) {
	for _, c := range data {
		m, err := eq(condition, c)
		if err != nil {
			return false, err
		} else if m {
			return true, nil
		}
	}
	return false, nil
}
