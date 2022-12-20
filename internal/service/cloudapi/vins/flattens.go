/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Terraform DECORT provider - manage resources provided by DECORT (Digital Energy Cloud
Orchestration Technology) with Terraform by Hashicorp.

Source code: https://github.com/rudecs/terraform-provider-decort

Please see README.md to learn where to place source code so that it
builds seamlessly.

Documentation: https://github.com/rudecs/terraform-provider-decort/wiki
*/

package vins

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func flattenMGMT(mgmt *VNFConfigMGMT) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"ip_addr":  mgmt.IPAddr,
		"password": mgmt.Password,
		"ssh_key":  mgmt.SSHKey,
		"user":     mgmt.User,
	}
	res = append(res, temp)
	return res
}

func flattenResources(resources *VNFConfigResources) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":      resources.CPU,
		"ram":      resources.RAM,
		"stack_id": resources.StackID,
		"uuid":     resources.UUID,
	}
	res = append(res, temp)
	return res
}

func flattenConfig(config VNFConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"mgmt":      flattenMGMT(&config.MGMT),
		"resources": flattenResources(&config.Resources),
	}
	res = append(res, temp)
	return res
}

func flattenQOS(qos QOS) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"e_rate":   qos.ERate,
		"guid":     qos.GUID,
		"in_brust": qos.InBurst,
		"in_rate":  qos.InRate,
	}
	res = append(res, temp)
	return res
}

func flattenInterfaces(interfaces VNFInterfaceList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)

	for _, vnfInterface := range interfaces {
		temp := map[string]interface{}{
			"conn_id":      vnfInterface.ConnID,
			"conn_type":    vnfInterface.ConnType,
			"def_gw":       vnfInterface.DefGW,
			"flipgroup_id": vnfInterface.FlipGroupID,
			"guid":         vnfInterface.GUID,
			"ip_address":   vnfInterface.IPAddress,
			"listen_ssh":   vnfInterface.ListenSSH,
			"mac":          vnfInterface.MAC,
			"name":         vnfInterface.Name,
			"net_id":       vnfInterface.NetID,
			"net_mask":     vnfInterface.NetMask,
			"net_type":     vnfInterface.NetType,
			"pci_slot":     vnfInterface.PCISlot,
			"qos":          flattenQOS(vnfInterface.QOS),
			"target":       vnfInterface.Target,
			"type":         vnfInterface.Type,
			"vnfs":         vnfInterface.VNFS,
		}
		res = append(res, temp)
	}

	return res
}

func flattenVNFDev(vnfDev VNFDev) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"_ckey":          vnfDev.CKey,
		"account_id":     vnfDev.AccountID,
		"capabilities":   vnfDev.Capabilities,
		"config":         flattenConfig(vnfDev.Config), //in progress
		"config_saved":   vnfDev.ConfigSaved,
		"custom_pre_cfg": vnfDev.CustomPreConfig,
		"desc":           vnfDev.Description,
		"gid":            vnfDev.GID,
		"guid":           vnfDev.GUID,
		"vnf_id":         vnfDev.ID,
		"interfaces":     flattenInterfaces(vnfDev.Interfaces),
		"lock_status":    vnfDev.LockStatus,
		"milestones":     vnfDev.Milestones,
		"vnf_name":       vnfDev.Name,
		"status":         vnfDev.Status,
		"tech_status":    vnfDev.TechStatus,
		"type":           vnfDev.Type,
		"vins":           vnfDev.VINS,
	}

	res = append(res, temp)

	return res
}

func flattenComputes(computes VINSComputeList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, compute := range computes {
		temp := map[string]interface{}{
			"compute_id":   compute.ID,
			"compute_name": compute.Name,
		}
		res = append(res, temp)
	}

	return res
}

func flattenReservations(reservations ReservationList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, reservation := range reservations {
		temp := map[string]interface{}{
			"client_type": reservation.ClientType,
			"desc":        reservation.Description,
			"domainname":  reservation.DomainName,
			"hostname":    reservation.HostName,
			"ip":          reservation.IP,
			"mac":         reservation.MAC,
			"type":        reservation.Type,
			"vm_id":       reservation.VMID,
		}
		res = append(res, temp)
	}

	return res
}

func flattenDHCPConfig(config DHCPConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"default_gw":   config.DefaultGW,
		"dns":          config.DNS,
		"ip_end":       config.IPEnd,
		"ip_start":     config.IPStart,
		"lease":        config.Lease,
		"netmask":      config.Netmask,
		"network":      config.Network,
		"reservations": flattenReservations(config.Reservations),
	}
	res = append(res, temp)

	return res
}

func flattenPrimary(primary DevicePrimary) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"dev_id":  primary.DevID,
		"iface01": primary.IFace01,
		"iface02": primary.IFace02,
	}
	res = append(res, temp)

	return res
}

