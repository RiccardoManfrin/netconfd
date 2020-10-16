package comm

// ErrorString is an error string
type ErrorString struct {
	S string
}

// Error returns the string of the error (wow!)
func (e *ErrorString) Error() string {
	return e.S
}

// ToJSONString converts error string to json string
func ToJSONString(e error) string {
	return "{\"str\":\"" + e.Error() + "\"}"
}
