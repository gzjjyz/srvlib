/**
 * @Author: zjj
 * @Date: 2026/2/2
 * @Desc:
**/

package utils

import "strings"

var moduleName []string

func GetModuleName() string {
	return strings.Join(moduleName, ",")
}

func ClearModuleName() {
	moduleName = nil
}