func flattenDevices(devices Devices) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"primary": flattenPrimary(devices.Primary),
	}
	res = append(res, temp)

	return res
}

func flattenDHCP(dhcp DHCP) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"_ckey":        dhcp.CKey,
		"account_id":   dhcp.AccountID,
		"config":       flattenDHCPConfig(dhcp.Config),
		"created_time": dhcp.CreatedTime,
		"devices":      flattenDevices(dhcp.Devices),
		"gid":          dhcp.GID,
		"guid":         dhcp.GUID,
		"dhcp_id":      dhcp.ID,
		"lock_status":  dhcp.LockStatus,
		"milestones":   dhcp.Milestones,
		"owner_id":     dhcp.OwnerID,
		"owner_type":   dhcp.OwnerType,
		"pure_virtual": dhcp.PureVirtual,
		"status":       dhcp.Status,
		"tech_status":  dhcp.TechStatus,
		"type":         dhcp.Type,
	}
	res = append(res, temp)

	return res
}

func flattenGWConfig(config GWConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"default_gw":  config.DefaultGW,
		"ext_net_id":  config.ExtNetID,
		"ext_net_ip":  config.ExtNetIP,
		"ext_netmask": config.ExtNetMask,
		"qos":         flattenQOS(config.QOS),
	}
	res = append(res, temp)

	return res
}

func flattenGW(gw GW) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"_ckey":        gw.CKey,
		"account_id":   gw.AccountID,
		"config":       flattenGWConfig(gw.Config),
		"created_time": gw.CreatedTime,
		"devices":      flattenDevices(gw.Devices),
		"gid":          gw.GID,
		"guid":         gw.GUID,
		"gw_id":        gw.ID,
		"lock_status":  gw.LockStatus,
		"milestones":   gw.Milestones,
		"owner_id":     gw.OwnerID,
		"owner_type":   gw.OwnerType,
		"pure_virtual": gw.PureVirtual,
		"status":       gw.Status,
		"tech_status":  gw.TechStatus,
		"type":         gw.Type,
	}
	res = append(res, temp)

	return res
}

func flattenRules(rules ListNATRules) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmp := map[string]interface{}{
			"rule_id":           rule.ID,
			"local_ip":          rule.LocalIP,
			"local_port":        rule.LocalPort,
			"protocol":          rule.Protocol,
			"public_port_end":   rule.PublicPortEnd,
			"public_port_start": rule.PublicPortStart,
			"vm_id":             rule.VMID,
			"vm_name":           rule.VMName,
		}
		res = append(res, tmp)
	}
	return res
}

func flattenNATConfig(config NATConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"net_mask": config.NetMask,
		"network":  config.Network,
		"rules":    flattenRules(config.Rules),
	}
	res = append(res, temp)

	return res

}

func flattenNAT(nat NAT) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"_ckey":        nat.CKey,
		"account_id":   nat.AccountID,
		"created_time": nat.CreatedTime,
		"config":       flattenNATConfig(nat.Config),
		"devices":      flattenDevices(nat.Devices),
		"gid":          nat.GID,
		"guid":         nat.GUID,
		"nat_id":       nat.ID,
		"lock_status":  nat.LockStatus,
		"milestones":   nat.Milestones,
		"owner_id":     nat.OwnerID,
		"owner_type":   nat.OwnerType,
		"pure_virtual": nat.PureVirtual,
		"status":       nat.Status,
		"tech_status":  nat.TechStatus,
		"type":         nat.Type,
	}
	res = append(res, temp)

	return res
}

func flattenVNFS(vnfs VNFS) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"dhcp": flattenDHCP(vnfs.DHCP),
		"gw":   flattenGW(vnfs.GW),
		"nat":  flattenNAT(vnfs.NAT),
	}
	res = append(res, temp)

	return res
}

func flattenRuleBlock(rules ListNATRules) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmp := map[string]interface{}{
			"int_ip":         rule.LocalIP,
			"int_port":       rule.LocalPort,
			"ext_port_start": rule.PublicPortStart,
			"ext_port_end":   rule.PublicPortEnd,
			"proto":          rule.Protocol,
			"rule_id":        rule.ID,
		}
		res = append(res, tmp)
	}
	return res
}

