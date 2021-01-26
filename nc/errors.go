package nc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

//ErrorCode describes the error type via enumeration
type ErrorCode int

const (
	//CONFLICT error type (inconsistency with respect to the existing state)
	CONFLICT ErrorCode = iota
	//NOT_FOUND error types encodes a restful resource not found by its ID
	NOT_FOUND
	//SEMANTIC error type of the requested operation in the syntax or logical content
	SEMANTIC
	//SYNTAX error type is for synctactical errors
	SYNTAX
	//UNKNOWN_TYPE error type (the value type is not recognized/supported)
	UNKNOWN_TYPE
	//UNEXPECTED_CORNER_CASE error type describes an error that was not meant to appear
	UNEXPECTED_CORNER_CASE
	//RESERVED can be used for outer error enum cohexistence
	RESERVED = 1000
)

var errorCodeToString = map[ErrorCode]string{
	CONFLICT:     "Conflict Error",
	SEMANTIC:     "Semantic Error",
	SYNTAX:       "Syntax Error",
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

//NewGenericError returns a generic error
func NewGenericError(err error) error {
	return &GenericError{Code: UNKNOWN_TYPE, Reason: fmt.Sprintf("Generic uncharted error :%v", err.Error())}
}

//NewGenericErrorWithReason returns a generic semantic error
func NewGenericErrorWithReason(reason string) error {
	return &GenericError{Code: UNKNOWN_TYPE, Reason: reason}
}

//SemanticError is a logical error on the content of the operation requested to be performed
type SemanticError GenericError

func (e *SemanticError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewGenericSemanticError returns a generic semantic error
func NewGenericSemanticError() error {
	return &SemanticError{Code: SEMANTIC, Reason: "Generic Semantic Error"}
}

//NewUnknownLinkKindError returns a SemanticError error on link layer type interfaces
func NewUnknownLinkKindError(linkKind string) error {
	return &SemanticError{Code: SEMANTIC, Reason: "LinkKind " + string(linkKind) + " not known"}
}

//NewBadAddressError returns a bad address error on link layer interfaces
func NewBadAddressError(c CIDRAddr) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Bad IP address " + c.String()}
}

//NewEINVALError returns a bad address error on link layer interfaces
func NewEINVALError() error {
	return &SemanticError{Code: SEMANTIC, Reason: "Syscall EINVAL error (check dmesg)"}
}

//NewActiveSlaveIfaceNotFoundForActiveBackupBondError Returns an error if an active interface is not found for an Active-Backup type bond
func NewActiveSlaveIfaceNotFoundForActiveBackupBondError(bondIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Active Slave Iface not found for Active-Backup type bond " + string(bondIfname)}
}

//NewParentLinkNotFoundForVlan returns a Not found error on link layer interfaces
func NewParentLinkNotFoundForVlan(ifname LinkID, parentIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Parent Link " + string(parentIfname) + " not found for Vlan Link " + string(ifname)}
}

//NewMultipleActiveSlaveIfacesFoundForActiveBackupBondError Returns an error if an active interface is not found for an Active-Backup type bond
func NewMultipleActiveSlaveIfacesFoundForActiveBackupBondError(bondIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Multiple Active Slave Ifaces found for Active-Backup type bond " + string(bondIfname)}
}

//NewBackupSlaveIfaceFoundForNonActiveBackupBondError Returns an error if a backup interface is found for a non Active-Backup type bond
func NewBackupSlaveIfaceFoundForNonActiveBackupBondError(backupIfname LinkID, bondIfname LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Backup Slave Iface " + string(backupIfname) + " found for non Active-Backup type bond " + string(bondIfname)}
}

//NewRouteLinkDeviceNotFoundError describes a link device not found for a route to create
func NewRouteLinkDeviceNotFoundError(routeID RouteID, linkID LinkID) error {
	return &SemanticError{Code: SEMANTIC, Reason: "Route " + string(routeID) + " Link Device " + string(linkID) + " not found"}
}

//NewENETUNREACHError returns a network unreachable error
func NewENETUNREACHError(r Route) error {
	return &SemanticError{Code: SEMANTIC, Reason: fmt.Sprintf("Got ENETUNREACH error: network is not reachable for route %v", r)}
}

//NewTooManyDNSServersError describes a link device not found for a route to create
func NewTooManyDNSServersError() error {
	return &SemanticError{Code: SEMANTIC, Reason: "More than two entries found in DNS servers config"}
}

//SyntaxError is a logical error on the content of the operation requested to be performed
type SyntaxError GenericError

func (e *SyntaxError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewInvalidIPAddressError Returns an error if a backup interface is found for a non Active-Backup type bond
func NewInvalidIPAddressError(addr string) error {
	return &SyntaxError{Code: SYNTAX, Reason: "Invalid IP Address/Network  " + addr}
}

//UnknownTypeError is a logical error on the content of the operation requested to be performed
type UnknownTypeError GenericError

func (e *UnknownTypeError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewLinkUnknownFlagTypeError returns a Conflict error on link layer interfaces
func NewLinkUnknownFlagTypeError(flag LinkFlag) error {
	return &UnknownTypeError{Code: UNKNOWN_TYPE, Reason: "Link Flag Type" + string(flag) + " unknown/unsupported"}
}

//ConflictError describes a conflict with the network state and requested changes
type ConflictError GenericError

func (e *ConflictError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewLinkExistsConflictError returns a Conflict error on link layer interfaces
func NewLinkExistsConflictError(linkID LinkID) error {
	return &ConflictError{Code: CONFLICT, Reason: "Link " + string(linkID) + " exists"}
}

//NewLinkDeviceDoesNotExistError returns a Conflict error on link layer interfaces
func NewLinkDeviceDoesNotExistError(linkID LinkID) error {
	return &ConflictError{Code: CONFLICT, Reason: "Link device " + string(linkID) + " does not exist"}
}

//NewNonBondMasterLinkTypeError returns an error for non bond master link type
func NewNonBondMasterLinkTypeError(ifname LinkID) error {
	return &ConflictError{Code: CONFLICT, Reason: "Master link interface " + string(ifname) + " is not a bond"}
}

//NewCannotStopDHCPError returns an error for DHCP related stop errors
func NewCannotStopDHCPError(ifname LinkID, e error) error {
	return &ConflictError{Code: CONFLICT, Reason: fmt.Sprintf("Failed to stop DHCP for interface %v: %v", string(ifname), e)}
}

//NewCannotStartDHCPError returns an error for DHCP related stop errors
func NewCannotStartDHCPError(ifname LinkID, e error) error {
	return &ConflictError{Code: CONFLICT, Reason: fmt.Sprintf("Failed to start DHCP for interface %v: %v", string(ifname), e)}
}

//NewCannotStatusDHCPError returns an error for DHCP related status errors
func NewCannotStatusDHCPError(ifname LinkID, e error) error {
	return &ConflictError{Code: CONFLICT, Reason: fmt.Sprintf("Failed to get DHCP status for interface %v: %v", string(ifname), e)}
}

//NewDHCPAlreadyRunningConflictError returns an error for DHCP that is requested for an interface where it's already running
func NewDHCPAlreadyRunningConflictError(ifname LinkID) error {
	return &ConflictError{Code: CONFLICT, Reason: fmt.Sprintf("DHCP is alreay running for interface %v", string(ifname))}
}

//NewEPERMError returns a missing permissions error
func NewEPERMError(context interface{}) error {
	return &ConflictError{
		Code: CONFLICT,
		Reason: fmt.Sprintf("Got EPERM error: insufficient permissions to perform action on %v: %v",
			reflect.TypeOf(context),
			context)}
}

//NewRouteExistsConflictError returns a Conflict error on link layer interfaces
func NewRouteExistsConflictError(routeID RouteID) error {
	return &ConflictError{Code: CONFLICT, Reason: "Route " + string(routeID) + " exists"}
}

//NotFoundError is a logical error on the content of the operation requested to be performed
type NotFoundError ConflictError

func (e *NotFoundError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewLinkNotFoundError returns a Not found error on link layer interfaces
func NewLinkNotFoundError(linkID LinkID) error {
	return &NotFoundError{Code: NOT_FOUND, Reason: "Link " + string(linkID) + " not found"}
}

//NewRouteByIDNotFoundError returns a Not found error on link layer interfaces
func NewRouteByIDNotFoundError(routeid RouteID) error {
	return &NotFoundError{Code: NOT_FOUND, Reason: "Route ID " + string(routeid) + " not found"}
}

//NewDHCPRunningNotFoundError returns a Not found error on link layer interfaces not managed by DHCP
func NewDHCPRunningNotFoundError(linkID LinkID) error {
	return &NotFoundError{Code: NOT_FOUND, Reason: "DHCP for Link ID " + string(linkID) + " not found"}
}

//NewDNSServerNotFoundError returns a Not found error on link layer interfaces not managed by DHCP
func NewDNSServerNotFoundError(dnsID DnsID) error {
	return &NotFoundError{Code: NOT_FOUND, Reason: "DNS ID " + string(dnsID) + " not found"}
}

//UnexpetecdCornerCaseError is fundamentally an implementation error catch exception
//It makes explitic to developer that he did not think of a case that instead happened
type UnexpetecdCornerCaseError GenericError

func (e *UnexpetecdCornerCaseError) Error() string {
	strerr, _ := json.Marshal(*e)
	return string(strerr)
}

//NewUnexpectedCornerCaseError returns a Conflict error on link layer interfaces
func NewUnexpectedCornerCaseError(reason string) error {
	return &UnexpetecdCornerCaseError{Code: UNEXPECTED_CORNER_CASE, Reason: reason}
}
