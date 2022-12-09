package k8s

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
	"github.com/rudecs/terraform-provider-decort/internal/controller"
	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

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

func utilityK8sWgListCheckPresence(ctx context.Context, d *schema.ResourceData, m interface{}) (K8SGroupList, error) {

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", strconv.Itoa(d.Get("k8s_id").(int)))

	resp, err := c.DecortAPICall(ctx, "POST", K8sGetAPI, urlValues)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, nil
	}

	var k8s K8SRecord
	if err := json.Unmarshal([]byte(resp), &k8s); err != nil {
		return nil, err
	}

	return k8s.K8SGroups.Workers, nil
}

func dataSourceK8sWgListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	wgList, err := utilityK8sWgListCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(d.Get("k8s_id").(int)))

	workersComputeList := make(map[uint64][]kvmvm.ComputeGetResp)
	for _, worker := range wgList {
		workersComputeList[worker.ID] = make([]kvmvm.ComputeGetResp, 0, len(worker.DetailedInfo))
		for _, info := range worker.DetailedInfo {
			compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			workersComputeList[worker.ID] = append(workersComputeList[worker.ID], *compute)
		}
	}
	flattenItemsWg(d, wgList, workersComputeList)
	return nil
}

func wgSchemaMake() map[string]*schema.Schema {
	wgSchema := dataSourceK8sWgSchemaMake()
	delete(wgSchema, "k8s_id")
	wgSchema["wg_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "ID of k8s worker Group.",
	}
	return wgSchema
}

func dataSourceK8sWgListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: wgSchemaMake(),
			},
		},
	}
}

func DataSourceK8sWgList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sWgListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sWgListSchemaMake(),
	}
}