func flattenVins(d *schema.ResourceData, vins VINSDetailed) {
	d.Set("vins_id", vins.ID)
	d.Set("vnf_dev", flattenVNFDev(vins.VNFDev))
	d.Set("_ckey", vins.CKey)
	d.Set("account_id", vins.AccountID)
	d.Set("account_name", vins.AccountName)
	d.Set("computes", flattenComputes(vins.Computes))
	d.Set("default_gw", vins.DefaultGW)
	d.Set("default_qos", flattenQOS(vins.DefaultQOS))
	d.Set("desc", vins.Description)
	d.Set("gid", vins.GID)
	d.Set("guid", vins.GUID)
	d.Set("lock_status", vins.LockStatus)
	d.Set("manager_id", vins.ManagerID)
	d.Set("manager_type", vins.ManagerType)
	d.Set("milestones", vins.Milestones)
	d.Set("name", vins.Name)
	d.Set("net_mask", vins.NetMask)
	d.Set("network", vins.Network)
	d.Set("pre_reservations_num", vins.PreReservaionsNum)
	d.Set("redundant", vins.Redundant)
	d.Set("rg_id", vins.RGID)
	d.Set("rg_name", vins.RGName)
	d.Set("sec_vnf_dev_id", vins.SecVNFDevID)
	d.Set("status", vins.Status)
	d.Set("user_managed", vins.UserManaged)
	d.Set("vnfs", flattenVNFS(vins.VNFS))
	d.Set("vxlan_id", vins.VXLanID)
	d.Set("nat_rule", flattenRuleBlock(vins.VNFS.NAT.Config.Rules))
}

func flattenVinsData(d *schema.ResourceData, vins VINSDetailed) {
	d.Set("vins_id", vins.ID)
	d.Set("vnf_dev", flattenVNFDev(vins.VNFDev))
	d.Set("_ckey", vins.CKey)
	d.Set("account_id", vins.AccountID)
	d.Set("account_name", vins.AccountName)
	d.Set("computes", flattenComputes(vins.Computes))
	d.Set("default_gw", vins.DefaultGW)
	d.Set("default_qos", flattenQOS(vins.DefaultQOS))
	d.Set("desc", vins.Description)
	d.Set("gid", vins.GID)
	d.Set("guid", vins.GUID)
	d.Set("lock_status", vins.LockStatus)
	d.Set("manager_id", vins.ManagerID)
	d.Set("manager_type", vins.ManagerType)
	d.Set("milestones", vins.Milestones)
	d.Set("name", vins.Name)
	d.Set("net_mask", vins.NetMask)
	d.Set("network", vins.Network)
	d.Set("pre_reservations_num", vins.PreReservaionsNum)
	d.Set("redundant", vins.Redundant)
	d.Set("rg_id", vins.RGID)
	d.Set("rg_name", vins.RGName)
	d.Set("sec_vnf_dev_id", vins.SecVNFDevID)
	d.Set("status", vins.Status)
	d.Set("user_managed", vins.UserManaged)
	d.Set("vnfs", flattenVNFS(vins.VNFS))
	d.Set("vxlan_id", vins.VXLanID)
}

func flattenVinsAudits(auidts VINSAuditsList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, audit := range auidts {
		temp := map[string]interface{}{
			"call":          audit.Call,
			"response_time": audit.ResponseTime,
			"statuscode":    audit.StatusCode,
			"timestamp":     audit.Timestamp,
			"user":          audit.User,
		}
		res = append(res, temp)
	}

	return res
}

func flattenVinsExtNetList(extNetList ExtNetList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, extNet := range extNetList {
		temp := map[string]interface{}{
			"default_gw":  extNet.DefaultGW,
			"ext_net_id":  extNet.ExtNetID,
			"ip":          extNet.IP,
			"prefix_len":  extNet.PrefixLen,
			"status":      extNet.Status,
			"tech_status": extNet.TechStatus,
		}
		res = append(res, temp)
	}

	return res
}

func flattenVinsIpList(ips IPList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, ip := range ips {
		temp := map[string]interface{}{
			"client_type": ip.ClientType,
			"domainname":  ip.DomainName,
			"hostname":    ip.HostName,
			"ip":          ip.IP,
			"mac":         ip.MAC,
			"type":        ip.Type,
			"vm_id":       ip.VMID,
		}
		res = append(res, temp)
	}

	return res
}

func flattenVinsList(vl VINSList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, v := range vl {
		temp := map[string]interface{}{
			"account_id":   v.AccountID,
			"account_name": v.AccountName,
			"created_by":   v.CreatedBy,
			"created_time": v.CreatedTime,
			"deleted_by":   v.DeletedBy,
			"deleted_time": v.DeletedTime,
			"external_ip":  v.ExternalIP,
			"vins_id":      v.ID,
			"vins_name":    v.Name,
			"network":      v.Network,
			"rg_id":        v.RGID,
			"rg_name":      v.RGName,
			"status":       v.Status,
			"updated_by":   v.UpdatedBy,
			"updated_time": v.UpdatedTime,
			"vxlan_id":     v.VXLANID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenVinsNatRuleList(natRules NATRuleList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, natRule := range natRules {
		temp := map[string]interface{}{
			"id":                natRule.ID,
			"local_ip":          natRule.LocalIP,
			"local_port":        natRule.LocalPort,
			"protocol":          natRule.Protocol,
			"public_port_end":   natRule.PublicPortEnd,
			"public_port_start": natRule.PublicPortStart,
			"vm_id":             natRule.VMID,
			"vm_name":           natRule.VMName,
		}
		res = append(res, temp)
	}

	return res
}
