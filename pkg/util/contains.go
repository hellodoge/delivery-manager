package util

func ContainsInt(value int, arr []int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}