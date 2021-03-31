package openapi

import (
	"fmt"
	"runtime/debug"

	"gitlab.lan.athonet.com/core/netconfd/nc"
)

//NewAttributeDoesntBelongToLinkKindSemanticError returns an error for an attribute not belonging to the discriminated object
func NewAttributeDoesntBelongToLinkKindSemanticError(attribKey string, infoKind string) error {
	if nc.NetconfdDebugTrace {
		debug.PrintStack()
	}
	return &nc.SemanticError{Code: nc.SEMANTIC, Reason: fmt.Sprintf("Attribute %s doesn't belong to Link Kind %s", attribKey, infoKind)}
}

//NewMissingRequiredAttributeForLinkKindSemanticError returns an error for a missing attribute
func NewMissingRequiredAttributeForLinkKindSemanticError(attribKey string, infoKind string) error {
	if nc.NetconfdDebugTrace {
		debug.PrintStack()
	}
	return &nc.SemanticError{Code: nc.SEMANTIC, Reason: fmt.Sprintf("Missing required attribute %s for Link Kind %s", attribKey, infoKind)}
}
