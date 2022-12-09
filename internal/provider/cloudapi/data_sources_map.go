/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>

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

package cloudapi

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/account"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/bservice"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/disks"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/extnet"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/image"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/k8s"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/lb"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/locations"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/rg"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/snapshot"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/vins"
)

func NewDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"decort_account":                        account.DataSourceAccount(),
		"decort_resgroup":                       rg.DataSourceResgroup(),
		"decort_kvmvm":                          kvmvm.DataSourceCompute(),
		"decort_k8s":                            k8s.DataSourceK8s(),
		"decort_k8s_list":                       k8s.DataSourceK8sList(),
		"decort_k8s_list_deleted":               k8s.DataSourceK8sListDeleted(),
		"decort_k8s_wg":                         k8s.DataSourceK8sWg(),
		"decort_k8s_wg_list":                    k8s.DataSourceK8sWgList(),
		"decort_vins":                           vins.DataSourceVins(),
		"decort_snapshot_list":                  snapshot.DataSourceSnapshotList(),
		"decort_disk":                           disks.DataSourceDisk(),
		"decort_disk_list":                      disks.DataSourceDiskList(),
		"decort_rg_list":                        rg.DataSourceRgList(),
		"decort_disk_list_types_detailed":       disks.DataSourceDiskListTypesDetailed(),
		"decort_disk_list_types":                disks.DataSourceDiskListTypes(),
		"decort_disk_list_deleted":              disks.DataSourceDiskListDeleted(),
		"decort_disk_list_unattached":           disks.DataSourceDiskListUnattached(),
		"decort_disk_snapshot":                  disks.DataSourceDiskSnapshot(),
		"decort_disk_snapshot_list":             disks.DataSourceDiskSnapshotList(),
		"decort_account_list":                   account.DataSourceAccountList(),
		"decort_account_computes_list":          account.DataSourceAccountComputesList(),
		"decort_account_disks_list":             account.DataSourceAccountDisksList(),
		"decort_account_vins_list":              account.DataSourceAccountVinsList(),
		"decort_account_audits_list":            account.DataSourceAccountAuditsList(),
		"decort_account_rg_list":                account.DataSourceAccountRGList(),
		"decort_account_consumed_units":         account.DataSourceAccountConsumedUnits(),
		"decort_account_consumed_units_by_type": account.DataSourceAccountConsumedUnitsByType(),
		"decort_account_reserved_units":         account.DataSourceAccountReservedUnits(),
		"decort_account_templates_list":         account.DataSourceAccountTemplatessList(),
		"decort_account_deleted_list":           account.DataSourceAccountDeletedList(),
		"decort_account_flipgroups_list":        account.DataSourceAccountFlipGroupsList(),
		"decort_bservice_list":                  bservice.DataSourceBasicServiceList(),
		"decort_bservice":                       bservice.DataSourceBasicService(),
		"decort_bservice_snapshot_list":         bservice.DataSourceBasicServiceSnapshotList(),
		"decort_bservice_group":                 bservice.DataSourceBasicServiceGroup(),
		"decort_bservice_deleted_list":          bservice.DataSourceBasicServiceDeletedList(),
		"decort_extnet_list":                    extnet.DataSourceExtnetList(),
		"decort_extnet_computes_list":           extnet.DataSourceExtnetComputesList(),
		"decort_extnet":                         extnet.DataSourceExtnet(),
		"decort_extnet_default":                 extnet.DataSourceExtnetDefault(),
		"decort_vins_list":                      vins.DataSourceVinsList(),
		"decort_locations_list":                 locations.DataSourceLocationsList(),
		"decort_location_url":                   locations.DataSourceLocationUrl(),
		"decort_image_list":                     image.DataSourceImageList(),
		"decort_image":                          image.DataSourceImage(),
		"decort_lb":                             lb.DataSourceLB(),
		"decort_lb_list":                        lb.DataSourceLBList(),
		"decort_lb_list_deleted":                lb.DataSourceLBListDeleted(),
		// "decort_pfw": dataSourcePfw(),
	}

}
