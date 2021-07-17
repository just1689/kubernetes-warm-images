package util

import (
	"fmt"
	"os"
)

var HostIP = os.Getenv("HOST_IP")

func LogPrepend(level int, text string) (result string) {
	result = ""
	if level == 1 {
		result += ">   "
	} else if level == 2 {
		result += ">>  "
	} else {
		result += "    "
	}
	result = fmt.Sprintf("%s %s: %s", HostIP, result, text)
	return
}
