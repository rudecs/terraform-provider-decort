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
	log "github.com/sirupsen/logrus"

	"github.com/rudecs/terraform-provider-decort/internal/service/cloudapi/kvmvm"
)

func dataSourceK8sRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	k8s, err := utilityDataK8sCheckPresence(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatUint(k8s.ID, 10))

	k8sList, err := utilityK8sListCheckPresence(ctx, d, m, K8sListAPI)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	curK8s := K8SItem{}
	for _, k8sCluster := range k8sList {
		if k8sCluster.ID == k8s.ID {
			curK8s = k8sCluster
		}
	}
	if curK8s.ID == 0 {
		return diag.Errorf("Cluster with id %d not found in List clusters", k8s.ID)
	}
	d.Set("vins_id", curK8s.VINSID)

	masterComputeList := make([]kvmvm.ComputeGetResp, 0, len(k8s.K8SGroups.Masters.DetailedInfo))
	workersComputeList := make([]kvmvm.ComputeGetResp, 0, len(k8s.K8SGroups.Workers))
	for _, masterNode := range k8s.K8SGroups.Masters.DetailedInfo {
		compute, err := utilityComputeCheckPresence(ctx, d, m, masterNode.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		masterComputeList = append(masterComputeList, *compute)
	}
	for _, worker := range k8s.K8SGroups.Workers {
		for _, info := range worker.DetailedInfo {
			compute, err := utilityComputeCheckPresence(ctx, d, m, info.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			workersComputeList = append(workersComputeList, *compute)
		}
	}

	c := m.(*controller.ControllerCfg)
	urlValues := &url.Values{}
	urlValues.Add("k8sId", d.Id())
	kubeconfig, err := c.DecortAPICall(ctx, "POST", K8sGetConfigAPI, urlValues)
	if err != nil {
		log.Warnf("could not get kubeconfig: %v", err)
	}
	d.Set("kubeconfig", kubeconfig)

	urlValues = &url.Values{}
	urlValues.Add("lbId", strconv.FormatUint(k8s.LBID, 10))
	resp, err := c.DecortAPICall(ctx, "POST", LbGetAPI, urlValues)
	if err != nil {
		return diag.FromErr(err)
	}

	var lb LbRecord
	if err := json.Unmarshal([]byte(resp), &lb); err != nil {
		return diag.FromErr(err)
	}
	d.Set("extnet_id", lb.ExtNetID)
	d.Set("lb_ip", lb.PrimaryNode.FrontendIP)

	flattenK8sData(d, *k8s, masterComputeList, workersComputeList)
	return nil
}

func aclListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"explicit": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"right": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user_group_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func aclGroupSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclListSchemaMake(),
			},
		},
		"k8s_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclListSchemaMake(),
			},
		},
		"rg_acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclListSchemaMake(),
			},
		},
	}
}

func detailedInfoSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compute_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"interfaces": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: interfacesSchemaMake(),
			},
		},
		"natable_vins_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"natable_vins_network": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func interfacesSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"def_gw": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func masterGroupSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"detailed_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: detailedInfoSchemaMake(),
			},
		},
		"disk": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"master_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func k8sGroupListSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"annotations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"detailed_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: detailedInfoSchemaMake(),
			},
		},
		"disk": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
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

func dataSourceK8sSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"k8s_id": {
			Type:     schema.TypeInt,
			Required: true,
		},

		"acl": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: aclGroupSchemaMake(),
			},
		},
		"account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bservice_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"k8sci_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deleted_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"extnet_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the external network to connect workers to. If omitted network will be chosen by the platfom.",
		},
		"k8s_ci_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"masters": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: masterGroupSchemaMake(),
			},
		},
		"workers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: k8sGroupListSchemaMake(),
			},
		},
		"lb_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"lb_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address of default load balancer.",
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rg_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rg_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tech_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_by": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"kubeconfig": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kubeconfig for cluster access.",
		},
		"vins_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func DataSourceK8s() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sSchemaMake(),
	}
}
