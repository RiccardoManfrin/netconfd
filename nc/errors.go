package nc

import (
	"encoding/json"
	"fmt"
)

type ErrorCode int

const (
	//CONFLICT error type (inconsistency with respect to the existing state)
	CONFLICT ErrorCode = iota
	//SEMANTIC error type of the requested operation in the syntax or logical content
	SEMANTIC
	//UNKNOWN_TYPE error type (the value type is not recognized/supported)
	UNKNOWN_TYPE
	//RESERVED can be used for outer error enum cohexistence
	RESERVED = 1000
)

var errorCodeToString = map[ErrorCode]string{
	CONFLICT:     "Conflict Error",
	SEMANTIC:     "Semantic Error",
	UNKNOWN_TYPE: "UnknownType Error",
}

//GenericError describes a generic error of the library
type GenericError struct {
	//code error type
	Code ErrorCode `json:"code"`
	//reason describes the specific reason for the error
	Reason string `json:"reason"`
}

func (e *GenericError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//ConflictError describes a conflict with the network state and requested changes
type ConflictError GenericError

func (e *ConflictError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//SemanticError is a logical error on the content of the operation requested to be performed
type SemanticError GenericError

func (e *SemanticError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//UnknownTypeError is a logical error on the content of the operation requested to be performed
type UnknownTypeError GenericError

func (e *UnknownTypeError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewGenericSemanticError returns a generic semantic error
func NewGenericSemanticError() error {
	return &SemanticError{Code: SEMANTIC, Reason: "Generic Semantic Error"}
}

//NewBadAddressError returns a bad address error on link layer interfaces
func NewBadAddressError(c CIDRAddr) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Bad IP address " + c.String()}
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(linkID LinkID) error {
	return &ConflictError{Code: CONFLICT, Reason: "Link " + string(linkID) + " exists"}
}

//NonBondMasterLinkTypeError returns an error for non bond master link type
func NonBondMasterLinkTypeError(ifname LinkID) error {
	return &ConflictError{Code: SEMANTIC, Reason: "Master link interface " + string(ifname) + " is not a bond"}
}

//NewLinkNotFoundError returns a Not found error on link layer interfaces
func NewLinkNotFoundError(linkID LinkID) error {
	return &NotFoundError{Code: CONFLICT, Reason: "Link " + string(linkID) + " not found"}
}

//NewUnknownLinkKindError returns a UnknownTypeError error on link layer type interfaces
func NewUnknownLinkKindError(linkKind string) error {
	return &SemanticError{Code: SEMANTIC, Reason: "LinkKind " + string(linkKind) + " not known"}
}

//NotFoundError is a logical error on the content of the operation requested to be performed
type NotFoundError ConflictError

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found: %s", e.Reason)
}

//NewRouteByIDNotFoundError returns a Conflict error on link layer interfaces
func NewRouteByIDNotFoundError(routeid RouteID) error {
	return &NotFoundError{Code: CONFLICT, Reason: "Route ID" + string(routeid) + " did not match"}
}

//NewActiveSlaveIfaceNotFoundForActiveBackupBondError Returns an error if an active interface is not found for an Active-Backup type bond
func NewActiveSlaveIfaceNotFoundForActiveBackupBondError(bondIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Active Slave Iface not found for Active-Backup type bond " + string(bondIfname)}
}

//NewBackupSlaveIfaceFoundForNonActiveBackupBondError Returns an error if a backup interface is found for a non Active-Backup type bond
func NewBackupSlaveIfaceFoundForNonActiveBackupBondError(backupIfname LinkID, bondIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Backup Slave Iface " + string(backupIfname) + " found for non Active-Backup type bond " + string(bondIfname)}
}
