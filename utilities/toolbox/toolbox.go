package toolbox

func SliceContains[T comparable](arr []T, elem T) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ArePointersEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a != nil && b != nil {
		return a == b
	}

	return false
}
