package kvmvm

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/status"
	log "github.com/sirupsen/logrus"
)

func flattenComputeDisksDemo(disksList []DiskRecord, extraDisks []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, disk := range disksList {
		if disk.Name == "bootdisk" || findInExtraDisks(disk.ID, extraDisks) { //skip main bootdisk and extraDisks
			continue
		}
		temp := map[string]interface{}{
			"disk_name": disk.Name,
			"disk_id":   disk.ID,
			"disk_type": disk.Type,
			"sep_id":    disk.SepID,
			"shareable": disk.Shareable,
			"size_max":  disk.SizeMax,
			"size_used": disk.SizeUsed,
			"pool":      disk.Pool,
			"desc":      disk.Desc,
			"image_id":  disk.ImageID,
			"size":      disk.SizeMax,
		}
		res = append(res, temp)
	}
	return res
}

func flattenCompute(d *schema.ResourceData, compFacts string) error {
	// This function expects that compFacts string contains response from API compute/get,
	// i.e. detailed information about compute instance.
	//
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceComputeExists(...) method
	model := ComputeGetResp{}
	log.Debugf("flattenCompute: ready to unmarshal string %s", compFacts)
	err := json.Unmarshal([]byte(compFacts), &model)
	if err != nil {
		return err
	}

	log.Debugf("flattenCompute: ID %d, RG ID %d", model.ID, model.RgID)

	d.SetId(fmt.Sprintf("%d", model.ID))
	// d.Set("compute_id", model.ID) - we should NOT set compute_id in the schema here: if it was set - it is already set, if it wasn't - we shouldn't
	d.Set("name", model.Name)
	d.Set("rg_id", model.RgID)
	d.Set("rg_name", model.RgName)
	d.Set("account_id", model.AccountID)
	d.Set("account_name", model.AccountName)
	d.Set("driver", model.Driver)
	d.Set("cpu", model.Cpu)
	d.Set("ram", model.Ram)
	// d.Set("boot_disk_size", model.BootDiskSize) - bootdiskSize key in API compute/get is always zero, so we set boot_disk_size in another way
	if model.VirtualImageID != 0 {
		d.Set("image_id", model.VirtualImageID)
	} else {
		d.Set("image_id", model.ImageID)
	}
	d.Set("description", model.Desc)
	d.Set("enabled", false)
	if model.Status == status.Enabled {
		d.Set("enabled", true)
	}

	//d.Set("cloud_init", "applied") // NOTE: for existing compute we hard-code this value as an indicator for DiffSuppress fucntion
	//d.Set("status", model.Status)
	//d.Set("tech_status", model.TechStatus)
	d.Set("started", false)
	if model.TechStatus == "STARTED" {
		d.Set("started", true)
	}

	bootDisk := findBootDisk(model.Disks)

	d.Set("boot_disk_size", bootDisk.SizeMax)
	d.Set("boot_disk_id", bootDisk.ID) // we may need boot disk ID in resize operations
	d.Set("sep_id", bootDisk.SepID)
	d.Set("pool", bootDisk.Pool)

	//if len(model.Disks) > 0 {
	//log.Debugf("flattenCompute: calling parseComputeDisksToExtraDisks for %d disks", len(model.Disks))
	//if err = d.Set("extra_disks", parseComputeDisksToExtraDisks(model.Disks)); err != nil {
	//return err
	//}
	//}

	if len(model.Interfaces) > 0 {
		log.Debugf("flattenCompute: calling parseComputeInterfacesToNetworks for %d interfaces", len(model.Interfaces))
		if err = d.Set("network", parseComputeInterfacesToNetworks(model.Interfaces)); err != nil {
			return err
		}
	}

	if len(model.OsUsers) > 0 {
		log.Debugf("flattenCompute: calling parseOsUsers for %d logins", len(model.OsUsers))
		if err = d.Set("os_users", parseOsUsers(model.OsUsers)); err != nil {
			return err
		}
	}

	err = d.Set("disks", flattenComputeDisksDemo(model.Disks, d.Get("extra_disks").(*schema.Set).List()))
	if err != nil {
		return err
	}

	return nil
}
