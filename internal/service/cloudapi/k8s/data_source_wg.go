package k8s

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
	log "github.com/sirupsen/logrus"
)

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

func dataSourceK8sWgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debugf("dataSourceK8sWgRead: called with k8s id %d", d.Get("k8s_id").(int))

	k8s, err := utilityDataK8sCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(d.Get("wg_id").(int)))

	var id int
	if d.Id() != "" {
		id, err = strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		id = d.Get("wg_id").(int)
	}

	curWg := K8SGroup{}
	for _, wg := range k8s.K8SGroups.Workers {
		if wg.ID == uint64(id) {
			curWg = wg
			break
		}
	}
	if curWg.ID == 0 {
		return diag.Errorf("Not found wg with id: %v in k8s cluster: %v", id, k8s.ID)
	}

	workersComputeList := make([]kvmvm.ComputeGetResp, 0, 0)
	for _, info := range curWg.DetailedInfo {
		compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		workersComputeList = append(workersComputeList, *compute)
	}

	flattenWgData(d, curWg, workersComputeList)

	return nil
}

func dataSourceK8sWgSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of k8s instance.",
		},
		"wg_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of k8s worker Group.",
		},

		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the worker group.",
		},

		"num": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of worker nodes to create.",
		},

		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Worker node CPU count.",
		},

		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Worker node RAM in MB.",
		},

		"disk": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Worker node boot disk size. If unspecified or 0, size is defined by OS image size.",
		},
		"detailed_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: detailedInfoSchemaMake(),
			},
		},
		"labels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"annotations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"taints": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func DataSourceK8sWg() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sWgRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sWgSchemaMake(),
	}
}
