package conditions

type ConditionsBase struct {
	cacheKey *string
}

func (cb ConditionsBase) isArrayOrSlice(i interface{}) bool {
	switch i.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}
