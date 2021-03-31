package openapi

import (
	"encoding/json"
)

type DiscriminatedLinkInfo struct {
	LinkLinkinfo
}

func huntVlanInfoDataParams(v *LinkLinkinfo) error {
	if _, ok := v.InfoData.GetIdOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("id", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetProtocolOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("protocol", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetFlagsOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("flags", v.GetInfoKind())
	}
	return nil
}

func huntGreInfoDataParams(v *LinkLinkinfo) error {
	if _, ok := v.InfoData.GetLocalOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("local", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetRemoteOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("remote", v.GetInfoKind())
	}
	return nil
}

func huntBondInfoDataParams(v *LinkLinkinfo) error {
	if _, ok := v.InfoData.GetModeOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("mode", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetMiimonOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("miimon", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetUpdelayOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("updelay", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetDowndelayOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("downdelay", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetPeerNotifyDelayOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("peer_notify_delay", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetUseCarrierOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("use_carrier", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetArpIntervalOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("arp_interval", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetArpValidateOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("arp_validate", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetArpAllTargetsOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("arp_all_targets", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetPrimaryReselectOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("primary_reselect", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetFailOverMacOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("fail_over_mac", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetXmitHashPolicyOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("xmit_hash_policy", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetResendIgmpOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("resend_igmp", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetAllSlavesActiveOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("all_slaves_active", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetMinLinksOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("min_links", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetLpIntervalOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("lp_interval", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetPacketsPerSlaveOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("packets_per_slave", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetAdLacpRateOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("ad_lacp_rate", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetAdSelectOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("ad_select", v.GetInfoKind())
	}
	if _, ok := v.InfoData.GetTlbDynamicLbOk(); ok {
		return NewAttributeDoesntBelongToLinkKindSemanticError("tlb_dynamic_lb", v.GetInfoKind())
	}
	return nil
}

func (v *LinkLinkinfo) Validate() error {
	if v.InfoKind == nil {
		return nil
	}
	infoKind := v.GetInfoKind()
	switch infoKind {
	case "bond":
		{
			if v.InfoData == nil {
				return NewMissingRequiredAttributeForLinkKindSemanticError("info_data", infoKind)
			}
			if err := huntVlanInfoDataParams(v); err != nil {
				return err
			}
			if err := huntGreInfoDataParams(v); err != nil {
				return err
			}
		}
	case "gre":
		{
			if v.InfoData == nil {
				return NewMissingRequiredAttributeForLinkKindSemanticError("info_data", infoKind)
			}
			if err := huntVlanInfoDataParams(v); err != nil {
				return err
			}
			if err := huntBondInfoDataParams(v); err != nil {
				return err
			}
		}
	case "vlan":
		{
			if v.InfoData == nil {
				return NewMissingRequiredAttributeForLinkKindSemanticError("info_data", infoKind)
			}
			if err := huntBondInfoDataParams(v); err != nil {
				return err
			}
			if err := huntGreInfoDataParams(v); err != nil {
				return err
			}
		}
	case "device",
		"bridge",
		"dummy",
		"ppp",
		"tun",
		"tap":
		{
			if v.InfoData != nil {
				return NewAttributeDoesntBelongToLinkKindSemanticError("info_data", infoKind)
			}
		}
	}
	return nil
}

func (v *DiscriminatedLinkInfo) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &v.LinkLinkinfo)
	if err != nil {
		return err
	}
	err = v.LinkLinkinfo.Validate()
	return err
}
