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

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/account"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/bservice"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/disks"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/extnet"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/rg"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/snapshot"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/vgpu"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/vins"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/grid"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/image"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/pcidevice"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudbroker/sep"
)

func NewDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"decort_account":  account.DataSourceAccount(),
		"decort_resgroup": rg.DataSourceResgroup(),
		// "decort_kvmvm":                          dataSourceCompute(),
		"decort_image":                          image.DataSourceImage(),
		"decort_disk":                           disks.DataSourceDisk(),
		"decort_vins":                           vins.DataSourceVins(),
		"decort_grid":                           grid.DataSourceGrid(),
		"decort_grid_list":                      grid.DataSourceGridList(),
		"decort_image_list":                     image.DataSourceImageList(),
		"decort_image_list_stacks":              image.DataSourceImageListStacks(),
		"decort_snapshot_list":                  snapshot.DataSourceSnapshotList(),
		"decort_vgpu":                           vgpu.DataSourceVGPU(),
		"decort_pcidevice":                      pcidevice.DataSourcePcidevice(),
		"decort_pcidevice_list":                 pcidevice.DataSourcePcideviceList(),
		"decort_sep_list":                       sep.DataSourceSepList(),
		"decort_sep":                            sep.DataSourceSep(),
		"decort_sep_consumption":                sep.DataSourceSepConsumption(),
		"decort_sep_disk_list":                  sep.DataSourceSepDiskList(),
		"decort_sep_config":                     sep.DataSourceSepConfig(),
		"decort_sep_pool":                       sep.DataSourceSepPool(),
		"decort_disk_list":                      disks.DataSourceDiskList(),
		"decort_rg_list":                        rg.DataSourceRgList(),
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
		// "decort_pfw": dataSourcePfw(),
	}

}
