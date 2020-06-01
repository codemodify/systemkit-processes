package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
)

func AsJSONString(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("ERROR: AsJSONString(), details [%s]", err.Error())
	}
	return string(bytes)
}

func IsNullOrEmpty(value string) bool {
	return (len(strings.TrimSpace(value)) <= 0)
}

