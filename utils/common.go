package utils

func Contains(slice interface{}, value interface{}) bool {
	switch v := slice.(type) {
	case []int:
		for _, item := range v {
			if item == value {
				return true
			}
		}
	case []string:
		for _, item := range v {
			if item == value {
				return true
			}
		}
	default:
		return false
	}
	return false
}
