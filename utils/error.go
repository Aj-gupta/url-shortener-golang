package utils

type Error interface {
	error
	Status() int
	ErrorArr() []string
	ErrorLength() int
}

// CustomServiceErr represents an error with an associated HTTP status code.
type CustomAPIErr struct {
	Code      int
	Err       error
	Errors    map[string]interface{}
	ErrLength int
}

// Allows CustomServiceErr to satisfy the error interface.
func (e CustomAPIErr) Error() string {
	return e.Err.Error()
}

func (e CustomAPIErr) ErrorArr() map[string]interface{} {
	return e.Errors
}

// Returns the HTTP status code.
func (e CustomAPIErr) Status() int {
	return e.Code
}

// Returns the HTTP error length.
func (e CustomAPIErr) ErrorLength() int {
	return e.ErrLength
}
