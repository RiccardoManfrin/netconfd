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

//NewBadAddressError returns a Conflict error on link layer interfaces
func NewBadAddressError(c CIDRAddr) error {
	return &SemanticError{reason: "Bad IP address " + c.String()}
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(linkID LinkID) error {
	return &ConflictError{reason: "Link " + string(linkID) + " exists"}
}

//NotFoundError is a logical error on the content of the operation requested to be performed
type NotFoundError ConflictError

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", e.reason)
}

//NewRouteByIDNotFoundError returns a Conflict error on link layer interfaces
func NewRouteByIDNotFoundError(routeid RouteID) error {
	return &ConflictError{reason: "Route ID" + string(routeid) + " did not match"}
}
