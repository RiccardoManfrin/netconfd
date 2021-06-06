// Copyright (c) 2021, Athonet S.r.l. All rights reserved.
// riccardo.manfrin@athonet.com

package nc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"syscall"

	"github.com/riccardomanfrin/netconfd/logger"
	"github.com/vishvananda/netlink"
)

// PortRange represents rule sport/dport range.
type PortRange struct {
	Start uint16
	End   uint16
}

func (p PortRange) IsSingle() bool {
	return p.Start == p.End
}

func (p *PortRange) String() string {
	return fmt.Sprintf("%v-%v", p.Start, p.End)
}

func ParsePortRange(prange string) (error, PortRange) {
	var pl, pr uint16
	found, err := fmt.Sscanf(prange, "%v-%v", pl, pr)
	if found != 2 {
		return fmt.Errorf("Invalid port range %v", prange), PortRange{}
	}
	if err != nil {
		return err, PortRange{}
	}
	if pl < pr {
		return nil, PortRange{Start: pl, End: pr}
	} else {
		return nil, PortRange{Start: pr, End: pl}
	}
}

// Rule represents a netlink rule.
type Rule struct {
	ID                RuleID
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

//Print implements rule print
func (r *Rule) Print() string {
	data, _ := json.Marshal(r)
	return fmt.Sprintf("%v", string(data))
}

func ruleFormat(rule Rule) (netlink.Rule, error) {
	nlrule := netlink.Rule{}
	nlrule.Priority = rule.Priority
	nlrule.Family = rule.Family
	nlrule.Table = rule.Table
	nlrule.Mark = rule.Mark
	nlrule.Mask = rule.Mask
	nlrule.Tos = rule.Tos
	nlrule.TunID = rule.TunID
	nlrule.Goto = rule.Goto
	nlrule.Src = rule.Src
	nlrule.Dst = rule.Dst
	nlrule.Flow = rule.Flow
	nlrule.IifName = rule.IifName
	nlrule.OifName = rule.OifName
	nlrule.SuppressIfgroup = rule.SuppressIfgroup
	nlrule.SuppressPrefixlen = rule.SuppressPrefixlen
	nlrule.Invert = rule.Invert
	nlrule.Dport = &netlink.RulePortRange{
		Start: rule.Dport.Start,
		End:   rule.Dport.End,
	}
	nlrule.Sport = &netlink.RulePortRange{
		Start: rule.Sport.Start,
		End:   rule.Sport.End,
	}
	return nlrule, nil
}

//RuleID identifies a rule via MD5 of its content
type RuleID string

//RuleCreate create and add a new rule
func RuleCreate(rule Rule) (RuleID, error) {
	ruleid := RuleIDGet(rule)
	if isUnmanaged(UnmanagedID(rule.IifName), LINKTYPE) {
		logger.Log.Info(fmt.Sprintf("Skipping Unmanaged Link %v rule configuration", rule.IifName))
		return ruleid, NewUnmanagedLinkRuleCannotBeModifiedError(rule)
	}
	rule.ID = ruleid
	rules, err := RulesGet()
	if err != nil {
		return ruleid, err
	}
	for _, r := range rules {
		if r.ID == ruleid {
			return ruleid, NewRuleExistsConflictError(ruleid)
		}
	}

	nlrule, err := ruleFormat(rule)
	if err != nil {
		return ruleid, err
	}
	logger.Log.Debug(fmt.Sprintf("Creating rule %v", rule))
	err = netlink.RuleAdd(&nlrule)
	if err != nil {
		if err.(syscall.Errno) == syscall.EEXIST {
			logger.Log.Warning(fmt.Sprintf("Skipping rule %v creation: rule exists", ruleid))
		} else {
			return ruleid, mapNetlinkError(err, &rule)
		}
	}

	return ruleid, nil
}

func RuleIDGet(r Rule) RuleID {
	md := md5.New()
	serialized, err := json.Marshal(r)
	if err != nil {
		logger.Fatal(err)
		return RuleID("")
	}
	md.Write(serialized)
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
	return ncrules, nil
}

func ruleParse(rule netlink.Rule) (Rule, error) {
	r := Rule{
		Priority:          rule.Priority,
		Family:            rule.Family,
		Table:             rule.Table,
		Mark:              rule.Mark,
		Mask:              rule.Mask,
		Tos:               rule.Tos,
		TunID:             rule.TunID,
		Goto:              rule.Goto,
		Src:               rule.Src,
		Dst:               rule.Dst,
		Flow:              rule.Flow,
		IifName:           rule.IifName,
		OifName:           rule.OifName,
		SuppressIfgroup:   rule.SuppressIfgroup,
		SuppressPrefixlen: rule.SuppressPrefixlen,
		Invert:            rule.Invert,
	}
	if rule.Dport != nil {
		r.Dport = &PortRange{
			Start: rule.Dport.Start,
			End:   rule.Dport.End,
		}
	}
	if rule.Sport != nil {
		r.Sport = &PortRange{
			Start: rule.Sport.Start,
			End:   rule.Sport.End,
		}
	}

	return r, nil
}
