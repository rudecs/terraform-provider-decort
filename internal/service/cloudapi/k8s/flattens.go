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

package k8s

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

func flattenAclList(aclList ACLList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, acl := range aclList {
		temp := map[string]interface{}{
			"explicit":      acl.Explicit,
			"guid":          acl.GUID,
			"right":         acl.Right,
			"status":        acl.Status,
			"type":          acl.Type,
			"user_group_id": acl.UserGroupID,
		}
		res = append(res, temp)
	}
	return res
}

func flattenAcl(acl ACLGroup) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"account_acl": flattenAclList(acl.AccountACL),
		"k8s_acl":     flattenAclList(acl.K8SACL),
		"rg_acl":      flattenAclList(acl.RGACL),
	}

	res = append(res, temp)
	return res
}

func flattenInterfaces(interfaces []kvmvm.InterfaceRecord) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, interfaceCompute := range interfaces {
		temp := map[string]interface{}{
			"def_gw":     interfaceCompute.DefaultGW,
			"ip_address": interfaceCompute.IPAddress,
		}
		res = append(res, temp)
	}

	return res
}

func flattenDetailedInfo(detailedInfoList DetailedInfoList, computes []kvmvm.ComputeGetResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if computes != nil {
		for i, detailedInfo := range detailedInfoList {
			temp := map[string]interface{}{
				"compute_id":           detailedInfo.ID,
				"name":                 detailedInfo.Name,
				"status":               detailedInfo.Status,
				"tech_status":          detailedInfo.TechStatus,
				"interfaces":           flattenInterfaces(computes[i].Interfaces),
				"natable_vins_ip":      computes[i].NatableVinsIP,
				"natable_vins_network": computes[i].NatableVinsNet,
			}
			res = append(res, temp)
		}
	} else {
		for _, detailedInfo := range detailedInfoList {
			temp := map[string]interface{}{
				"compute_id":  detailedInfo.ID,
				"name":        detailedInfo.Name,
				"status":      detailedInfo.Status,
				"tech_status": detailedInfo.TechStatus,
			}
			res = append(res, temp)
		}
	}

	return res
}

func flattenMasterGroup(mastersGroup MasterGroup, masters []kvmvm.ComputeGetResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"cpu":           mastersGroup.CPU,
		"detailed_info": flattenDetailedInfo(mastersGroup.DetailedInfo, masters),
		"disk":          mastersGroup.Disk,
		"master_id":     mastersGroup.ID,
		"name":          mastersGroup.Name,
		"num":           mastersGroup.Num,
		"ram":           mastersGroup.RAM,
	}

	res = append(res, temp)
	return res
}

func flattenK8sGroup(k8SGroupList K8SGroupList, workers []kvmvm.ComputeGetResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, k8sGroup := range k8SGroupList {
		temp := map[string]interface{}{
			"annotations":   k8sGroup.Annotations,
			"cpu":           k8sGroup.CPU,
			"detailed_info": flattenDetailedInfo(k8sGroup.DetailedInfo, workers),
			"disk":          k8sGroup.Disk,
			"guid":          k8sGroup.GUID,
			"id":            k8sGroup.ID,
			"labels":        k8sGroup.Labels,
			"name":          k8sGroup.Name,
			"num":           k8sGroup.Num,
			"ram":           k8sGroup.RAM,
			"taints":        k8sGroup.Taints,
		}

		res = append(res, temp)
	}
	return res
}

func flattenK8sGroups(k8sGroups K8SGroups, masters []kvmvm.ComputeGetResp, workers []kvmvm.ComputeGetResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"masters": flattenMasterGroup(k8sGroups.Masters, masters),
		"workers": flattenK8sGroup(k8sGroups.Workers, workers),
	}
	res = append(res, temp)
	return res
}

func flattenK8sData(d *schema.ResourceData, k8s K8SRecord, masters []kvmvm.ComputeGetResp, workers []kvmvm.ComputeGetResp) {
	d.Set("acl", flattenAcl(k8s.ACL))
	d.Set("account_id", k8s.AccountID)
	d.Set("account_name", k8s.AccountName)
	d.Set("bservice_id", k8s.BServiceID)
	d.Set("k8sci_id", k8s.CIID)
	d.Set("created_by", k8s.CreatedBy)
	d.Set("created_time", k8s.CreatedTime)
	d.Set("deleted_by", k8s.DeletedBy)
	d.Set("deleted_time", k8s.DeletedTime)
	d.Set("k8s_ci_name", k8s.K8CIName)
	d.Set("masters", flattenMasterGroup(k8s.K8SGroups.Masters, masters))
	d.Set("workers", flattenK8sGroup(k8s.K8SGroups.Workers, workers))
	d.Set("lb_id", k8s.LBID)
	d.Set("name", k8s.Name)
	d.Set("rg_id", k8s.RGID)
	d.Set("rg_name", k8s.RGName)
	d.Set("status", k8s.Status)
	d.Set("tech_status", k8s.TechStatus)
	d.Set("updated_by", k8s.UpdatedBy)
	d.Set("updated_time", k8s.UpdatedTime)
}

