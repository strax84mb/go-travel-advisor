package utils

func IsPresent[V int](list []V, check func(V) bool) bool {
	for _, v := range list {
		if check(v) {
			return true
		}
	}
	return false
}
