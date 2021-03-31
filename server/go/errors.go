package openapi

import (
	"gitlab.lan.athonet.com/core/netconfd/nc"
)

//NewGenericSemanticError returns a generic semantic error
func NewGenericSemanticError() error {

	return &nc.SemanticError{Code: nc.SEMANTIC, Reason: "Generic Semantic Error"}
}
