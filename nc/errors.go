package nc

import "fmt"

//ConflictError describes a conflict with the network state and requested changes
type ConflictError struct {
	reason string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict error: %s", e.reason)
}

//SemanticError is a logical error on the content of the operation requested to be performed
type SemanticError ConflictError

func (e *SemanticError) Error() string {
	return fmt.Sprintf("Semantic error: %s", e.reason)
}

//NewBadLinkNameError returns a Conflict error on link layer interfaces
func NewBadLinkNameError(ifname string) error {
	return &SemanticError{reason: "Link name" + ifname + " is unacceptable"}
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(ifname string) error {
	return &ConflictError{reason: "Link " + ifname + " exists"}
}
