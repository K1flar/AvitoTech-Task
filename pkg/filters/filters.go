package filters

import (
	"strconv"
)

func parseIntWithDefaultValue(str string, defaultValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = defaultValue
	}
	return val
}
