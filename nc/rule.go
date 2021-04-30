package nc

import (
	"crypto/md5"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// PortRange represents rule sport/dport range.
type PortRange struct {
	Start uint16
	End   uint16
}

// Rule represents a netlink rule.
type Rule struct {
	Priority          int
	Family            int
	Table             int
	Mark              int
	Mask              int
	Tos               uint
	TunID             uint
	Goto              int
	Src               *net.IPNet
	Dst               *net.IPNet
	Flow              int
	IifName           string
	OifName           string
	SuppressIfgroup   int
	SuppressPrefixlen int
	Invert            bool
	Dport             *PortRange
	Sport             *PortRange
}

//RuleID identifies a rule via MD5 of its content
type RuleID string

//RuleCreate create and add a new rule
func RuleCreate(rule Rule) (RuleID, error) {
	id := RuleID("IMPLEMENT")
	return id, NewUnsupportedError("TODO")
}

func RuleIDGet(r Rule) RuleID {
	md := md5.New()
	md.Write([]byte(r.IifName))
	md.Write([]byte(r.OifName))
	return RuleID(fmt.Sprintf("%x", md.Sum(nil)))
}

//RuleDelete deletes a rule by ID
func RuleDelete(ruleid RuleID) error {
	return NewUnsupportedError("TODO")
}

//RuleGet Returns a rule if it exists
func RuleGet(_ruleID RuleID) (Rule, error) {
	r := Rule{}
	return r, NewUnsupportedError("TODO")
}

//RulesGet returns the array of rules
func RulesGet() ([]Rule, error) {
	rules, err := netlink.RuleListFiltered(netlink.FAMILY_ALL, nil, 0)
	if err != nil {
		return nil, mapNetlinkError(err, nil)
	}
	ncrules := make([]Rule, len(rules))
	for i, r := range rules {
		ncrules[i], err = ruleParse(r)
		if err != nil {
			return ncrules, err
		}
	}
	return ncrules, NewUnsupportedError("TODO")
}

func ruleParse(rule netlink.Rule) (Rule, error) {
	r := Rule{}
	return r, NewUnsupportedError("TODO")
}
