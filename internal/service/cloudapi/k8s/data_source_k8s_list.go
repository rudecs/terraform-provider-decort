package k8s

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rudecs/terraform-provider-decort/internal/constants"
)

func dataSourceK8sListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	k8sList, err := utilityK8sListCheckPresence(ctx, d, m, K8sListAPI)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())
	flattenK8sList(d, k8sList)

	return nil
}

func serviceAccountSchemaMake() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"username": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func k8sWorkersGroupsSchemaMake() map[string]*schema.Schema {
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
		"detailed_info_id": {
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

func createK8sListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"includedeleted": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"page": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"items": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"account_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"acl": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"bservice_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"ci_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"config": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
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
					"desc": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"extnet_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"gid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guid": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"k8s_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"lb_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"milestones": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"k8s_name": {
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
					"service_account": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: serviceAccountSchemaMake(),
						},
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
					"vins_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"workers_groups": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: k8sWorkersGroupsSchemaMake(),
						},
					},
				},
			},
		},
	}
}

func dataSourceK8sListSchemaMake() map[string]*schema.Schema {
	k8sListSchema := createK8sListSchema()
	return k8sListSchema
}

func DataSourceK8sList() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		ReadContext: dataSourceK8sListRead,

		Timeouts: &schema.ResourceTimeout{
			Read:    &constants.Timeout30s,
			Default: &constants.Timeout60s,
		},

		Schema: dataSourceK8sListSchemaMake(),
	}
}
