// responsible for creating and managing custom errors
package errors

// To be used when HTTP Status Unauthorized - 401 is to be sent back
type unauthorizedError struct {
	message string
}

func (e unauthorizedError) Error() string {
	return e.message
}

func GetUnauthorizedError(message string) unauthorizedError {
	return unauthorizedError{message: message}
}
