package customerror

import "fmt"

type NoPermissionError struct {
	User string
	Area string
}

// Error реализует интерфейс error для NoPermissionError.
func (e *NoPermissionError) Error() string {
	return fmt.Sprintf("пользователь '%s' не имеет доступа к '%s'", e.User, e.Area)
}
