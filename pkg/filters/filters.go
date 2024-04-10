package filters

import (
	"net/url"
	"strconv"
)

func getFromQuery(q url.Values, param string) (string, bool) {
	if val := q.Get(param); val != "" {
		return val, true
	}
	return "", false
}

func parseIntWithDefaultValue(str string, defaultValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = defaultValue
	}
	return val
}
