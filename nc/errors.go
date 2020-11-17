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

//NewGenericSemanticError returns a generic semantic error
func NewGenericSemanticError() error {
	return &SemanticError{reason: "Generic Semantic Error"}
}

//NewBadAddressError returns a bad address error on link layer interfaces
func NewBadAddressError(c CIDRAddr) error {
	return &SemanticError{reason: "Bad IP address " + c.String()}
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(linkID LinkID) error {
	return &ConflictError{reason: "Link " + string(linkID) + " exists"}
}

//UnknownTypeError is a logical error on the content of the operation requested to be performed
type UnknownTypeError SemanticError

func (e *UnknownTypeError) Error() string {
	return fmt.Sprintf("Unknown type: %s", e.reason)
}

//NewUnknownLinkKindError returns a UnknownTypeError error on link layer type interfaces
func NewUnknownLinkKindError(linkKind string) error {
	return &SemanticError{reason: "LinkKind " + string(linkKind) + " not known"}
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
