package bankmuscatpg

import (
	"fmt"
	"net/url"
	"strings"
)

func stringToMap(requestData string) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{})

	keyValuePairs := strings.Split(requestData, "&")

	for _, keyValue := range keyValuePairs {
		parts := strings.SplitN(keyValue, "=", 2)
		key, err := url.QueryUnescape(parts[0])
		if err != nil {
			return nil, err
		}
		value, err := url.QueryUnescape(parts[1])
		if err != nil {
			return nil, err
		}
		requestMap[key] = value
	}

	return requestMap, nil
}

func mapToString(requestMap map[string]interface{}) string {
	var requestData string
	for key, value := range requestMap {
		requestData += key + "=" + url.QueryEscape(fmt.Sprint(value)) + "&"
	}

	if len(requestData) > 0 {
		requestData = requestData[:len(requestData)-1]
	}
	return requestData
}