func flattenServiceAccount(serviceAccount ServiceAccount) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	temp := map[string]interface{}{
		"guid":     serviceAccount.GUID,
		"password": serviceAccount.Password,
		"username": serviceAccount.Username,
	}
	res = append(res, temp)
	return res
}

func flattenWorkersGroup(workersGroups K8SGroupList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, worker := range workersGroups {
		temp := map[string]interface{}{
			"annotations":      worker.Annotations,
			"cpu":              worker.CPU,
			"detailed_info":    flattenDetailedInfo(worker.DetailedInfo, nil),
			"disk":             worker.Disk,
			"guid":             worker.GUID,
			"detailed_info_id": worker.ID,
			"labels":           worker.Labels,
			"name":             worker.Name,
			"num":              worker.Num,
			"ram":              worker.RAM,
			"taints":           worker.Taints,
		}
		res = append(res, temp)
	}
	return res
}

func flattenConfig(config interface{}) map[string]interface{} {
	return config.(map[string]interface{})
}

func flattenK8sItems(k8sItems K8SList) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, item := range k8sItems {
		temp := map[string]interface{}{
			"account_id":      item.AccountID,
			"account_name":    item.Name,
			"acl":             item.ACL,
			"bservice_id":     item.BServiceID,
			"ci_id":           item.CIID,
			"created_by":      item.CreatedBy,
			"created_time":    item.CreatedTime,
			"deleted_by":      item.DeletedBy,
			"deleted_time":    item.DeletedTime,
			"desc":            item.Description,
			"extnet_id":       item.ExtNetID,
			"gid":             item.GID,
			"guid":            item.GUID,
			"k8s_id":          item.ID,
			"lb_id":           item.LBID,
			"milestones":      item.Milestones,
			"k8s_name":        item.Name,
			"rg_id":           item.RGID,
			"rg_name":         item.RGName,
			"service_account": flattenServiceAccount(item.ServiceAccount),
			"status":          item.Status,
			"tech_status":     item.TechStatus,
			"updated_by":      item.UpdatedBy,
			"updated_time":    item.UpdatedTime,
			"vins_id":         item.VINSID,
			"workers_groups":  flattenWorkersGroup(item.WorkersGroup),
		}

		res = append(res, temp)
	}
	return res
}

func flattenK8sList(d *schema.ResourceData, k8sItems K8SList) {
	d.Set("items", flattenK8sItems(k8sItems))
}

func flattenResourceK8s(d *schema.ResourceData, k8s K8SRecord, masters []kvmvm.ComputeGetResp, workers []kvmvm.ComputeGetResp) {
	d.Set("acl", flattenAcl(k8s.ACL))
	d.Set("account_id", k8s.AccountID)
	d.Set("account_name", k8s.AccountName)
	d.Set("bservice_id", k8s.BServiceID)
	d.Set("created_by", k8s.CreatedBy)
	d.Set("created_time", k8s.CreatedTime)
	d.Set("deleted_by", k8s.DeletedBy)
	d.Set("deleted_time", k8s.DeletedTime)
	d.Set("k8s_ci_name", k8s.K8CIName)
	d.Set("masters", flattenMasterGroup(k8s.K8SGroups.Masters, masters))
	d.Set("workers", flattenK8sGroup(k8s.K8SGroups.Workers, workers))
	d.Set("lb_id", k8s.LBID)
	d.Set("rg_id", k8s.RGID)
	d.Set("rg_name", k8s.RGName)
	d.Set("status", k8s.Status)
	d.Set("tech_status", k8s.TechStatus)
	d.Set("updated_by", k8s.UpdatedBy)
	d.Set("updated_time", k8s.UpdatedTime)
	d.Set("default_wg_id", k8s.K8SGroups.Workers[0].ID)
}

func flattenWgData(d *schema.ResourceData, wg K8SGroup, computes []kvmvm.ComputeGetResp) {
	d.Set("annotations", wg.Annotations)
	d.Set("cpu", wg.CPU)
	d.Set("detailed_info", flattenDetailedInfo(wg.DetailedInfo, computes))
	d.Set("disk", wg.Disk)
	d.Set("guid", wg.GUID)
	d.Set("labels", wg.Labels)
	d.Set("name", wg.Name)
	d.Set("num", wg.Num)
	d.Set("ram", wg.RAM)
	d.Set("taints", wg.Taints)
}

func flattenWgList(wgList K8SGroupList, computesMap map[uint64][]kvmvm.ComputeGetResp) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, wg := range wgList {
		computes := computesMap[wg.ID]
		temp := map[string]interface{}{
			"annotations":   wg.Annotations,
			"cpu":           wg.CPU,
			"wg_id":         wg.ID,
			"detailed_info": flattenDetailedInfo(wg.DetailedInfo, computes),
			"disk":          wg.Disk,
			"guid":          wg.GUID,
			"labels":        wg.Labels,
			"name":          wg.Name,
			"num":           wg.Num,
			"ram":           wg.RAM,
			"taints":        wg.Taints,
		}

		res = append(res, temp)
	}
	return res
}

func flattenItemsWg(d *schema.ResourceData, wgList K8SGroupList, computes map[uint64][]kvmvm.ComputeGetResp) {
	d.Set("items", flattenWgList(wgList, computes))
}
