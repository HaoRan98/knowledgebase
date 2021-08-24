package util

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJson(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func FormJson(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	var str string
	for _, k := range strs {
		if len(k) == 0 {
			return ""
		}

		str = strs[0]
		for !strings.HasPrefix(k, str) {
			str = str[:len(str)-1]
		}
		if len(str) == 0 {
			return ""
		}
	}
	return str
}
