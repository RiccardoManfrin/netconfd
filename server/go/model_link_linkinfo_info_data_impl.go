package openapi

import "encoding/json"

type DiscriminatedInfoData struct {
	LinkLinkinfoInfoData
}

func (v DiscriminatedInfoData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.LinkLinkinfoInfoData)
}

func (v *DiscriminatedInfoData) UnmarshalJSON(src []byte) error {
	return json.Unmarshal(src, &v.LinkLinkinfoInfoData)
}

func (o LinkLinkinfoInfoData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Protocol != nil {
		toSerialize["protocol"] = o.Protocol
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Flags != nil {
		toSerialize["flags"] = o.Flags
	}
	if o.Mode != nil {
		toSerialize["mode"] = o.Mode
	}
	if o.Miimon != nil {
		toSerialize["miimon"] = o.Miimon
	}
	if o.Updelay != nil {
		toSerialize["updelay"] = o.Updelay
	}
	if o.Downdelay != nil {
		toSerialize["downdelay"] = o.Downdelay
	}
	if o.PeerNotifyDelay != nil {
		toSerialize["peer_notify_delay"] = o.PeerNotifyDelay
	}
	if o.UseCarrier != nil {
		toSerialize["use_carrier"] = o.UseCarrier
	}
	if o.ArpInterval != nil {
		toSerialize["arp_interval"] = o.ArpInterval
	}
	if o.ArpValidate != nil {
		toSerialize["arp_validate"] = o.ArpValidate
	}
	if o.ArpAllTargets != nil {
		toSerialize["arp_all_targets"] = o.ArpAllTargets
	}
	if o.PrimaryReselect != nil {
		toSerialize["primary_reselect"] = o.PrimaryReselect
	}
	if o.FailOverMac != nil {
		toSerialize["fail_over_mac"] = o.FailOverMac
	}
	if o.XmitHashPolicy != nil {
		toSerialize["xmit_hash_policy"] = o.XmitHashPolicy
	}
	if o.ResendIgmp != nil {
		toSerialize["resend_igmp"] = o.ResendIgmp
	}
	if o.AllSlavesActive != nil {
		toSerialize["all_slaves_active"] = o.AllSlavesActive
	}
	if o.MinLinks != nil {
		toSerialize["min_links"] = o.MinLinks
	}
	if o.LpInterval != nil {
		toSerialize["lp_interval"] = o.LpInterval
	}
	if o.PacketsPerSlave != nil {
		toSerialize["packets_per_slave"] = o.PacketsPerSlave
	}
	if o.AdLacpRate != nil {
		toSerialize["ad_lacp_rate"] = o.AdLacpRate
	}
	if o.AdSelect != nil {
		toSerialize["ad_select"] = o.AdSelect
	}
	if o.TlbDynamicLb != nil {
		toSerialize["tlb_dynamic_lb"] = o.TlbDynamicLb
	}
	if o.Local != nil {
		toSerialize["local"] = o.Local
	}
	if o.Remote != nil {
		toSerialize["remote"] = o.Remote
	}
	if o.Table != nil {
		toSerialize["table"] = *o.Table
	}
	return json.Marshal(toSerialize)
}
