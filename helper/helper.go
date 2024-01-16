package helper

type Variable struct {
	Key, Value, Description, Category string
	HCL, Sensitive                    bool
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
