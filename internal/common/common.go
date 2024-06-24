// пакет общих ф-ий проекта
package common

import (
	"strings"
)

// ф-я проверки что строка не является пробелом / не состоит только из пробелов
func ValidString(str string) bool {
	return strings.TrimSpace(str) != ""
}
