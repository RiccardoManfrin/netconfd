package comm

// ErrorString is an error string
type ErrorString struct {
	S string
}

// Error returns the string of the error (wow!)
func (e *ErrorString) Error() string {
	return e.S
}
