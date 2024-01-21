package helper

import "log"

type Variable struct {
	Key, Value, Description, Category string
	HCL, Sensitive                    bool
}
type TfConfig struct {
	Address      string `yaml:"address"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
