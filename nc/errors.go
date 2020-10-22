package nc

import "fmt"

//ConflictError describes a conflict with the network state and requested changes
type ConflictError struct {
	reason string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict error: %s", e.reason)
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(ifname string) error {
	return &ConflictError{reason: "Link " + ifname + " exists"}
}
